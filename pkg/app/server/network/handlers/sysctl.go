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
 * LastEditTime: 2022-02-28 15:16:13
 * Description: provide Kernel configuration.
 ******************************************************************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
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
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "配置内核参数"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = args
	logParent.Action = model.SysctlChange

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	sysctl_change, err := agent.ChangeSysctl(args)
	if err != nil {
		response.Fail(c, nil, "修改内核运行时参数失败!")
		logParent.StatusCode = 400
		logParent.Message = "修改内核运行时参数失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"sysctl_change": sysctl_change}, "Success")
	logParent.StatusCode = 200
	logParent.Message = "修改成功"
	mysqlmanager.DB.Save(&logParent)
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
