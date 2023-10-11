// 插件系统对外提供的api

package pluginapi

import (
	"time"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"gitee.com/openeuler/PilotGo/utils"
	"github.com/gin-gonic/gin"
)

// 检查plugin接口调用权限
func AuthCheck(c *gin.Context) {
	// TODO
	c.Next()
}

type RunResult struct {
	*utils.CmdResult
	MachineUUID string
	MachineIP   string
}

// 远程运行脚本
func RunCommandHandler(c *gin.Context) {
	logger.Debug("process get agent request")

	d := &client.CmdStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	logger.Debug("run command on agents :%v", d.Batch.MachineUUIDs)

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			data, err := agent.RunCommand(d.Command)
			if err != nil {
				logger.Error("run command error, agent:%s, command:%s", uuid, d.Command)
			}
			logger.Debug("run command on agent result:%v", data)
			re := RunResult{
				CmdResult:   data,
				MachineUUID: uuid,
				MachineIP:   agent.IP,
			}
			return re
		}
		return RunResult{}
	}

	result := batch.BatchProcess(d.Batch, f, d.Command)
	response.Success(c, result, "run cmd succeed")
}

// 远程运行脚本
func RunScriptHandler(c *gin.Context) {
	logger.Debug("process get agent request")

	d := &client.ScriptStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	logger.Debug("run script on agents :%v", d.Batch.MachineUUIDs)

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			data, err := agent.RunScript(d.Script, d.Params)
			if err != nil {
				logger.Error("run script error, agent:%s, command:%s", uuid, d.Script)
			}
			logger.Debug("run script on agent result:%v", data)
			re := RunResult{
				CmdResult:   data,
				MachineUUID: uuid,
				MachineIP:   agent.IP,
			}
			return re
		}
		return RunResult{}
	}

	result := batch.BatchProcess(d.Batch, f, d.Script, d.Params)
	response.Success(c, result, "run script succeed")
}

// 异步远程运行脚本回调
func RunCommandAsyncCallbackHandler(c *gin.Context) {
	name := c.Query("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		response.Fail(c, nil, "plugin not found: %v"+err.Error())
		return
	}
	url := p.Url + "/command_result"

	d := &client.CmdStruct{}
	if err := c.ShouldBind(d); err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	taskId := time.Now().Format("20060102150405")
	for _, macuuid := range batch.GetMachineUUIDS(d.Batch) {
		uuid := macuuid
		go runCommandCallback(uuid, d.Command, taskId, func(result client.TaskCmdResult) {
			_, err := httputils.Post(url, &httputils.Params{
				Body: result,
			})
			if err != nil {
				logger.Error("agent %v 结果返回失败", uuid)
				return
			}
			logger.Info("agent %v 执行命令结果已经返回", uuid)
		})
	}

	logger.Info("批次agents正在远程执行命令: %v", d.Command)
	response.Success(c, nil, "远程命令已经发送")
}

func runCommandCallback(uuid string, command string, taskId string, callback func(client.TaskCmdResult)) {
	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunCommand(command)
		if err != nil {
			logger.Error("run command error, agent:%s, command:%s", uuid, command)
		}
		re := client.TaskCmdResult{
			TaskId: taskId,
			Result: client.CmdResult{
				MachineUUID: uuid,
				MachineIP:   agent.IP,
				RetCode:     data.RetCode,
				Stdout:      data.Stdout,
				Stderr:      data.Stderr,
			},
		}

		callback(re)
	}
	callback(client.TaskCmdResult{
		TaskId: taskId,
		Result: client.CmdResult{MachineUUID: uuid},
	})
}
