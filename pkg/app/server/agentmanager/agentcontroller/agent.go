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
 * Date: 2022-01-24 15:08:08
 * LastEditTime: 2022-04-11 14:03:11
 * Description: Get the basic information of the machine
 ******************************************************************************/
package agentcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func AgentInfoHandler(c *gin.Context) {
	logger.Debug("process get agent request")
	// TODO: process agent info
	agent := agentmanager.GetAgent("uuid")
	if agent == nil {
		c.JSON(http.StatusOK, `{"status":-1}`)
	}

	agent.AgentInfo()
	// TODO: 此处处理并返回agent信息

	c.JSON(http.StatusOK, `{"status":0}`)
}

func AgentListHandler(c *gin.Context) {
	logger.Debug("process get agent list request")

	agent_list := agentmanager.GetAgentList()

	c.JSON(http.StatusOK, agent_list)
}

func OsBasic(c *gin.Context) {
	uuid := c.Query("uuid")
	ip, state, dept := dao.MachineBasic(uuid)
	response.Response(c, http.StatusOK, http.StatusOK, gin.H{"IP": ip, "state": state, "depart": dept}, "Success")
}
