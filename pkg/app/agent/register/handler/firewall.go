package handler

import (
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeuler.org/PilotGo/PilotGo/pkg/utils/os"
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
	zone := zps[0]
	port := zps[1]
	proto := zps[2]
	del, err := uos.OS().DelZonePort(zone, port, proto)

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
		Data:   del,
	}
	return c.Send(resp_msg)
}
