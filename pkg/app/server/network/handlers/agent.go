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
	agent := agentmanager.GetAgentManager().GetAgent("uuid")
	if agent != nil {
		agent.GetInfo()
		// TODO: 此处处理并返回agent信息
	}

	c.JSON(http.StatusOK, `{"status":-1}`)
}
