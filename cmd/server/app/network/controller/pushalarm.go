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

	go client.Read()
	go client.Write()

	// 用户连接事件
	websocket.CliManager.Register <- client
}
