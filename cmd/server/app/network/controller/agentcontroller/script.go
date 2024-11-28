/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"regexp"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func RunCmd(c *gin.Context) {
	logger.Debug("process get agent request")
	// TODO: process agent info
	uuid := c.Query("uuid")
	cmd := c.Query("cmd")

	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunCommand(cmd)
		if err != nil {
			logger.Error("run command error, agent:%s, cmd:%s", uuid, cmd)
			response.Fail(c, gin.H{"status": false}, err.Error())
		}
		logger.Info("run command on agent result:%v", data)
		response.Success(c, nil, "run command success")
		return
	}

	logger.Info("unknown agent:%s", uuid)
	response.Fail(c, gin.H{"status": false}, "unknown agent")
}

func RunScriptWithBooleanCheck(c *gin.Context) {
	logger.Debug("process get agent script request")
	uuid := c.Query("uuid")
	cmd := c.Query("cmd")

	// 调用检测高危命令
	if containsDangerousCommand(cmd) {
		logger.Warn("Detected dangerous command")
		response.Fail(c, gin.H{"status": false}, "Dangerous command detected in script.")
		return
	}

	agent := agentmanager.GetAgent(uuid)
	if agent != nil {
		data, err := agent.RunCommand(cmd)
		if err != nil {
			logger.Error("run script error, agent:%s, cmd:%s", uuid, cmd)
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}
		logger.Info("run script on agent result:%v", data)
		response.Success(c, nil, "run script success")
		return
	}

	logger.Info("unknown agent:%s", uuid)
	response.Fail(c, gin.H{"status": false}, "unknown agent")
}

func containsDangerousCommand(content string) bool {
	for _, pattern := range dangerousCommandsList {
		matched, err := regexp.MatchString(pattern, content)
		if err != nil {
			logger.Error("Error matching pattern %s: %v\n", pattern, err)
			// TODO
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

var dangerousCommandsList = []string{
	`.*rm\s+-[r,f,rf].*`,
	`.*lvremove\s+-f.*`,
	`.*poweroff.*`,
	`.*shutdown\s+-[f,F,h,k,n,r,t,C].*`,
	`.*pvremove\s+-f.*`,
	`.*vgremove\s+-f.*`,
}
