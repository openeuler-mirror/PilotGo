/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-02-17 02:43:29
 * LastEditTime: 2022-04-13 01:51:51
 * Description: provide agent firewall manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"gitee.com/PilotGo/PilotGo/app/server/agentmanager"
	"gitee.com/PilotGo/PilotGo/sdk/response"
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
