// 插件系统对外提供的api

package pluginapi

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 检查plugin接口调用权限
func AuthCheck(c *gin.Context) {
	// TODO
	c.Next()
}

// 远程运行脚本
func RunScriptHandler(c *gin.Context) {
	logger.Debug("process get agent request")
	uuid := c.Query("uuid")
	script := c.Query("script")

	// TODO: support batch
	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunScript(script)
		if err != nil {
			logger.Error("run script error, agent:%s, script:%s", uuid, script)
			response.Fail(c, nil, err.Error())
		}
		logger.Debug("run script on agent result:%v", data)
		response.Success(c, data, "")
		return
	}

	logger.Warn("unknown agent:%s", uuid)
	response.Fail(c, nil, "unknown agent")
}
