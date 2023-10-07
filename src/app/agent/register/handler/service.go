package handler

import (
	"gitee.com/PilotGo/PilotGo/app/agent/network"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils/message/protocol"
	uos "gitee.com/PilotGo/PilotGo/utils/os"
)

func ServiceListHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	serviceList, _ := uos.OS().GetServiceList()

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   serviceList,
	}
	return c.Send(resp_msg)
}

func GetServiceHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	service := msg.Data.(string)
	serviceInfo, _ := uos.OS().GetService(service)

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   serviceInfo,
	}
	return c.Send(resp_msg)
}

func ServiceStatusHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	service := msg.Data.(string)
	serviceStatus, _ := uos.OS().GetServiceStatus(service)

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   serviceStatus,
	}
	return c.Send(resp_msg)
}

func ServiceRestartHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	service := msg.Data.(string)
	err := uos.OS().RestartService(service)

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
			Data:   "重启成功",
		}
		return c.Send(resp_msg)
	}
}

func ServiceStartHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	service := msg.Data.(string)
	err := uos.OS().StartService(service)
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
			Data:   "启动成功",
		}
		return c.Send(resp_msg)
	}
}

func ServiceStopHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	service := msg.Data.(string)
	err := uos.OS().StopService(service)

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
			Data:   "关闭服务成功",
		}
		return c.Send(resp_msg)
	}
}
