package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
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

// type Plugin = dao.PluginModel
func Init() error {
	if err := mysqlmanager.MySQL().AutoMigrate(&dao.PluginModel{}); err != nil {
		return err
	}

	if err := globalPluginManager.recovery(); err != nil {
		return err
	}
	return nil
}

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
	Extentions  []*common.Extention
}

func (p *Plugin) Clone() *Plugin {
	result := &Plugin{
		UUID:        p.UUID,
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Email:       p.Email,
		Url:         p.Url,
		PluginType:  p.PluginType,
		Enabled:     p.Enabled,
		Extentions:  []*common.Extention{},
	}
	for _, e := range p.Extentions {
		result.Extentions = append(result.Extentions, &common.Extention{
			PluginName: e.PluginName,
			Name:       e.Name,
			Type:       e.Type,
			URL:        e.URL,
		})
	}

	return result
}

type PluginManager struct {
	sync.Mutex

	Plugins []*Plugin
}

var globalPluginManager = &PluginManager{
	Mutex:   sync.Mutex{},
	Plugins: []*Plugin{},
}

// 从DB中恢复插件信息
func (m *PluginManager) recovery() error {
	plugins, err := dao.QueryPlugins()
	if err != nil {
		logger.Error("failed to recovery plugin info from db")
		return nil
	}

	for _, p := range plugins {
		logger.Debug("recovery plugin:%s %s %d", p.UUID, p.Url, p.Enabled)
		err := m.updatePlugin(p.UUID, p.Url, p.Enabled)
		if err != nil {
			logger.Error("failed to update plugin %s info", p.Name)
			// 插件离线等异常状况
			// 使用历史数据
			m.Lock()
			m.Plugins = append(m.Plugins, &Plugin{
				UUID:        p.UUID,
				Name:        p.Name,
				Version:     p.Version,
				Description: p.Description,
				Author:      p.Author,
				Email:       p.Email,
				Url:         p.Url,
				PluginType:  p.PluginType,
				Enabled:     p.Enabled,
			})
			m.Unlock()

			// 继续恢复下一个plugin
		}
	}

	logger.Debug("finish recovery")
	return nil
}

// 根据url查询最新的plugin信息，更新到指定的uuid记录当中
func (m *PluginManager) updatePlugin(uuid, url string, enabled int) error {
	// 查询最新的插件信息
	logger.Debug("update plugin")
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Error("failed to request plugin info:%s", err.Error())
		return err
	}

	p := &Plugin{
		UUID:        uuid,
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Email:       info.Email,
		Url:         info.Url,
		PluginType:  info.PluginType,
		Enabled:     enabled,
		Extentions:  info.Extentions,
	}

	if err := dao.UpdatePluginInfo(toPluginDao(p)); err != nil {
		return err
	}

	found := false
	m.Lock()
	for i, v := range m.Plugins {
		if v.UUID == uuid {
			m.Plugins[i] = p
			break
		}
	}
	if !found {
		m.Plugins = append(m.Plugins, p)
	}
	m.Unlock()

	return nil
}

func (m *PluginManager) addPlugin(url string) error {
	err := Handshake(url)
	if err != nil {
		return err
	}

	p, err := requestPluginInfo(url)
	if err != nil {
		return err
	}

	if p.UUID == "" {
		p.UUID = uuid.New().String()
	}

	if err := dao.RecordPlugin(toPluginDao(p)); err != nil {
		return err
	}

	m.Lock()
	m.Plugins = append(m.Plugins, p)
	m.Unlock()

	return nil
}

func (m *PluginManager) deletePlugin(uuid string) error {
	if err := dao.DeletePlugin(uuid); err != nil {
		logger.Error("failed to delete plugin info:%s", err.Error())
		return err
	}

	m.Lock()
	index := 0
	for i, v := range m.Plugins {
		if v.UUID == uuid {
			index = i
			break
		}
	}

	if index == 0 {
		m.Plugins = m.Plugins[1:]
	} else if index == len(m.Plugins)-1 {
		m.Plugins = m.Plugins[:index]
	} else {
		m.Plugins = append(m.Plugins[:index], m.Plugins[index+1:]...)
	}
	m.Unlock()
	return nil
}

func (m *PluginManager) togglePlugin(uuid string, enable int) error {
	if err := dao.UpdatePluginEnabled(&dao.PluginModel{
		UUID:    uuid,
		Enabled: enable,
	}); err != nil {
		return err
	}

	url := ""
	m.Lock()
	for _, v := range m.Plugins {
		if v.UUID == uuid {
			v.Enabled = enable
			url = v.Url
			break
		}
	}
	m.Unlock()

	// 更新最新的插件信息
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Error("failed to request plugin info:%s", err.Error())
		return errors.New("plugin offline or unreachable")
	}
	m.Lock()
	for _, v := range m.Plugins {
		if v.UUID == uuid {
			v.Name = info.Name
			v.Version = info.Version
			v.Description = info.Description
			v.Author = info.Author
			v.Email = info.Email
			v.Url = info.Url
			v.PluginType = info.PluginType
			v.Extentions = info.Extentions
			break
		}
	}
	m.Unlock()

	return nil
}

