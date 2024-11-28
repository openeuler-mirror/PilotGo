/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package handler

import (
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func CronStartHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	dataStr := msg.Data.(string)
	parts := strings.Split(dataStr, ",")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return c.Send(&protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "Invalid ID format",
		})
	}
	spec := parts[1]
	command := parts[2]

	err = common.CronStart(id, spec, command)
	status := 0
	if err != nil {
		status = -1
	}

	respMsg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: status,
	}
	if err != nil {
		respMsg.Error = err.Error()
	} else {
		respMsg.Data = "任务已开始"
	}

	return c.Send(respMsg)
}

func CronStopAndDelHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	msgg := msg.Data.(string)
	message := strings.Split(msgg, ",")
	id, _ := strconv.Atoi(message[0])

	err := common.StopAndDel(id)
	if err != nil {
		respMsg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(respMsg)
	} else {
		respMsg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "任务已暂停",
		}
		return c.Send(respMsg)
	}
}
