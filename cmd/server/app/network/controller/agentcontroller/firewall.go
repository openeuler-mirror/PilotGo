/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type ZonePort struct {
	UUID     string `json:"uuid"`
	Zone     string `json:"zone"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Source   string `json:"source"`
}

func FirewalldConfig(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	config, Err, err := agent.FirewalldConfig()
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "获取防火墙配置失败!")
		return
	}
	response.Success(c, gin.H{"firewalld_config": config}, "获取防火墙配置成功!")
}

func FirewalldZoneConfig(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	zone := c.Query("zone")
	config, Err, err := agent.FirewalldZoneConfig(zone)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "获取防火墙区域配置失败!")
		return
	}
	response.Success(c, gin.H{"firewalld_zone": config}, "获取防火墙区域配置成功!")
}

func FirewalldSetDefaultZone(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	_, Err, err := agent.FirewalldSetDefaultZone(zp.Zone)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "更改防火墙默认区域失败")
		return
	}
	response.Success(c, nil, "更改防火墙默认区域成功!")
}
func FirewalldServiceAdd(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	Err, err := agent.FirewalldServiceAdd(zp.Zone, zp.Service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, nil, "添加防火墙服务成功!")
}

func FirewalldServiceRemove(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	Err, err := agent.FirewalldServiceRemove(zp.Zone, zp.Service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, nil, "移除防火墙服务成功!")
}

func FirewalldSourceAdd(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	Err, err := agent.FirewalldSourceAdd(zp.Zone, zp.Source)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, nil, "添加防火墙来源地址成功!")
}

func FirewalldSourceRemove(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	Err, err := agent.FirewalldSourceRemove(zp.Zone, zp.Source)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, nil, "移除防火墙来源地址成功!")
}

func FirewalldRestart(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	restart, Err, err := agent.FirewalldRestart()
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "重启防火墙失败")
		return
	}
	response.Success(c, gin.H{"firewalld_restart": restart}, "重启防火墙成功!")
}

func FirewalldStop(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	stop, Err, err := agent.FirewalldStop()
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "关闭防火墙失败!")
		return
	}
	response.Success(c, gin.H{"firewalld_stop": stop}, "关闭防火墙成功!")
}

func FirewalldZonePortAdd(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)
	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	add, Err, err := agent.FirewalldZonePortAdd(zp.Zone, zp.Port, zp.Protocol)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, gin.H{"firewalld_add": add}, "添加成功!")
}

func FirewalldZonePortDel(c *gin.Context) {
	var zp ZonePort
	c.ShouldBind(&zp)

	agent := agentmanager.GetAgent(zp.UUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	del, Err, err := agent.FirewalldZonePortDel(zp.Zone, zp.Port, zp.Protocol)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, gin.H{"firewalld_del": del}, "删除成功!")
}
