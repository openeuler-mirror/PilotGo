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
