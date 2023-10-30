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
 * LastEditTime: 2022-04-11 17:07:55
 * Description: provide agent service manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
)

type AgentService struct {
	UUID         string `json:"uuid"`
	Service      string `json:"service"`
	UserName     string `json:"userName"`
	UserDeptName string `json:"userDept"`
}

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

	serviceInfo, err := agent.GetService(service)
	if err != nil {
		response.Fail(c, nil, "获取服务状态失败!")
		return
	}
	response.Success(c, gin.H{"service_status": serviceInfo.ServiceActiveStatus}, "Success")
}

func ServiceStartHandler(c *gin.Context) {
	var agentservice AgentService
	if err := c.Bind(&agentservice); err != nil {
		response.Fail(c, nil, "parameter error")
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
		Module:     auditlog.ModuleMachine,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "ServiceStart",
	}
	auditlog.Add(log)

	agent := agentmanager.GetAgent(agentservice.UUID)
	if agent == nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+"获取uuid失败")
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "获取uuid失败")
		return
	}

	service_start, Err, err := agent.ServiceStart(agentservice.Service)
	if len(Err) != 0 || err != nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+err.Error())
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"error": Err}, "Failed!")
		return
	}
	response.Success(c, gin.H{"service_start": service_start}, "Success")
}

func ServiceStopHandler(c *gin.Context) {
	var agentservice AgentService
	if err := c.Bind(&agentservice); err != nil {
		response.Fail(c, nil, "parameter error")
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
		Module:     auditlog.ModuleMachine,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "ServiceStop",
	}
	auditlog.Add(log)

	agent := agentmanager.GetAgent(agentservice.UUID)
	if agent == nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+"获取uuid失败")
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "获取uuid失败")
		return
	}

	service_stop, Err, err := agent.ServiceStop(agentservice.Service)
	if len(Err) != 0 || err != nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+err.Error())
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		return
	}
	response.Success(c, gin.H{"service_stop": service_stop}, "Success")
}

func ServiceRestartHandler(c *gin.Context) {
	var agentservice AgentService
	if err := c.Bind(&agentservice); err != nil {
		response.Fail(c, nil, "parameter error")
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
		Module:     auditlog.ModuleMachine,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "ServiceRestart",
	}
	auditlog.Add(log)

	agent := agentmanager.GetAgent(agentservice.UUID)
	if agent == nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+"获取uuid失败")
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "获取uuid失败")
		return
	}

	service_restart, Err, err := agent.ServiceRestart(agentservice.Service)
	if len(Err) != 0 || err != nil {
		auditlog.UpdateMessage(log, "agentuuid:"+agentservice.UUID+err.Error())
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"error": Err}, "Failed!")
		return
	}
	response.Success(c, gin.H{"service_restart": service_restart}, "Success")
}
