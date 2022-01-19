package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
)

func OSInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		c.JSON(http.StatusOK, `{"status":-1}`)
	}

	os_info, err := agent.GetOSInfo()
	if err != nil {
		c.JSON(http.StatusOK, `{"status":-1}`)
	}

	c.JSON(http.StatusOK, os_info)
}
