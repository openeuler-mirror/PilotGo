package pluginapi

import (
	"github.com/gin-gonic/gin"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func Service(ctx *gin.Context) {
	uuid := ctx.Query("uuid")
	service := ctx.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(ctx, nil, "获取uuid失败!")
		return
	}

	service_status, err := agent.ServiceStatus(service)
	if err != nil {
		response.Fail(ctx, nil, "获取服务状态失败!")
		return
	}
	response.Success(ctx, gin.H{"service_status": service_status}, "Success")
}

func StartService(c *gin.Context) {}

func StopService(c *gin.Context) {}
