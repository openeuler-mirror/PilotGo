package handler

import (
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func CronStartHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	msgg := msg.Data.(string)
	message := strings.Split(msgg, ",")
	id, _ := strconv.Atoi(message[0])
	spec := message[1]
	command := message[2]

	err := common.CronStart(id, spec, command)
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
			Data:   "任务已开始",
		}
		return c.Send(resp_msg)
	}
}

func CronStopAndDelHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	msgg := msg.Data.(string)
	message := strings.Split(msgg, ",")
	id, _ := strconv.Atoi(message[0])

	err := common.StopAndDel(id)
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
			Data:   "任务已暂停",
		}
		return c.Send(resp_msg)
	}
}
