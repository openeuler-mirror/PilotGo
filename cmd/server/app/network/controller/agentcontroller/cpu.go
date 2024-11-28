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
