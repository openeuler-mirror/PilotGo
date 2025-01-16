/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"bytes"
	"context"
	"sync"

	Websocket "gitee.com/openeuler/PilotGo/cmd/server/app/network/websocket"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 终端连接功能删除
func WebTerminal(c *gin.Context) {
	// 升级协议并获得socket连接
	conn, err := Websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("获取websocket连接失败:%s", err.Error())
		return
	}
	defer conn.Close()

	// 后端获取到前端传来的主机信息,以此建立ssh客户端
	msg := c.DefaultQuery("msg", "")
	ssh, err := Websocket.DecodedMsgToSSHClient(msg)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	// 生成ssh socket客户端，建立session、client、channel
	client, err := ssh.NewSSHClient()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		conn.Close()
		return
	}
	defer client.Close()

	terminal, err := Websocket.NewTerminal(conn, client)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	defer terminal.Close()

	var bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	var logBuff = bufPool.Get().(*bytes.Buffer)
	defer func() {
		logBuff.Reset()
		defer bufPool.Put(logBuff)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := terminal.LoopRead(logBuff, ctx)
		if err != nil {
			logger.Error("LoopRead 退出：%#v", err)
		}
		cancel()
	}()
	go func() {
		defer wg.Done()
		err := terminal.SessionWait()
		if err != nil {
			logger.Error("SessionWait 退出：%#v", err)
		}
		cancel()
	}()
	wg.Wait()
}
