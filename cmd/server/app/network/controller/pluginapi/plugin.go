/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package pluginapi

import (
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginList(c *gin.Context) {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		response.Fail(c, nil, "查询插件错误："+err.Error())
		return
	}
	response.Success(c, plugins, "插件查询成功")
}

func PluginHeartbeat(c *gin.Context) {
	ClientID := &struct {
		ClientID string `json:"clientID"`
	}{}
	if err := c.ShouldBindJSON(&ClientID); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 更新心跳时间
	key := client.HeartbeatKey + ClientID.ClientID
	value := client.PluginStatus{
		Connected:   true,
		LastConnect: time.Now(),
	}
	err := redismanager.Set(key, value)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "Heartbeat received")
}

func HasPermission(c *gin.Context) {
	p := &common.Permission{}
	if err := c.ShouldBind(&p); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	//TODO:解析发送请求插件的uuid
	var uuid string
	ok, err := auth.CheckAuth(user.Username, p.Resource, p.Operate, uuid)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, ok, "include permission")
}
