/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package handler

import (
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/pkg/utils/os"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func FirewalldConfigHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	config, err := uos.OS().Config()
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   config,
	}
	return c.Send(resp_msg)
}

func FirewalldDefaultZoneHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zone := msg.Data.(string)
	default_zone, err := uos.OS().FirewalldSetDefaultZone(zone)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   default_zone,
	}
	return c.Send(resp_msg)
}

func FirewalldZoneConfigHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zone := msg.Data.(string)
	default_zone, err := uos.OS().FirewalldZoneConfig(zone)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   default_zone,
	}
	return c.Send(resp_msg)
}

func FirewalldServiceAddHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	zone := zps[0]
	service := zps[1]
	err := uos.OS().FirewalldServiceAdd(zone, service)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}
	return c.Send(resp_msg)
}

func FirewalldServiceRemoveHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	zone := zps[0]
	service := zps[1]
	err := uos.OS().FirewalldServiceRemove(zone, service)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}
	return c.Send(resp_msg)
}

func FirewalldSourceAddHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	zone := zps[0]
	source := zps[1]
	err := uos.OS().FirewalldSourceAdd(zone, source)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}
	return c.Send(resp_msg)
}

func FirewalldSourceRemoveHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	zone := zps[0]
	source := zps[1]
	err := uos.OS().FirewalldSourceRemove(zone, source)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}
	return c.Send(resp_msg)
}

func FirewalldRestartHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	Restart := uos.OS().Restart()
	if !Restart {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "重启防火墙失败",
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   Restart,
	}
	return c.Send(resp_msg)
}

func FirewalldStopHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	Stop := uos.OS().Stop()
	if !Stop {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "关闭防火墙失败",
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   Stop,
	}
	return c.Send(resp_msg)
}

func FirewalldZonePortAddHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	zone := zps[0]
	port := zps[1]
	proto := zps[2]
	add, err := uos.OS().AddZonePort(zone, port, proto)

	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   add,
	}
	return c.Send(resp_msg)
}

func FirewalldZonePortDelHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	zp := msg.Data.(string)
	zps := strings.Split(zp, ",")
	if len(zps) < 3 {
		return c.Send(&protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "Invalid data format",
		})
	}

	zone := zps[0]
	port := zps[1]
	proto := zps[2]
	// 调用删除端口函数
	result, err := uos.OS().DelZonePort(zone, port, proto)
	if err != nil {
		return c.Send(&protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		})
	}
	// 发送成功响应
	return c.Send(&protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   result,
	})
}
