/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/config"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/redismanager"
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

func Init(stopCh <-chan struct{}) error {
	if err := mysqlmanager.MySQL().AutoMigrate(&dao.PluginModel{}); err != nil {
		return err
	}

	if err := globalPluginManager.recovery(); err != nil {
		return err
	}

	// 检查插件状态，重新绑定plugin与pilotgo
	CheckPluginHeartbeats(stopCh)

	return nil
}

type Plugin struct {
	UUID              string              `json:"uuid"`
	CustomName        string              `json:"custom_name"`
	Name              string              `json:"name"`
	Version           string              `json:"version"`
	Description       string              `json:"description"`
	Author            string              `json:"author"`
	Email             string              `json:"email"`
	Url               string              `json:"url"`
	PluginType        string              `json:"plugin_type"`
	ConnectStatus     bool                `json:"status"`
	LastHeartbeatTime string              `json:"lastheatbeat"`
	Enabled           int                 `json:"enabled"`
	Extentions        []common.Extention  `json:"extentions"`
	Permissions       []common.Permission `json:"permissions"`
}

func (p *Plugin) Clone() *Plugin {
	result := &Plugin{
		CustomName:  p.CustomName,
		UUID:        p.UUID,
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Email:       p.Email,
		Url:         p.Url,
		PluginType:  p.PluginType,
		Enabled:     p.Enabled,
		Permissions: p.Permissions,
	}

	es := make([]common.Extention, len(p.Extentions))
	for i, e := range p.Extentions {
		cloned := e.Clone() // 调用每个成员的 Clone 方法，假设实现了 Clone 方法
		es[i] = cloned
	}
	result.Extentions = es

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
		err := m.updatePlugin(p.UUID, &PluginParam{CustomName: p.CustomName, Url: p.Url}, p.Enabled)
		if err != nil {
			logger.Error("failed to update plugin %s info", p.Name)
			// 插件离线等异常状况
			// 使用历史数据
			m.Lock()
			m.Plugins = append(m.Plugins, &Plugin{
				CustomName:  p.CustomName,
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
func (m *PluginManager) updatePlugin(uuid string, pp *PluginParam, enabled int) error {
	// 查询最新的插件信息
	logger.Debug("update plugin")
	info, err := requestPluginInfo(pp)
	if err != nil {
		logger.Error("failed to request plugin info:%s", err.Error())
		return err
	}
	info.UUID = uuid
	err = Handshake(pp.Url, info)
	if err != nil {
		return err
	}

	p := &Plugin{
		UUID:        uuid,
		CustomName:  info.CustomName,
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Email:       info.Email,
		Url:         pp.Url,
		PluginType:  info.PluginType,
		Enabled:     enabled,
		Extentions:  info.Extentions,
		Permissions: info.Permissions,
	}

	if err := dao.UpdatePluginInfo(toPluginDao(p)); err != nil {
		return err
	}
	//更新插件管理器的插件列表
	found := false
	m.Lock()
	for i, v := range m.Plugins {
		if v.UUID == uuid {
			m.Plugins[i] = p
			found = true
			break
		}
	}
	if !found {
		m.Plugins = append(m.Plugins, p)
	}
	m.Unlock()

	return nil
}

func (m *PluginManager) addPlugin(addPlugin *PluginParam) (*Plugin, error) {
	p, err := requestPluginInfo(addPlugin)
	if err != nil {
		return nil, err
	}

	if p.UUID == "" {
		p.UUID = uuid.New().String()
	}

	err = Handshake(addPlugin.Url, p)
	if err != nil {
		return nil, err
	}

	// use accessible url as plugin url
	p.Url = addPlugin.Url
	if err := dao.RecordPlugin(toPluginDao(p)); err != nil {
		return nil, err
	}

	m.Lock()
	m.Plugins = append(m.Plugins, p)
	m.Unlock()

	return p, nil
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
	custom_name := ""
	m.Lock()
	for _, v := range m.Plugins {
		if v.UUID == uuid {
			v.Enabled = enable
			url = v.Url
			custom_name = v.CustomName
			break
		}
	}
	m.Unlock()
	if url == "" {
		logger.Error("get plugin url error")
		return errors.New("get plugin url error")
	}
	// 开启插件的时候更新最新的插件信息
	if enable == 1 {
		err := m.updatePlugin(uuid, &PluginParam{CustomName: custom_name, Url: url}, enable)
		if err != nil {
			logger.Error("failed to update plugin info:%s", err.Error())
			return err
		}
	}
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

func (m *PluginManager) GetPluginByUUID(uuid string) (*Plugin, error) {
	var result *Plugin
	found := false
	m.Lock()
	for _, v := range m.Plugins {
		if v.UUID == uuid {
			// 使用深拷贝避免指针泄露
			result = v.Clone()
			found = true
			break
		}
	}
	m.Unlock()

	if !found {
		return nil, fmt.Errorf("plugin %s not found", uuid)
	}
	return result, nil
}

// 获取所有插件
func (m *PluginManager) getPlugins() ([]*Plugin, error) {
	result := []*Plugin{}
	m.Lock()
	for _, v := range m.Plugins {
		// 使用深拷贝避免指针泄露
		p := v.Clone()
		result = append(result, p)
	}
	m.Unlock()

	return result, nil
}

func toPluginDao(p *Plugin) *dao.PluginModel {
	return &dao.PluginModel{
		CustomName:  p.CustomName,
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
func Handshake(url string, p *Plugin) error {
	index := strings.Index(url, "plugin")
	if index > 0 {
		url = url[:index]
	}
	port := strings.Split(config.OptionsConfig.HttpServer.Addr, ":")[1]

	url = strings.TrimRight(url, "/") + "/plugin_manage/bind?port=" + port
	logger.Debug("plugin url is:%s", url)

	token, err := jwt.GeneratePluginToken(p.Name, p.UUID)
	if err != nil {
		logger.Error("generate plugin token error:%s", err.Error())
		return err
	}

	resp, err := httputils.Put(url, &httputils.Params{
		Cookie: map[string]string{
			"PluginToken": token,
		},
	})
	if err != nil {
		logger.Error("request plugin info error:%s", err.Error())
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
func requestPluginInfo(plugin *PluginParam) (*Plugin, error) {
	url := plugin.Url
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

	//先解析插件信息
	PluginInfo := &client.PluginInfo{}
	err = json.Unmarshal(resp.Body, PluginInfo)
	if err != nil {
		logger.Error("unmarshal request plugin info error:%s", err.Error())
		return nil, err
	}

	//解析插件扩展点信息
	data := struct {
		Extentions []map[string]interface{} `json:"extentions"`
	}{}
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		logger.Error("request plugin extentions error:%s", err.Error())
		return nil, err
	}
	extentions := common.ParseParameters(data.Extentions)

	permissions := struct {
		Permissions []common.Permission `json:"permissions"`
	}{}
	if err := json.Unmarshal(resp.Body, &permissions); err != nil {
		logger.Error("request plugin permissions error:%s", err.Error())
		return nil, err
	}

	// TODO: check info valid
	return &Plugin{
		CustomName:  plugin.CustomName,
		Name:        PluginInfo.Name,
		Version:     PluginInfo.Version,
		Description: PluginInfo.Description,
		Author:      PluginInfo.Author,
		Email:       PluginInfo.Email,
		Url:         url,
		PluginType:  PluginInfo.PluginType,
		Extentions:  extentions,
		Permissions: permissions.Permissions,
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

// 获取单个plugin
func GetPluginByUUID(uuid string) (*Plugin, error) {
	plugin, err := globalPluginManager.GetPluginByUUID(uuid)
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
	result := []*Plugin{}
	for _, p := range plugins {
		plugin_status, err := GetPluginConnectStatus(p.Url)
		if err != nil {
			logger.Error("plugin status get failed %s, %s", err, p.Url)
			continue
		}
		p.ConnectStatus = plugin_status.Connected
		p.LastHeartbeatTime = plugin_status.LastConnect.Format("2006-01-02 15:04:05")
		result = append(result, p)
	}
	return result, nil
}

// 分页查询
func GetPluginPaged(offset, size int) (int64, []*Plugin, error) {
	// 借助db的分页功能实现分页查询
	total, plugins, error := dao.GetPluginPaged(offset, size)

	result := []*Plugin{}
	for _, p := range plugins {
		plugin, err := globalPluginManager.GetPluginByUUID(p.UUID)
		if err != nil {
			logger.Error("manager get plugin %s", err)
			continue
		}
		plugin_status, err := GetPluginConnectStatus(p.Url)
		if err != nil {
			logger.Error("plugin status get failed %s", err)
			continue
		}
		plugin.ConnectStatus = plugin_status.Connected
		plugin.LastHeartbeatTime = plugin_status.LastConnect.Format("2006-01-02 15:04:05")
		result = append(result, plugin)
	}

	return total, result, error
}

type PluginParam struct {
	CustomName string `json:"custom_name"`
	Url        string `json:"url"`
}

func AddPlugin(param *PluginParam) (*Plugin, error) {
	url := param.Url
	logger.Debug("add plugin from %s", url)

	ok, err := dao.IsExistCustomName(param.CustomName)
	if err != nil {
		return &Plugin{}, err
	}
	if ok {
		return &Plugin{}, errors.New("已存在相同插件名称")
	}

	p, err := globalPluginManager.addPlugin(param)
	if err != nil {
		return p, err
	}

	if err := addPluginInRedis(p.UUID); err != nil {
		logger.Error("failed to add plugin heartbeat in redis:%s", err.Error())
		return p, err
	}

	//向数据库添加admin用户的插件权限
	err = auth.AddPluginPermission("admin", p.Permissions, p.UUID)
	return p, err
}

func DeletePlugin(uuid string, plugin *Plugin) error {
	logger.Debug("delete plugin: %s", uuid)

	//删除插件权限
	logger.Debug("delete %s plugin permission", plugin.Name)

	err := auth.DeletePluginPermission(plugin.Permissions, plugin.UUID)
	if err != nil {
		logger.Error("failed to delete plugin permissions in mysql:%s", err.Error())
		return err
	}

	if err := deletePluginInRedis(uuid); err != nil {
		logger.Error("failed to delete plugin heartbeat in redis:%s", err.Error())
		return err
	}

	if err := globalPluginManager.deletePlugin(uuid); err != nil {
		logger.Error("failed to delete plugin info:%s", err.Error())
		return err
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

func deletePluginInRedis(uuid string) error {
	p, err := dao.QueryPluginById(uuid)
	if err != nil {
		return err
	}
	logger.Debug("delete %v plugin heartbeat in redis: %v", p.Name, p.Url)
	key := client.HeartbeatKey + p.Url
	err = redismanager.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func addPluginInRedis(uuid string) error {
	p, err := dao.QueryPluginById(uuid)
	if err != nil {
		return err
	}
	logger.Debug("add %v plugin heartbeat in redis: %v", p.Name, p.Url)
	key := client.HeartbeatKey + p.Url
	value := client.PluginStatus{
		Connected:   true,
		LastConnect: time.Now(),
	}
	err = redismanager.Set(key, value)
	return err
}

func GetPluginConnectStatus(url string) (*client.PluginStatus, error) {
	key := client.HeartbeatKey + url
	var valueObj client.PluginStatus
	plugin_status, err := redismanager.Get(key, &valueObj)
	if err != nil {
		return nil, err
	}
	return plugin_status.(*client.PluginStatus), nil
}

// 从db获取某个角色的所有插件权限
func GetRolePluginPermission(role string) map[string]interface{} {
	p2p := make(map[string]interface{})
	for _, v := range globalPluginManager.Plugins {
		pers := []string{}
		policys := auth.GetFilteredPolicy(role, "", "", v.UUID)
		for _, p := range policys {
			pers = append(pers, p[1]+"/"+p[2])
		}
		p2p[v.UUID] = struct {
			UUID        string
			Name        string
			Permassions []string
		}{
			UUID:        v.UUID,
			Name:        v.Name,
			Permassions: pers,
		}
	}
	return p2p
}

type PluginPermission struct {
	UUID        string              `json:"uuid"`
	Permissions []common.Permission `json:"permissions"`
}

// 更新插件角色权限
func UpdatePluginPermissions(role string, PluginPermissions []PluginPermission) error {
	for _, p := range PluginPermissions {
		err := auth.AddPluginPermission(role, p.Permissions, p.UUID)
		if err != nil {
			logger.Error("add role:%s buttion policy failed:%s", role, err)
			return err
		}
	}
	return nil
}
