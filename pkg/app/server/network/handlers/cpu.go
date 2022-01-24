package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func CPUInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	cpu_info, err := agent.GetCPUInfo()
	if err != nil {
		response.Fail(c, nil, "获取系统CPU信息失败!")
		return
	}
	response.Success(c, gin.H{"CPU_info": cpu_info}, "Success")
}
