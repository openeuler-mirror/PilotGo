package handler

import (
	"gitee.com/PilotGo/PilotGo/app/agent/network"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils/message/protocol"
	uos "gitee.com/PilotGo/PilotGo/utils/os"
)

func SysctlInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process sysctl info command:%s", msg.String())

	// TODO: process error
	sysctlInfo, _ := uos.OS().GetSysctlConfig()

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   sysctlInfo,
	}
	return c.Send(resp_msg)
}

func SysctlChangeHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process sysctl change command:%s", msg.String())
	args := msg.Data.(string)
	sysctlChange, _ := uos.OS().TempModifyPar(args)

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   sysctlChange,
	}
	return c.Send(resp_msg)
}

func SysctlViewHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process sysctl view command:%s", msg.String())
	args := msg.Data.(string)
	sysctlView, _ := uos.OS().GetVarNameValue(args)

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   sysctlView,
	}
	return c.Send(resp_msg)
}
