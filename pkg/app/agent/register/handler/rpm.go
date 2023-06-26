package handler

import (
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeuler.org/PilotGo/PilotGo/pkg/utils/os"
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
