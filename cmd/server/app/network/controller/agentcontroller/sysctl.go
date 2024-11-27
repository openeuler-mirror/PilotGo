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

func SysInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	sysctl_info, err := agent.GetSysctlInfo()
	if err != nil {
		response.Fail(c, nil, "获取内核配置失败!")
		return
	}
	response.Success(c, gin.H{"sysctl_info": sysctl_info}, "Success")
}

func SysctlChangeHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	args := c.Query("args")
	//username := c.Query("userName")
	//userDeptName := c.Query("userDept")

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID: uuidservice.New().String(),
		Module:  auditlog.ModuleMachine,
		Status:  auditlog.StatusOK,
		UserID:  u.ID,
		Action:  "SysctlChange",
		Message: "agentuuid:" + uuid,
	}
	auditlog.Add(log)

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		auditlog.UpdateMessage(log, "agentuuid:"+uuid+"获取uuid失败")
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "获取uuid失败")
		return
	}

	sysctl_change, err := agent.ChangeSysctl(args)
	if err != nil {
		auditlog.UpdateMessage(log, "agentuuid:"+uuid+"修改内核运行时参数失败")
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"error": err}, "修改内核运行时参数失败!")
		return
	}
	response.Success(c, gin.H{"sysctl_change": sysctl_change}, "Success")
}

func SysctlViewHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	args := c.Query("args")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	sysctl_view, err := agent.SysctlView(args)
	if err != nil {
		response.Fail(c, nil, "获取该参数的值失败!")
		return
	}
	response.Success(c, gin.H{"sysctl_view": sysctl_view}, "Success")
}
