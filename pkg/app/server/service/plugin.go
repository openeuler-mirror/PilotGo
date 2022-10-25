package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	splugin "openeuler.org/PilotGo/plugin-sdk/plugin"
)

type Plugin struct {
	Name        string
	Version     string
	Description string
	Url         string
	Status      int
}

// 初始化插件服务
func PluginServiceInit() error {
	return restorePluginInfo()
}

// TODO： 替换成concurrent hashmap
type PluginManager struct {
	loadedPlugin map[string]*Plugin
	lock         sync.RWMutex
}

var globalManager = &PluginManager{
	loadedPlugin: map[string]*Plugin{},
	lock:         sync.RWMutex{},
}

// 从db中恢复插件信息
func restorePluginInfo() error {
	plugins, err := dao.QueryPlugins()
	if err != nil {
		return err
	}

	for _, p := range plugins {
		np := &Plugin{
			Name:        p.Name,
			Version:     p.Version,
			Description: p.Version,
			Url:         p.Url,
			Status:      p.Status,
		}

		globalManager.Add(np)
	}

	return nil
}

// 添加一个插件
func (pm *PluginManager) Add(p *Plugin) error {
	if p.Url == "" || p.Name == "" {
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

// 移除注册的插件
func (pm *PluginManager) Remove(name string) error {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	if _, ok := pm.loadedPlugin[name]; ok {
		delete(pm.loadedPlugin, name)
		return nil
	}

	return errors.New("plugin not found")
}

// 获取注册的插件
func (pm *PluginManager) Get(name string) (string, error) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	if p, ok := pm.loadedPlugin[name]; ok {
		return p.Url, nil
	}

	return "", errors.New("plugin not found")
}

// 获取所有的插件
func (pm *PluginManager) GetAll() []*Plugin {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	pluginLen := len(pm.loadedPlugin)
	if pluginLen == 0 {
		return nil
	}

	plugins := make([]*Plugin, pluginLen)
	i := 0
	for _, value := range pm.loadedPlugin {
		plugins[i].Name = value.Name
		plugins[i].Version = value.Version
		plugins[i].Description = value.Description
		plugins[i].Url = value.Url
		i++
	}

	return plugins
}

// 检查插件是否注册
func (pm *PluginManager) Check(name string) bool {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	_, ok := pm.loadedPlugin[name]
	return ok
}

func GetManager() *PluginManager {
	return globalManager
}

// 请求plugin的接口服务，获取接口信息
func CheckPlugin(url string) (*splugin.PluginInfo, error) {
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Debug("")
		return nil, err
	}

	return info, nil
}

// 发起http请求获取到插件的基本信息
func requestPluginInfo(url string) (*splugin.PluginInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Debug("request plugin info error:%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("read request plugin info body error:%s", err.Error())
		return nil, err
	}

	info := &splugin.PluginInfo{}
	err = json.Unmarshal(body, info)
	if err != nil {
		logger.Debug("unmarshal request plugin info error:%s", err.Error())
	}
	return info, nil
}

// 获取plugin清单
func GetPlugins() []*Plugin {
	return globalManager.GetAll()
}

func AddPlugin(url string) error {
	return nil
}
