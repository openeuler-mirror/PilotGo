package handler

import (
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeuler.org/PilotGo/PilotGo/pkg/utils/os"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
)

func NetTCPHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	nettcp, err := uos.OS().GetTCP()
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
		Data:   nettcp,
	}
	return c.Send(resp_msg)
}

func NetUDPHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	netudp, err := uos.OS().GetUDP()

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
		Data:   netudp,
	}
	return c.Send(resp_msg)
}

func NetIOCounterHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	netio, err := uos.OS().GetIOCounter()

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
		Data:   netio,
	}
	return c.Send(resp_msg)
}

func NetNICConfigHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	netnic, err := uos.OS().GetNICConfig()

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
		Data:   netnic,
	}
	return c.Send(resp_msg)
}

func GetNetWorkConnectInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	network, err := uos.OS().ConfigNetworkConnect()
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
		Data:   network,
	}
	return c.Send(resp_msg)
}

func GetNetWorkConnInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	network, err := uos.OS().GetNetworkConnInfo()
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
		Data:   network,
	}
	return c.Send(resp_msg)
}

func RestartNetWorkHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	msgg := msg.Data.(string)
	err := uos.OS().RestartNetwork(msgg)
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

func GetNICNameHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	nic_name, err := uos.OS().GetNICName()
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
		Data:   nic_name,
	}
	return c.Send(resp_msg)
}
