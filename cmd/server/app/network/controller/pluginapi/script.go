/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
// 插件系统对外提供的api

package pluginapi

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/controller"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/gin-gonic/gin"
)

// 检查plugin接口调用权限
func AuthCheck(c *gin.Context) {
	_, err := jwt.ParsePluginClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "plugin token check error:" + err.Error()})
		c.Abort()
		return
	}

	c.Next()
}

// 远程运行脚本
func RunCommandHandler(c *gin.Context) {
	logger.Debug("process get agent request")

	d := &common.CmdStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	/*u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "run command",
	}
	auditlog.Add(log)*/
	logger.Debug("run command on agents :%v", d.Batch.MachineUUIDs)

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			/*log_s := &auditlog.AuditLog{
				LogUUID:    uuidservice.New().String(),
				ParentUUID: log.LogUUID,
				Module:     auditlog.ModulePlugin,
				Status:     auditlog.StatusOK,
				UserID:     u.ID,
				Action:     "run command",
				Message:    "agentuuid:" + uuid,
			}
			auditlog.Add(log_s)*/
			data, err := agent.RunCommand(d.Command)
			if err != nil {
				//auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
				//auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
				logger.Error("run command error, agent:%s, command:%s", uuid, d.Command)
			}
			logger.Debug("run command on agent result:%v", data)
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
	/*
		u, err := jwt.ParseUser(c)
		if err != nil {
			response.Fail(c, nil, "user token error:"+err.Error())
			return
		}
		log := &auditlog.AuditLog{
			LogUUID:    uuidservice.New().String(),
			ParentUUID: "",
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "run script",
		}
		auditlog.Add(log)*/

	logger.Debug("run script on agents :%v", d.Batch.MachineUUIDs)
	// Enabled according to the needs of the plugin
	positions, matchedCommands := controller.FindDangerousCommandsPos(d.Script)
	if len(positions) > 0 {
		logger.Debug("Matched Commands: %v", matchedCommands)
		str := strings.Join(matchedCommands, "\n")
		response.Fail(c, nil, "Dangerous commands detected in script: "+str)
		return
	}
	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil { /*
				log_s := &auditlog.AuditLog{
					LogUUID:    uuidservice.New().String(),
					ParentUUID: log.LogUUID,
					Module:     auditlog.ModulePlugin,
					Status:     auditlog.StatusOK,
					UserID:     u.ID,
					Action:     "run script",
					Message:    "agentuuid:" + uuid,
				}
				auditlog.Add(log_s)*/
			data, err := agent.RunScript(d.Script, d.Params)
			if err != nil {
				//auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
				//auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
				logger.Error("run script error, agent:%s, command:%s", uuid, d.Script)
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

	result := batch.BatchProcess(d.Batch, f, d.Script, d.Params)
	response.Success(c, result, "run script succeed")
}

// 异步远程运行脚本回调
func RunCommandAsyncHandler(c *gin.Context) {
	d := &common.CmdStruct{}
	if err := c.ShouldBind(d); err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}
	/*
		u, err := jwt.ParseUser(c)
		if err != nil {
			response.Fail(c, nil, "user token error:"+err.Error())
			return
		}
		log := &auditlog.AuditLog{
			LogUUID:    uuidservice.New().String(),
			ParentUUID: "",
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "Run Command Async",
		}
		auditlog.Add(log)*/

	// 获取插件地址和回调url
	name := c.Query("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		response.Fail(c, nil, "plugin not found: "+err.Error())
		return
	}
	parsedURL, err := url.Parse(p.Url)
	if err != nil {
		logger.Error("URL解析失败:%v", err)
		response.Fail(c, nil, "解析插件url失败")
		return
	}
	caller := "http://" + parsedURL.Host + "/plugin_manage/api/v1/command_result"

	taskId := time.Now().Format("20060102150405")
	macuuids := batch.GetBatchMachineUUIDS(d.Batch)
	for _, uuid := range macuuids {
		/*log_s := &auditlog.AuditLog{
			LogUUID:    uuidservice.New().String(),
			ParentUUID: log.LogUUID,
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "Run Command Async",
			Message:    "agentuuid:" + uuid,
		}
		auditlog.Add(log_s)*/
		go asyncCommandRunner(uuid, d.Command, taskId, caller)
	}

	logger.Info("批次agents正在远程执行命令: %v", d.Command)
	response.Success(c, struct {
		TaskID  string `json:"task_id"`
		TaskLen int    `json:"task_len"`
	}{
		TaskID:  taskId,
		TaskLen: len(macuuids),
	}, "远程命令已经发送")
}

func asyncCommandRunner(macuuid string, command string, taskId string, caller string) {
	var result common.AsyncCmdResult
	var res []*common.CmdResult

	agent := agentmanager.GetAgent(macuuid)
	if agent == nil {
		logger.Error("agent %v 不存在,请检查机器是否已经连接", macuuid)
	} else {
		data, err := agent.RunCommand(command)
		if err != nil {
			logger.Error("run command error, agent:%s, command:%s", macuuid, command)
		}
		r := &common.CmdResult{
			MachineUUID: macuuid,
			MachineIP:   agent.IP,
			RetCode:     data.RetCode,
			Stdout:      data.Stdout,
			Stderr:      data.Stderr,
		}
		res = append(res, r)
	}

	result = common.AsyncCmdResult{
		TaskID: taskId,
		Result: res,
	}

	_, err := httputils.Put(caller, &httputils.Params{
		Body: result,
	})
	if err != nil {
		logger.Error("agent %v 结果返回失败：%v", macuuid, err)
		return
	}
	logger.Info("agent %v 执行命令结果已经返回", macuuid)
}
