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
 * Date: 2022-02-16 09:28:46
 * LastEditTime: 2022-03-01 11:18:35
 * Description: provide agent service manager functions.
 ******************************************************************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func ServiceListHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_list, err := agent.ServiceList()
	if err != nil {
		response.Fail(c, nil, "获取服务列表失败!")
		return
	}
	response.Success(c, gin.H{"service_list": service_list}, "Success")
}
func ServiceStatusHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_status, err := agent.ServiceStatus(service)
	if err != nil {
		response.Fail(c, nil, "获取服务状态失败!")
		return
	}
	response.Success(c, gin.H{"service_status": service_status}, "Success")
}

func ServiceStartHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "运行服务"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = service
	logParent.Action = model.ServiceStart

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_start, Err, err := agent.ServiceStart(service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "启动服务失败!")
		logParent.StatusCode = 400
		logParent.Message = Err
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"service_start": service_start}, "Success")
	logParent.StatusCode = 200
	logParent.Message = "启动服务成功"
	mysqlmanager.DB.Save(&logParent)
}
func ServiceStopHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "运行服务"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = service
	logParent.Action = model.ServiceStop

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_stop, Err, err := agent.ServiceStop(service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "关闭服务失败!")
		logParent.StatusCode = 400
		logParent.Message = Err
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"service_stop": service_stop}, "Success")
	logParent.StatusCode = 200
	logParent.Message = "关闭服务成功"
	mysqlmanager.DB.Save(&logParent)
}
func ServiceRestartHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "运行服务"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = service
	logParent.Action = model.ServiceRestart

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_restart, Err, err := agent.ServiceRestart(service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "重启服务失败!")
		logParent.StatusCode = 400
		logParent.Message = Err
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"service_restart": service_restart}, "Success")
	logParent.StatusCode = 200
	logParent.Message = "重启服务成功"
	mysqlmanager.DB.Save(&logParent)
}
