/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	eventSDK "gitee.com/openeuler/PilotGo-plugins/event/sdk"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/common"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	commonSDK "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
)

// 查询插件清单
func GetPluginsHandler(c *gin.Context) {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		response.Fail(c, nil, "查询插件错误："+err.Error())
		return
	}

	logger.Info("find %d plugins", len(plugins))
	response.Success(c, plugins, "插件查询成功")
}

// 分页查询插件清单
func GetPluginsPagedHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := plugin.GetPluginPaged(num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 添加插件
func AddPluginHandler(c *gin.Context) {
	param := plugin.PluginParam{}
	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Add Plugin",
	}
	auditlog.Add(log)

	p, err := plugin.AddPlugin(&param)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "add plugin failed:"+err.Error())
		return
	}

	msgData := commonSDK.MessageData{
		MsgType:     eventSDK.MsgPluginAdd,
		MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgPluginAdd),
		TimeStamp:   time.Now(),
		Data: eventSDK.MDPluginChange{
			PluginName:  p.Name,
			Version:     p.Version,
			Url:         p.Url,
			Description: p.Description,
		},
	}
	msgDataString, err := msgData.ToMessageDataString()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	ms := commonSDK.EventMessage{
		MessageType: eventSDK.MsgPluginAdd,
		MessageData: msgDataString,
	}
	plugin.PublishEvent(ms)

	response.Success(c, nil, "插件添加成功")
}

// 停用/启动插件
func TogglePluginHandler(c *gin.Context) {
	param := struct {
		UUID   string `json:"uuid"`
		Enable int    `json:"enable"`
	}{}

	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Toggle Plugin",
	}
	auditlog.Add(log)

	logger.Info("toggle plugin:%s to enable %d", param.UUID, param.Enable)
	if err = plugin.TogglePlugin(param.UUID, param.Enable); err != nil {
		response.Fail(c, nil, "toggle plugin error:"+err.Error())
		return
	}
	//根据uuid从pluginmanage中获取本插件信息，并返回插件扩展点信息
	plugin, err := plugin.GetPluginByUUID(param.UUID)
	if err != nil {
		response.Fail(c, nil, "get plugin by uuid error:"+err.Error())
	}
	response.Success(c, plugin.Extentions, "插件信息更新成功")
}

// 卸载插件
func UnloadPluginHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "undefined" {
		response.Fail(c, nil, "参数错误")
		return
	}

	p, err := plugin.GetPluginByUUID(uuid) //获取插件信息
	if err != nil {
		logger.Error("get plugin by uuid error:%s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Unload Plugin",
	}
	auditlog.Add(log)

	msgData := commonSDK.MessageData{
		MsgType:     eventSDK.MsgPluginRemove,
		MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgPluginRemove),
		TimeStamp:   time.Now(),
		Data: eventSDK.MDPluginChange{
			PluginName:  p.Name,
			Version:     p.Version,
			Url:         p.Url,
			Description: p.Description,
		},
	}
	msgDataString, err := msgData.ToMessageDataString()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	ms := commonSDK.EventMessage{
		MessageType: eventSDK.MsgPluginRemove,
		MessageData: msgDataString,
	}
	plugin.PublishEvent(ms)

	logger.Info("unload plugin:%s", uuid)
	if err := plugin.DeletePlugin(uuid, p); err != nil {
		return
	}
	response.Success(c, nil, "插件信息更新成功")
}

func PluginGatewayHandler(c *gin.Context) {
	// TODO
	name := c.Param("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		c.String(http.StatusNotFound, "plugin not found: "+err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "parse Plugin",
	}
	auditlog.Add(log)

	s := strings.Replace(p.Url, "/plugin/"+name, "", 1)
	ishttp, err := httputils.ServerIsHttp(s)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	if ishttp && strings.Split(s, "://")[0] == "https" {
		s = "http://" + strings.Split(s, "://")[1]
	}
	if !ishttp && strings.Split(s, "://")[0] == "http" {
		s = "https://" + strings.Split(s, "://")[1]
	}

	target, err := url.Parse(s)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	logger.Debug("proxy plugin request to: %s", target)
	c.Request.Host = target.Host

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
