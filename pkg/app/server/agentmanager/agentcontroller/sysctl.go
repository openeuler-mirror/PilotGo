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
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
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
	userDeptName := c.Query("userDept")

	logParent := model.AgentLogParent{
		UserName:   username,
		DepartName: userDeptName,
		Type:       global.LogTypeSysctl,
	}
	logParentId := dao.ParentAgentLog(logParent)

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {

		log := model.AgentLog{
			LogParentID:     logParentId,
			IP:              dao.UUID2MacIP(uuid),
			OperationObject: args,
			Action:          global.SysctlChange,
			StatusCode:      http.StatusBadRequest,
			Message:         "获取uuid失败",
		}
		dao.AgentLog(log)
		response.Fail(c, nil, "获取uuid失败")

		dao.UpdateParentAgentLog(logParentId, global.ActionFalse)
		return
	}

	sysctl_change, err := agent.ChangeSysctl(args)
	if err != nil {
		log := model.AgentLog{
			LogParentID:     logParentId,
			IP:              dao.UUID2MacIP(uuid),
			OperationObject: args,
			Action:          global.SysctlChange,
			StatusCode:      http.StatusBadRequest,
			Message:         err.Error(),
		}
		dao.AgentLog(log)
		response.Fail(c, gin.H{"error": err}, "修改内核运行时参数失败!")

		dao.UpdateParentAgentLog(logParentId, global.ActionFalse)
		return
	}
	log := model.AgentLog{
		LogParentID:     logParentId,
		IP:              dao.UUID2MacIP(uuid),
		OperationObject: args,
		Action:          global.SysctlChange,
		StatusCode:      http.StatusOK,
		Message:         "修改成功",
	}
	dao.AgentLog(log)
	dao.UpdateParentAgentLog(logParentId, global.ActionOK)

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
