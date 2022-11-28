package agentcontroller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/global"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func ConfigfileHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}
	//传入需要监控的文件的信息
	var ConMess global.ConfigMessage
	ConMess.ConfigName = "config_server"
	ConMess.ConfigType = "yaml"
	ConMess.ConfigPath = "."
	err := agent.ConfigfileInfo(ConMess)
	if err != nil {
		response.Fail(c, nil, "配置文件监控失败!")
		return
	}
	response.Success(c, nil, "Success")
}
