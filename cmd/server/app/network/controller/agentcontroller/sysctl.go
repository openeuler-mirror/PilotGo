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
 * LastEditTime: 2022-03-25 02:00:34
 * Description: provide Kernel configuration.
 ******************************************************************************/
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
