package pluginapi

import (
	"github.com/gin-gonic/gin"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func Service(ctx *gin.Context) {
	// TODO: support batch
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

func StartService(c *gin.Context) {
	// TODO: support batch
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_start, Err, err := agent.ServiceStart(service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "Failed!")
		return
	}

	response.Success(c, gin.H{"service_start": service_start}, "Success")
}

func StopService(c *gin.Context) {
	// TODO: support batch
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_stop, Err, err := agent.ServiceStop(service)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "Failed!")
		return
	}

	response.Success(c, gin.H{"service_stop": service_stop}, "Success")

}
