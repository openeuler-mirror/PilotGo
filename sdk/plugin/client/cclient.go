/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Aug 07 16:18:53 2025 +0800
 */
package client

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"github.com/gin-gonic/gin"
)

type RunCommandCallback func([]*common.CmdResult)

type CallbackHandler struct {
	RunCommandCallback RunCommandCallback
	TaskLen            int
}
type GetTagsCallback func([]string) []common.Tag

type CClient struct {
	token    string
	Registry registry.Registry

	// 用于处理主机标签
	getTagsCallback GetTagsCallback

	// 用于event消息处理
	EventChan        chan *common.EventMessage
	EventCallbackMap map[int]common.EventCallback

	// 用于异步command及script执行结果处理机
	asyncCmdResultChan      chan *common.AsyncCmdResult
	cmdProcessorCallbackMap map[string]CallbackHandler
}

func NewClient(token string, reg registry.Registry) *CClient {
	return &CClient{
		token:    token,
		Registry: reg,

		EventChan:        make(chan *common.EventMessage, 20),
		EventCallbackMap: make(map[int]common.EventCallback),

		asyncCmdResultChan:      make(chan *common.AsyncCmdResult, 20),
		cmdProcessorCallbackMap: make(map[string]CallbackHandler),
	}
}

// RegisterHandlers 注册一些插件标准的API接口，清单如下：
func (c *CClient) RegisterHandlers(router *gin.Engine) {
	api := router.Group("/plugin_manage/api/v1/")
	{
		api.GET("/gettags", func(c *gin.Context) {
			c.Set("__internal__client_instance", c)
		}, TagsHandler)

		api.PUT("/command_result", func(c *gin.Context) {
			c.Set("__internal__client_instance", c)
		}, RunCommandResultHandler)
	}

	// TODO: start command result process service
	// c.startEventProcessor()
	c.startCommandResultProcessor()
}
