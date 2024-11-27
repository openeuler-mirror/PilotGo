/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package handler

import (
	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/pkg/utils/os"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func AllRpmHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	allrpm, err := uos.OS().GetAllRpm()
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
		Data:   allrpm,
	}
	return c.Send(resp_msg)
}

func RpmSourceHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	rpmname := msg.Data.(string)
	rpmsource, err := uos.OS().GetRpmSource(rpmname)
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
		Data:   rpmsource,
	}
	return c.Send(resp_msg)
}

func RpmInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	rpmname := msg.Data.(string)
	rpminfo, err := uos.OS().GetRpmInfo(rpmname)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Data:   rpminfo,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   rpminfo,
		}
		return c.Send(resp_msg)
	}
}

func InstallRpmHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	rpmname := msg.Data.(string)

	err := uos.OS().InstallRpm(rpmname)

	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "",
		}
		return c.Send(resp_msg)
	}
}

func RemoveRpmHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	rpmname := msg.Data.(string)
	err := uos.OS().RemoveRpm(rpmname)

	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "",
		}
		return c.Send(resp_msg)
	}
}

func GetRepoSourceHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	repo, err := uos.OS().GetRepoSource()

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
		Data:   repo,
	}
	return c.Send(resp_msg)
}
