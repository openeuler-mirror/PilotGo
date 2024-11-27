/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package websocket

import (
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"k8s.io/klog/v2"
)

var (
	CliManager = NewClientManager() // 管理者
)

// 连接管理
type ClientManager struct {
	Clients     map[*Client]bool // 全部的连接
	ClientsLock sync.RWMutex     // 读写锁
	Register    chan *Client     // 连接处理
	Unregister  chan *Client     // 断开连接处理程序
	Broadcast   chan []byte      // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return
}

/**************************  manager  ***************************************/

// GetClients
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})

	return
}

// 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if !result {
			return
		}
	}
}

// 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
}

// 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	delete(manager.Clients, client)
}

// 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)
	logger.Debug("EventRegister 用户建立连接: %s", client.Addr)
}

// 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)
	logger.Debug("EventUnregister 用户断开连接: %s", client.Addr)
}

// 管道处理程序
func (manager *ClientManager) Start(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			klog.Warningln("websocket CliManager success exit")
			return
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)
		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)
		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************  manager info  ***************************************/
// 定时清理超时连接
func ClearTimeoutWebSocketConnections() {
	clients := CliManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout() {
			logger.Debug("心跳时间超时 关闭连接: %s, %d", client.Addr, client.HeartbeatTime)

			client.Socket.Close()
		}
	}
}
