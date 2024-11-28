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
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
)

type AgentService struct {
	UUID    string `json:"uuid"`
	Service string `json:"service"`
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
