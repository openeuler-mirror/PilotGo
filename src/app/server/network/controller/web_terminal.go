/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-04-20 16:48:55
 * LastEditTime: 2022-04-20 17:48:55
 * Description: web socket连接控制
 ******************************************************************************/

package controller

import (
	"bytes"
	"context"
	"sync"

	Websocket "gitee.com/openeuler/PilotGo/app/server/network/websocket"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 终端连接功能删除
func WS(c *gin.Context) {
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
	logBuff.Reset()
	defer bufPool.Put(logBuff)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := terminal.LoopRead(logBuff, ctx)
		if err != nil {
			logger.Error("%#v", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := terminal.SessionWait()
		if err != nil {
			logger.Error("%#v", err)
		}
		cancel()
	}()
	wg.Wait()
}
