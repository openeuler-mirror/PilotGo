/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"gitee.com/openeuler/PilotGo/cmd/agent/app/global"
	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ConfigfileHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		response.Fail(c, nil, "uuid参数缺失")
		return
	}
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
