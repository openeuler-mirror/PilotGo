/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/websocket"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/user"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func PushAlarmHandler(c *gin.Context) {
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("%s", err.Error())
		return
	}
	logger.Debug("webSocket 建立连接: %s", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := websocket.NewClient(conn.RemoteAddr().String(), conn, currentTime)

	u, ok := c.Get("user")
	if !ok {
		logger.Error("webSocket 连接失败: 找不到用户信息")
		return
	}

	go client.Read(u.(*user.User).Username)
	go client.Write()

	// 用户连接事件
	websocket.CliManager.Register <- client

	messages := websocket.CliManager.SendMsgBuffer.GetAll()
	for _, msg := range messages {
		_msg, ok := msg.(*global.WebsocketSendMsg)
		if !ok {
			logger.Error("websocketSendMsg assert error: %+v", msg)
			continue
		}
		client.SendMsg(_msg)
	}
}
