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
 * LastEditTime: 2022-04-20 15:51:51
 * Description: provide agent run script functions.
 ******************************************************************************/
package agentcontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

func RunScript(c *gin.Context) {
	logger.Debug("process get agent request")
	// TODO: process agent info
	uuid := c.Query("uuid")
	cmd := c.Query("cmd")
	fmt.Println(uuid, cmd)

	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunCommand(cmd)
		if err != nil {
			logger.Error("run script error, agent:%s, cmd:%s", uuid, cmd)
			c.JSON(http.StatusOK, `{"status":-1}`)
		}
		logger.Info("run command on agent result:%v", data)
		c.JSON(http.StatusOK, `{"status":0}`)
		return
	}

	logger.Info("unknown agent:%s", uuid)
	c.JSON(http.StatusOK, `{"status":-1}`)
}
