/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package network

import (
	"net"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/config"
	pnet "gitee.com/openeuler/PilotGo/pkg/utils/message/net"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type AgentMessageHandler func(*SocketClient, *protocol.Message) error

type SocketClient struct {
	conn             net.Conn
	messageChan      chan *protocol.Message
	MessageProcesser *protocol.MessageProcesser
	exitChan         chan struct{}
	exitError        error
}

func NewSocketClient() *SocketClient {
	return &SocketClient{
		MessageProcesser: protocol.NewMessageProcesser(),
		messageChan:      make(chan *protocol.Message, 100),
		exitChan:         make(chan struct{}),
	}
}

// 启动Socket
func (c *SocketClient) Run(conf *config.Server) error {
	if err := c.connect(conf.Addr); err != nil {
		logger.Error("connect server error:%s", err.Error())
		return err
	}

	<-c.exitChan
	c.conn.Close()
	return nil
}

// 启动Socket
func (c *SocketClient) exitWithError(err error) {
	c.exitError = err
	c.exitChan <- struct{}{}
}

func (c *SocketClient) connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.conn = conn

	go func(c *SocketClient) {
		for {
			msg := <-c.messageChan
			logger.Debug("send message response:%d, message length:%d", msg.Type, len(msg.String()))

			data := msg.Encode()
			sendData := protocol.TlvEncode(data)

			err := pnet.SendBytes(c.conn, sendData)
			if err != nil {
				logger.Error("send byte data error:%s", err.Error())
				c.exitWithError(err)
				return
			}
		}
	}(c)

	go func(c *SocketClient) {
		readBuff := []byte{}
		for {
			buff := make([]byte, 1024)
			n, err := c.conn.Read(buff)
			if err != nil {
				logger.Error("socket read data error:%s", err)
				c.exitWithError(err)
				return
			}
			readBuff = append(readBuff, buff[:n]...)

			//切割frame
			for {
				i, f := protocol.TlvDecode(&readBuff)
				if i != 0 {
					readBuff = readBuff[i:]
					go func(c *SocketClient, f *[]byte) {
						msg := protocol.ParseMessage(*f)
						c.MessageProcesser.ProcessMessage(c, msg)
					}(c, f)
				} else {
					break
				}
			}
		}
	}(c)
	return nil
}

func (c *SocketClient) Send(msg *protocol.Message) error {
	c.messageChan <- msg
	return nil
}

func (c *SocketClient) BindHandler(t int, f AgentMessageHandler) {
	// c.MessageProcesser.BindHandler(t, (protocol.MessageHandler)(f))
	c.MessageProcesser.BindHandler(t, func(c protocol.MessageContext, msg *protocol.Message) error {
		return f(c.(*SocketClient), msg)
	})
}

func (c *SocketClient) Close() error {
	return c.conn.Close()
}
