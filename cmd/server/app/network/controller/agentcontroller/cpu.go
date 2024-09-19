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
 * Date: 2022-02-16 15:13:25
 * LastEditTime: 2022-04-13 15:51:43
 * Description: provide agent cpu manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func CPUInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		response.Fail(c, nil, "uuid参数缺失")
		return
	}
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	cpuInfo, err := agent.GetCPUInfo()
	if err != nil {
		response.Fail(c, nil, "获取系统CPU信息失败!")
		return
	}
	response.Success(c, gin.H{"CPU_info": cpuInfo}, "Success")
}
