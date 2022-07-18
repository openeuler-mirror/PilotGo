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
 * Date: 2022-02-18 02:39:36
 * LastEditTime: 2022-03-04 02:25:56
 * Description: provide agent log manager functions.
 ******************************************************************************/
package network

import (
	"fmt"
	"net"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

type SocketServer struct {
	// MessageProcesser *protocol.MessageProcesser
	OnAccept func(net.Conn)
	OnStop   func()
}

func (s *SocketServer) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Debug("Waiting for agents")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		s.OnAccept(conn)
	}
}

func (s *SocketServer) Stop() {

}
