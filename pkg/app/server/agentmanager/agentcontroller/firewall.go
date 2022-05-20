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
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

type ZonePort struct {
	UUID string `json:"uuid"`
	Zone string `json:"zone"`
	Port string `json:"port"`
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

	add, Err, err := agent.FirewalldZonePortAdd(zp.Zone, zp.Port)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "添加失败!")
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

	del, Err, err := agent.FirewalldZonePortDel(zp.Zone, zp.Port)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "删除失败!")
		return
	}
	response.Success(c, gin.H{"firewalld_del": del}, "删除成功!")
}
