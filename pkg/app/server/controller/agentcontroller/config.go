package agentcontroller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/global"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/sdk/response"
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
	ConMess.Machine_uuid = uuid
	ConMess.ConfigName = "/home/wbj/PilotGo/config_server.yaml.templete"
	err := agent.ConfigfileInfo(ConMess)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "Success")
}
