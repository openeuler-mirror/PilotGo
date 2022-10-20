package plugin

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	go_plugin "github.com/hashicorp/go-plugin"
	plugin_sdk "github.com/liweifeng1/plugin-sdk"
	"io/ioutil"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"os"
	"os/exec"
	"sync"
)

type PluginState int64

const (
	Stopped = iota
	Running
)

type PluginInfo struct {
	Conn     plugin_sdk.PluginInterface
	Client   *go_plugin.Client
	State    PluginState
	Filename string
	Config   []plugin_sdk.PluginConfig
	Manifest plugin_sdk.PluginManifest
}

type PluginManager struct {
	loadedPlugin map[string]PluginInfo
	lock         sync.RWMutex
}

var globalManager *PluginManager
var pluginMap = map[string]go_plugin.Plugin{
	"pilotGo_plugin": &plugin_sdk.PilotGoPlugin{},
}

func init() {
	globalManager = &PluginManager{
		loadedPlugin: map[string]PluginInfo{},
		lock:         sync.RWMutex{},
	}
}

func newClient(name string, logger hclog.Logger) *go_plugin.Client {
	return go_plugin.NewClient(&go_plugin.ClientConfig{
		HandshakeConfig: plugin_sdk.HandshakeConfig,
		Plugins:         pluginMap,
		Managed:         true,
		Cmd:             exec.Command(name),
		Logger:          logger,
	})
}

func hostPlugin(name string) PluginInfo {
	hclogger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})
	client := newClient(fmt.Sprintf("./plugins/%s", name), hclogger)
	rpcClient, err := client.Client()
	if err != nil {
		logger.Error("load plugin %s failed:%s", name, err.Error())
	}
	raw, err := rpcClient.Dispense("pilotGo_plugin")
	if err != nil {
		logger.Error("load plugin %s failed:%s", name, err.Error())
	}
	p := raw.(plugin_sdk.PluginInterface)
	err = p.OnLoad()
	if err != nil {
		logger.Error("load plugin %s failed:%s", name, err.Error())
	}
	return PluginInfo{
		Conn:     p,
		Client:   client,
		State:    Running,
		Filename: name,
		Manifest: p.GetManifest(),
		Config:   p.GetConfiguration(),
	}
}

func (pm *PluginManager) Load(name string) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	p := hostPlugin(name)

	if pm.IsPluginLoaded(name) {
		logger.Error("Cannot load plugin '%s' as it is already loading!\n", name)
		p.Client.Kill()
		return errors.New("cannot load plugin twice")
	}
	pm.loadedPlugin[name] = p
	return nil
}

func (pm *PluginManager) UnLoad(name string) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	if p, ok := pm.loadedPlugin[name]; ok {
		if p.Client != nil {
			err := p.Conn.OnClose()
			if err != nil {
				return errors.New("cannot close plugin")
			}
			p.Client.Kill()
		}
		delete(pm.loadedPlugin, name)
		return nil
	}
	return errors.New("plugin not found")
}

func (pm *PluginManager) IsPluginLoaded(name string) bool {
	if p, ok := pm.loadedPlugin[name]; ok && p.Conn != nil {
		return true
	}
	return false
}

func (pm *PluginManager) GetWebExtension(name string) (string, error) {
	if p, ok := pm.loadedPlugin[name]; ok {
		address := "http://127.0.0.1"
		port := ""
		for _, config := range p.Config {
			if config.Type == plugin_sdk.PortValue {
				port = config.Values
			} else if config.Type == plugin_sdk.UrlValue {
				address = config.Values
			}
		}

		html := "/"
		for _, extension := range p.Conn.GetWebExtension() {
			if extension.Type == plugin_sdk.HTML {
				html = extension.PathMatchRegex
			}
		}
		if port == "" {
			return address + html, nil
		}
		return address + ":" + port + html, nil
	}
	return "", errors.New("plugin not found")
}

func (pm *PluginManager) GetAll() ([]string, error) {
	var plugins []string
	files, err := ioutil.ReadDir("plugins")
	if err != nil {
		return plugins, errors.New("scan plugins failed")
	}
	for _, file := range files {
		plugins = append(plugins, file.Name())
	}
	return plugins, nil
}

func GetManager() *PluginManager {
	return globalManager
}