func (m *PluginManager) getPlugin(name string) (*Plugin, error) {
	var result *Plugin
	found := false
	m.Lock()
	for _, v := range m.Plugins {
		if v.Name == name {
			// 使用深拷贝避免指针泄露
			result = v.Clone()
			found = true
			break
		}
	}
	m.Unlock()

	if !found {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return result, nil
}

// 获取所有插件
func (m *PluginManager) getPlugins() ([]*Plugin, error) {
	result := []*Plugin{}
	m.Lock()
	for _, v := range m.Plugins {
		// 使用深拷贝避免指针泄露
		result = append(result, v.Clone())
	}
	m.Unlock()

	return result, nil
}

func toPluginDao(p *Plugin) *dao.PluginModel {
	return &dao.PluginModel{
		UUID:        p.UUID,
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Email:       p.Email,
		Url:         p.Url,
		PluginType:  p.PluginType,
		Enabled:     p.Enabled,
	}
}

// 与plugin进行握手，绑定PilotGo与server端
func Handshake(url string) error {
	index := strings.Index(url, "plugin")
	if index > 0 {
		url = url[:index]
	}
	port := strings.Split(config.Config().HttpServer.Addr, ":")[1]

	url = strings.TrimRight(url, "/") + "/plugin_manage/bind?port=" + port
	logger.Debug("plugin url is:%s", url)

	resp, err := httputils.Put(url, nil)
	if err != nil {
		logger.Debug("request plugin info error:%s", err.Error())
		return err
	}

	d := &struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}{}
	err = json.Unmarshal(resp.Body, d)
	if err != nil {
		logger.Error("unmarshal request plugin info error:%s", err.Error())
		return err
	}
	if d.Code != http.StatusOK {
		return errors.New(d.Message)
	}

	return nil
}

// 获取到插件的基本信息
func requestPluginInfo(url string) (*Plugin, error) {
	index := strings.Index(url, "plugin")
	if index > 0 {
		url = url[:index]
	}
	url = strings.TrimRight(url, "/") + "/plugin_manage/info"
	logger.Debug("plugin url is:%s", url)

	resp, err := httputils.Get(url, nil)
	if err != nil {
		logger.Error("request plugin info error:%s", err.Error())
		return nil, err
	}

	info := &client.PluginFullInfo{}
	err = json.Unmarshal(resp.Body, info)
	if err != nil {
		logger.Error("unmarshal request plugin info error:%s", err.Error())
		return nil, err
	}
	// TODO: check info valid

	return &Plugin{
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Email:       info.Email,
		Url:         info.Url,
		PluginType:  info.PluginType,
		Extentions:  info.Extentions,
		// Status:      common.StatusLoaded,
	}, nil
}

// 获取单个plugin
func GetPlugin(name string) (*Plugin, error) {
	plugin, err := globalPluginManager.getPlugin(name)
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return nil, err
	}
	return plugin, nil
}

// 获取所有的plugin
func GetPlugins() ([]*Plugin, error) {
	plugins, err := globalPluginManager.getPlugins()
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return nil, err
	}

	return plugins, nil
}

// 分页查询
func GetPluginPaged(offset, size int) (int64, []*Plugin, error) {
	// 借助db的分页功能实现分页查询
	total, plugins, error := dao.GetPluginPaged(offset, size)

	result := []*Plugin{}
	for _, p := range plugins {
		plugin, err := globalPluginManager.getPlugin(p.Name)
		if err != nil {
			logger.Error("manager get plugin %s", err)
			continue
		}
		result = append(result, plugin)
	}

	return total, result, error
}

type AddPluginParam struct {
	Name string `json:"name"`
	Type string `json:"plugin_type"`
	Url  string `json:"url"`
}

func AddPlugin(param *AddPluginParam) error {
	url := param.Url
	logger.Debug("add plugin from %s", url)

	if err := globalPluginManager.addPlugin(url); err != nil {
		return err
	}
	return nil
}

func DeletePlugin(uuid string) error {
	logger.Debug("delete plugin: %s", uuid)

	if err := globalPluginManager.deletePlugin(uuid); err != nil {
		logger.Error("failed to delete plugin info:%s", err.Error())
	}
	return nil
}

func TogglePlugin(uuid string, enable int) error {
	logger.Debug("toggle plugin: %s to enable %d ", uuid, enable)
	if err := globalPluginManager.togglePlugin(uuid, enable); err != nil {
		return err
	}

	return nil
}
