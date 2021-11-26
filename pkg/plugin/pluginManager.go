package plugin

import (
	"errors"
	"sync"
)

type Plugin struct {
	Name        string
	Version     string
	Description string
	Depends     []*Plugin
	Url         string
	Port        string
	Protocol    string
	Status      string
}

type PluginManager struct {
	loadedPlugin map[string]*Plugin
	lock         sync.RWMutex
}

var globalManager *PluginManager

func init() {
	globalManager = &PluginManager{
		loadedPlugin: map[string]*Plugin{},
		lock:         sync.RWMutex{},
	}
}

// Load 注册一个插件
func (pm *PluginManager) Regist(p *Plugin) error {
	if p.Url == "" || p.Port == "" || p.Name == "" || p.Protocol == "" {
		return errors.New("invalid plugin parameter")
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()

	if _, ok := pm.loadedPlugin[p.Name]; ok {
		return errors.New("plugin already registered")
	}
	pm.loadedPlugin[p.Name] = p

	return nil
}

// Remove 移除注册的插件
func (pm *PluginManager) Remove(name string) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	if _, ok := pm.loadedPlugin[name]; ok {
		delete(pm.loadedPlugin, name)
		return nil
	}

	return errors.New("plugin not found")
}

// Remove 获取注册的插件
func (pm *PluginManager) Get(name string) (string, string, string, error) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	if p, ok := pm.loadedPlugin[name]; ok {
		return p.Url, p.Port, p.Protocol, nil
	}

	return "", "", "", errors.New("plugin not found")
}

func (pm *PluginManager) GetAll() []Plugin {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	pluginLen := len(pm.loadedPlugin)
	if pluginLen == 0 {
		return nil
	}

	plugins := make([]Plugin, pluginLen)
	i := 0
	for _, value := range pm.loadedPlugin {
		plugins[i].Name = value.Name
		plugins[i].Version = value.Version
		plugins[i].Description = value.Description
		plugins[i].Depends = value.Depends
		plugins[i].Url = value.Url
		plugins[i].Port = value.Port
		plugins[i].Protocol = value.Protocol
		i++
	}

	return plugins
}

// Check 检查插件是否注册
func (pm *PluginManager) Check(name string) bool {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	_, ok := pm.loadedPlugin[name]
	return ok
}

func GetManager() *PluginManager {
	return globalManager
}
