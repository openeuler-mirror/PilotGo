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
	"net"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"k8s.io/klog/v2"
)

type SocketServer struct {
	// MessageProcesser *protocol.MessageProcesser
	OnAccept func(net.Conn) error
	OnStop   func()
}

func SocketServerInit(conf *options.SocketServer, stopCh <-chan struct{}) error {
	server := &SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}

	go func() {
		if err := server.Run(conf.Addr, stopCh); err != nil {
			logger.Error("socket server init run failed: %s", err.Error())
		}
	}()
	return nil
}

func (s *SocketServer) Run(addr string, stopCh <-chan struct{}) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		<-stopCh
		klog.Warningln("SocketServer prepare stop")
		listener.Close()

	}()
	logger.Debug("Waiting for agents")
	for {
		select {
		case <-stopCh:
			klog.Warning("SocketServer prepare stop")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					klog.Warningln("SocketServer success exit")
				} else {
					klog.Errorf("SocketServer run  error:%v", err)
				}
				continue
			}
			if err := s.OnAccept(conn); err != nil {
				klog.Errorf("failed to add and run agent: %v", err)
			}
		}

	}
}

func (s *SocketServer) Stop() {
}
