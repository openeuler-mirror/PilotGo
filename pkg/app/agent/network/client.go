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
 * LastEditTime: 2022-04-06 19:23:17
 * Description: provide agent service functions.
 ******************************************************************************/
package network

import (
	"fmt"
	"net"

	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*SocketClient, *protocol.Message) error

type SocketClient struct {
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
}

func (c *SocketClient) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c.conn = conn

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
			i, f := protocol.TlvDecode(&readBuff)
			if i != 0 {
				readBuff = readBuff[i:]
				go func() {
					msg := protocol.ParseMessage(*f)
					c.MessageProcesser.ProcessMessage(c, msg)
				}()
			}
		}
	}(c)
	return nil
}

func (c *SocketClient) Send(msg *protocol.Message) error {
	data := msg.Encode()
	sendData := protocol.TlvEncode(data)

	data_length := len(sendData)
	send_count := 0
	for {
		n, err := c.conn.Write(sendData[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count = send_count + n
			break
		}
	}
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
