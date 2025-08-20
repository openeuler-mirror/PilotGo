/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Thu Aug 07 17:35:12 2025 +0800
 */
package client

import (
	"encoding/json"
	"io"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func RunCommandResultHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		return
	}
	var result common.AsyncCmdResult
	if err := json.Unmarshal(j, &result); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		return
	}

	v, ok := c.Get("__internal__client_instance")
	if !ok {
		logger.Error("%v", "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		logger.Error("%v", "client获取失败")
		return
	}

	client.ProcessCommandResult(&result)

}

func TagsHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		response.Fail(c, gin.H{"status": false}, "没获取到："+err.Error())
		return
	}
	uuidTags := &struct {
		UUIDS []string `json:"uuids"`
	}{}
	if err := json.Unmarshal(j, &uuidTags); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		response.Fail(c, gin.H{"status": false}, "反序列化结果失败："+err.Error())
		return
	}

	v, ok := c.Get("__internal__client_instance")
	if !ok {
		logger.Error("%v", "未获取到client值信息")
		response.Fail(c, gin.H{"status": false}, "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		logger.Error("%v", "client获取失败")
		response.Fail(c, gin.H{"status": false}, "client获取失败")
		return
	}

	if client.getTagsCallback != nil {
		result := client.getTagsCallback(uuidTags.UUIDS)
		response.Success(c, result, "")
	} else {
		logger.Error("get tags callback not set")
		response.Fail(c, gin.H{"status": false}, "get tags callback not set")
	}
}
