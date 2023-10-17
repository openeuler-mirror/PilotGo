package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/google/uuid"
)

const (
	PluginEnabled  = 1
	PluginDisabled = 0
)

type Plugin struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Url         string `json:"url"`
	PluginType  string `json:"plugin_type"`
	Enabled     int    `json:"enabled"`
	Status      string `json:"status"`
}

// 初始化插件服务
func ServiceInit() error {
	return globalManager.RestorePluginInfo()
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
func (pm *PluginManager) RestorePluginInfo() error {
	plugins, err := dao.QueryPlugins()
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return err
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()
	for _, p := range plugins {
		np := &Plugin{
			UUID:        p.UUID,
			Name:        p.Name,
			Version:     p.Version,
			Description: p.Description,
			Author:      p.Author,
			Email:       p.Email,
			Url:         p.Url,
			PluginType:  p.PluginType,
			Enabled:     p.Enabled,
			Status:      common.StatusOffline,
		}

		pm.loadedPlugin[np.Name] = np
	}

	return nil
}

// 添加一个插件
func (pm *PluginManager) Add(p *Plugin) error {
	if p.Url == "" || p.Name == "" || p.PluginType == "" {
		return errors.New("invalid plugin parameter")
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()
	if _, ok := pm.loadedPlugin[p.Name]; ok {
		return errors.New("plugin already registered")
	}

	// 记录到DB当中
	err := dao.RecordPlugin(&dao.PluginModel{
		UUID:        p.UUID,
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Email:       p.Email,
		Url:         p.Url,
		PluginType:  p.PluginType,
		Enabled:     PluginDisabled,
	})
	if err != nil {
		return err
	}

	// 记录db无报错之后才更新缓存数据
	pm.loadedPlugin[p.Name] = p

	return nil
}

// 移除注册的插件
func (pm *PluginManager) Remove(uuid string) error {
	if err := dao.DeletePlugin(uuid); err != nil {
		// TODO: 细化db删除错误处理
		logger.Error("failed to delete plugin info:%s", err.Error())
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()
	name := ""
	for _, p := range pm.loadedPlugin {
		if p.UUID == uuid {
			name = p.Name
			break
		}
	}
	delete(pm.loadedPlugin, name)

	return errors.New("plugin not found")
}

// 获取注册的插件
func (pm *PluginManager) Get(name string) (*Plugin, error) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	if p, ok := pm.loadedPlugin[name]; ok {
		return p, nil
	}

	return nil, errors.New("plugin not found")
}

// 获取所有的插件
func (pm *PluginManager) GetAll() []*Plugin {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	plugins := []*Plugin{}
	for _, value := range pm.loadedPlugin {
		p := &Plugin{
			UUID:        value.UUID,
			Name:        value.Name,
			Version:     value.Version,
			Description: value.Description,
			Author:      value.Author,
			Email:       value.Email,
			Url:         value.Url,
			PluginType:  value.PluginType,
			Enabled:     value.Enabled,
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
	var status int
	if enable == PluginEnabled {
		status = PluginDisabled
	} else {
		status = PluginEnabled
	}
	if err := dao.UpdatePluginEnabled(&dao.PluginModel{
		UUID:    uuid,
		Enabled: status,
	}); err != nil {
		return err
	}

	pm.lock.RLock()
	defer pm.lock.RUnlock()
	for _, p := range pm.loadedPlugin {
		if p.UUID == uuid {
			p.Enabled = status
		}
		return nil
	}

	return nil
}

// 检查插件是否使能
func (pm *PluginManager) IsPluginEnabled(uuid string) (int, error) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	pm.lock.RLock()
	defer pm.lock.RUnlock()
	for _, p := range pm.loadedPlugin {
		if p.UUID == uuid {
			return p.Enabled, nil
		}
	}
	return PluginDisabled, errors.New("plugin not found")
}

func GetManager() *PluginManager {
	return globalManager
}

// 与plugin进行握手，交换必要信息
func Handshake(url string) (*Plugin, error) {
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Debug("")
		return nil, err
	}

	plugin := &Plugin{
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Email:       info.Email,
		Url:         info.Url,
		PluginType:  info.PluginType,
		Status:      common.StatusLoaded,
	}

	return plugin, nil
}

// 发起http请求，提供server地址，同时获取到插件的基本信息
func requestPluginInfo(url string) (*client.PluginInfo, error) {
	conf := config.Config().HttpServer
	url = url + fmt.Sprintf("?server=%s", conf.Addr)
	resp, err := httputils.Get(url, nil)
	if err != nil {
		logger.Debug("request plugin info error:%s", err.Error())
		return nil, err
	}

	info := &client.PluginInfo{}
	err = json.Unmarshal(resp.Body, info)
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

func GetPlugin(name string) (*Plugin, error) {
	return globalManager.Get(name)
}

type AddPluginParam struct {
	Name string `json:"name"`
	Type string `json:"plugin_type"`
	Url  string `json:"url"`
}

func AddPlugin(param *AddPluginParam) error {
	url := param.Url
	logger.Debug("add plugin from %s", url)
	url = strings.TrimRight(url, "/")

	plugin, err := Handshake(url + "/plugin_manage/info")
	if err != nil {
		return err
	}
	plugin.UUID = uuid.New().String()
	plugin.PluginType = param.Type

	if err := globalManager.Add(plugin); err != nil {
		return err
	}
	return nil
}

func DeletePlugin(uuid string) error {
	logger.Debug("delete plugin: %s", uuid)

	if err := globalManager.Remove(uuid); err != nil {
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
