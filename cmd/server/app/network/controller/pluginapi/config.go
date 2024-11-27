/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package pluginapi

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	batchservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ApplyConfig(c *gin.Context) {

}

// 配置文件下发
func FileDeploy(c *gin.Context) {
	fd := &struct {
		DeployBatch    common.Batch `json:"deploybatch"`
		DeployPath     string       `json:"deploypath"`
		DeployFileName string       `json:"deployname"`
		DeployText     string       `json:"deployfile"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	path := fd.DeployPath
	filename := fd.DeployFileName
	text := fd.DeployText

	if len(path) == 0 {
		response.Fail(c, nil, "路径为空，请检查文件路径")
		return
	}
	if len(filename) == 0 {
		response.Fail(c, nil, "文件名为空，请检查文件名字")
		return
	}
	if len(text) == 0 {
		response.Fail(c, nil, "文件内容为空，请重新文件内容")
		return
	}
	f := func(uuid string) batchservice.R {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			return common.NodeResult{UUID: uuid,
				Error: "get agent failed"}
		}

		_, Err, err := agent.SaveFile(path, filename, text)
		if len(Err) != 0 || err != nil {
			return common.NodeResult{UUID: uuid,
				Error: Err + err.Error()}
		}
		return common.NodeResult{UUID: uuid}
	}

	result := batchservice.BatchProcess(&fd.DeployBatch, f, path, filename, text)
	response.Success(c, result, "文件下发结果")
}

func GetNodeFiles(c *gin.Context) {
	fd := &struct {
		DeployBatch common.Batch `json:"deploybatch"`
		Path        string       `json:"path"`
		FileName    string       `json:"filename"`
	}{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	f := func(uuid string) batchservice.R {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			logger.Error("get agent failed, agent:%s", uuid)
			return common.NodeResult{UUID: uuid,
				Error: "get agent failed"}
		}
		//获取此路径下的所有文件，检查文件是否符合正则表达式
		data, _, err := agent.ReadFilePattern(fd.Path, fd.FileName)
		if err != nil {
			logger.Error("failed to read the file, agent:%s,err:%s", uuid, err)
			return common.NodeResult{UUID: uuid,
				Error: "failed to read the file"}
		}
		return common.NodeResult{UUID: uuid,
			Data: data}
	}

	rs := batchservice.BatchProcess(&fd.DeployBatch, f)
	response.Success(c, rs, "文件获取完成")
}
