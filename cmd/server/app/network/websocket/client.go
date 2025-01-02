/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package websocket

import (
	"encoding/json"
	"runtime/debug"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// 用户连接
type Client struct {
	Addr          string                        // 客户端地址
	Socket        *websocket.Conn               // 用户连接
	Send          chan *global.WebsocketSendMsg // 待发送的数据
	HeartbeatTime uint64                        // 用户上次心跳时间
}

// 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan *global.WebsocketSendMsg, 100),
		HeartbeatTime: firstTime,
	}

	return
}

// 读取客户端数据
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("write stop: %s, %s", string(debug.Stack()), r)
		}
	}()

	for {
		_, _, err := c.Socket.ReadMessage()
		if err != nil {
			CliManager.Unregister <- c
			logger.Debug("读取客户端数据 关闭连接: %s, %s", c.Addr, err)
			return
		}
	}
}

// 向客户端写数据
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("write stop: %s, %s", string(debug.Stack()), r)
		}
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			// 发送数据错误 关闭连接
			logger.Debug("Client发送数据 关闭连接: %s, OK: %t", c.Addr, ok)
			return
		}
		msg_json, err := json.Marshal(message)
		if err != nil {
			logger.Error("fail to marshal json data in ws client.write(): %s, %+v, %s", c.Addr, message, err.Error())
			return
		}
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg_json); err != nil {
			logger.Error("err while writing message to websocket client(remote: %s): %s", c.Socket.RemoteAddr(), err.Error())
		}
	}
}

// 发送数据
func (c *Client) SendMsg(msg_type global.WebsocketSendMsgType, msg string) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("SendMsg stop: %s, %s", string(debug.Stack()), r)
		}
	}()

	c.Send <- &global.WebsocketSendMsg{
		MsgType: msg_type,
		Msg:     msg,
	}
}

// 监测系统日志警告推送到前端
func SendWarnMsgToWeb(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			klog.Warningln("SendWarnMsgToWeb success exit")
			return
		case data := <-global.WARN_MSG:
			CliManager.Broadcast <- data
		}
	}

}

// 心跳超时
func (c *Client) IsHeartbeatTimeout() (timeout bool) {
	currentTime := uint64(time.Now().Unix())
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}
