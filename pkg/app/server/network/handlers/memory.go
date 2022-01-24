package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func MemoryInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	memory_info, err := agent.GetMemoryInfo()
	if err != nil {
		response.Fail(c, nil, "获取系统内存信息失败!")
		return
	}
	response.Success(c, gin.H{"memory_info": memory_info}, "Success")
}
