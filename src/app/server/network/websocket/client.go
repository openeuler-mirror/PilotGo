package websocket

import (
	"runtime/debug"
	"time"

	"gitee.com/PilotGo/PilotGo/app/server/agentmanager"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"github.com/gorilla/websocket"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	HeartbeatTime uint64          // 用户上次心跳时间
}

// 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
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
		c.Socket.WriteMessage(websocket.TextMessage, message)
	}
}

// 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("SendMsg stop: %s, %s", string(debug.Stack()), r)
		}
	}()

	c.Send <- msg
}

// 监测系统日志警告推送到前端
func SendWarnMsgToWeb() {
	for {
		data := <-agentmanager.WARN_MSG
		CliManager.Broadcast <- []byte(data.(string))
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
