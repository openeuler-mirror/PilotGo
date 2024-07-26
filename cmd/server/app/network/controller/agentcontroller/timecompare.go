package agentcontroller

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
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
	resulttime, err := strconv.ParseInt(strings.Replace(fmt.Sprint(result), "\n", "", -1), 10, 64)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "数据有误")
		return
	}
	if math.Abs(float64(currentTime-resulttime)) < 100 {
		response.Success(c, nil, "Success")
		return
	}
	response.Fail(c, nil, "时间不一致")
}
