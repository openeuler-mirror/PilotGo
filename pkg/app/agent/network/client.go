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
 * Date: 2022-01-24 15:08:08
 * LastEditTime: 2022-04-08 13:08:25
 * Description: provide agent service functions.
 ******************************************************************************/
package network

import (
	"fmt"
	"net"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	pnet "openeluer.org/PilotGo/PilotGo/pkg/net"
	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*SocketClient, *protocol.Message) error

type SocketClient struct {
	conn             net.Conn
	messageChan      chan *protocol.Message
	MessageProcesser *protocol.MessageProcesser
}

func (c *SocketClient) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

	go func(c *SocketClient) {
		c.messageChan = make(chan *protocol.Message, 100)
		for {
			msg := <-c.messageChan
			logger.Debug("send message response:%d, message length:%d", msg.Type, len(msg.String()))

			data := msg.Encode()
			sendData := protocol.TlvEncode(data)

			err := pnet.Send(c.conn, sendData)
			if err != nil {
				logger.Error("send byte data error:%s", err.Error())
			}
		}
	}(c)

	go func(c *SocketClient) {
		readBuff := []byte{}
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				fmt.Println("read error:", err)
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
