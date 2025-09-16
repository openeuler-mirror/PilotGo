package pluginapi

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type ScriptsRun struct {
	Batch         *common.Batch `json:"batch"`
	ScriptType    string        `json:"script_type"`
	ScriptContent string        `json:"script_content"`
	Params        string        `json:"params"`
	TimeOutSec    int           `json:"timeoutSec"`
}

// 远程agent运行各种脚本
func AgentRunScriptsHandler(c *gin.Context) {
	logger.Debug("process get agent request")

	sr := &ScriptsRun{}
	err := c.ShouldBind(sr)
	if err != nil {
		logger.Debug("解析脚本基本信息失败:%s", err)
		response.Fail(c, nil, "参数失败")
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			data, err := agent.AgentRunScripts(sr.ScriptType, sr.ScriptContent, sr.Params, sr.TimeOutSec)
			if err != nil {
				logger.Error("run script error, agent:%s, command:%s", uuid, sr.ScriptContent)
			}
			logger.Debug("run script on agent result:%v", data)
			re := common.CmdResult{
				MachineUUID: uuid,
				MachineIP:   agent.IP,
				RetCode:     data.RetCode,
				Stdout:      data.Stdout,
				Stderr:      data.Stderr,
			}
			return re
		}
		return common.CmdResult{}
	}

	result := batch.BatchProcess(sr.Batch, f, sr.ScriptContent, sr.Params)
	response.Success(c, result, "run script succeed")
}
