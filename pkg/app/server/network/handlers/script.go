package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func RunScript(c *gin.Context) {
	logger.Debug("process get agent request")
	// TODO: process agent info
	uuid := c.Query("uuid")
	cmd := c.Query("cmd")
	fmt.Println(uuid, cmd)

	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunScript(cmd)
		if err != nil {
			logger.Error("run script error, agent:%s, cmd:%s", uuid, cmd)
			c.JSON(http.StatusOK, `{"status":-1}`)
		}
		logger.Info("run command on agent result:%v", data)
		c.JSON(http.StatusOK, `{"status":0}`)
		return
	}

	logger.Info("unknown agent:%s", uuid)
	c.JSON(http.StatusOK, `{"status":-1}`)
}
