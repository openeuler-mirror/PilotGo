/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugins licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 24 10:02:04 2024 +0800
 */
package sdk

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"github.com/gin-gonic/gin"
)

func UnPluginListenEventHandler() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		logger.Info("插件监听注册goroutine已启动")
		defer logger.Info("插件取消监听goroutine退出")
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.Info("接收到退出信号: %s", s.String())
				UnPluginListenEvent()
				os.Exit(0)
			default:
				logger.Info("接收到未知信号: %s", s.String())
			}
		}
	}()
}
func RegisterEventHandlers(router *gin.Engine, c *client.Client) {

	api := router.Group("/plugin_manage/api/v1/")
	{
		api.POST("/event", eventHandler)
	}
	plugin_client = c
	startEventProcessor(c)
}

func eventHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		return
	}
	var msg common.EventMessage
	if err := json.Unmarshal(j, &msg); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		return
	}

	ProcessEvent(&msg)
}

func eventPluginServer() (string, error) {
	plugins, err := plugin_client.GetPlugins()
	if err != nil {
		return "", err
	}

	var eventServer string
	for _, p := range plugins {
		if p.Name == "event" {
			eventServer = p.Url
			break
		}
	}

	if eventServer == "" {
		return "", errors.New("event plugin not found")
	}

	return eventServer, nil
}

func registerEventCallback(eventType int, callback common.EventCallback) {
	plugin_client.EventCallbackMap[eventType] = callback
}

func unregisterEventCallback(eventType int) {
	delete(plugin_client.EventCallbackMap, eventType)
}

func ProcessEvent(event *common.EventMessage) {
	plugin_client.EventChan <- event
}

func startEventProcessor(c *client.Client) {
	go func(c *client.Client) {
		for {
			e := <-c.EventChan

			// TODO: process event message
			cb, ok := c.EventCallbackMap[e.MessageType]
			if ok {
				cb(e)
			}
		}
	}(c)

}
