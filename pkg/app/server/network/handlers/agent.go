package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func AgentInfoHandler(c *gin.Context) {
	logger.Debug("process get agent request")
	// TODO: process agent info
	agent := agentmanager.GetAgent("uuid")
	if agent != nil {
		agent.AgentInfo()
		// TODO: 此处处理并返回agent信息

		c.JSON(http.StatusOK, `{"status":0}`)
		return
	}

	c.JSON(http.StatusOK, `{"status":-1}`)
}

func AgentListHandler(c *gin.Context) {
	logger.Debug("process get agent list request")

	agent_list := agentmanager.GetAgentList()

	c.JSON(http.StatusOK, agent_list)
}
