package agentcontroller

import (
	"time"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func CompareTimeHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}
	result, err := agent.GetTimeInfo()
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "时间获取失败")
		return
	}
	currentTime := time.Now().Unix()
	if currentTime-result.(int64) < 150 {
		response.Success(c, nil, "Success")
		return
	}
	response.Fail(c, nil, "时间不一致")
}
