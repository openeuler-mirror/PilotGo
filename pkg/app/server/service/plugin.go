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
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Enabled     int    `json:"enabled"`
	Status      string `json:"status"`
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
			Enabled:     p.Enabled,
			Status:      splugin.StatusOffline,
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

	plugins := []*Plugin{}
	for _, value := range pm.loadedPlugin {
		p := &Plugin{
			Name:        value.Name,
			Version:     value.Version,
			Description: value.Description,
			Url:         value.Url,
		}

		plugins = append(plugins, p)
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

// 更新插件使能状态
func (pm *PluginManager) UpdatePlugin(uuid string, enable int) error {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	// TODO: 更新插件状态
	return nil
}

// 检查插件是否使能
func (pm *PluginManager) IsPluginEnabled(uuid string) bool {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	// TODO: 查询插件状态
	return false
}

func GetManager() *PluginManager {
	return globalManager
}

// 请求plugin的接口服务，获取接口信息
func CheckPlugin(url string) (*Plugin, error) {
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Debug("")
		return nil, err
	}

	plugin := &Plugin{
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Url:         info.Url,
		Status:      splugin.StatusLoaded,
	}

	return plugin, nil
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
	// TODO: check info valid
	return info, nil
}

// 获取plugin清单
func GetPlugins() []*Plugin {
	return globalManager.GetAll()
}

func AddPlugin(url string) error {
	logger.Debug("add login from %s", url)

	plugin, err := CheckPlugin(url + "/plugin_manage/info")
	if err != nil {
		return err
	}

	if err := globalManager.Add(plugin); err != nil {
		return err
	}
	return nil
}

func DeletePlugin(id int) error {
	logger.Debug("delete login: %d", id)

	if err := globalManager.Remove("id"); err != nil {
		return err
	}
	return nil
}

func TogglePlugin(uuid string, enable int) error {
	logger.Debug("toggle plugin: %s to enable %d ", uuid, enable)

	if err := globalManager.UpdatePlugin(uuid, enable); err != nil {
		return err
	}
	return nil
}
