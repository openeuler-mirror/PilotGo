package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

func ServiceListHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_list, err := agent.ServiceList()
	if err != nil {
		response.Fail(c, nil, "获取服务列表失败!")
		return
	}
	response.Success(c, gin.H{"service_list": service_list}, "Success")
}
func ServiceStatusHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_status, err := agent.ServiceStatus(service)
	if err != nil {
		response.Fail(c, nil, "获取服务状态失败!")
		return
	}
	response.Success(c, gin.H{"service_status": service_status}, "Success")
}

func ServiceStartHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_start, err := agent.ServiceStart(service)
	if err != nil {
		response.Fail(c, nil, "启动服务失败!")
		return
	}
	response.Success(c, gin.H{"service_start": service_start}, "Success")
}
func ServiceStopHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_stop, err := agent.ServiceStop(service)
	if err != nil {
		response.Fail(c, nil, "关闭服务失败!")
		return
	}
	response.Success(c, gin.H{"service_stop": service_stop}, "Success")
}
func ServiceRestartHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_restart, err := agent.ServiceRestart(service)
	if err != nil {
		response.Fail(c, nil, "重启服务失败!")
		return
	}
	response.Success(c, gin.H{"service_restart": service_restart}, "Success")
}
