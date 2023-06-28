package handler

import (
	"fmt"

	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/global"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/localstorage"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeuler.org/PilotGo/PilotGo/pkg/utils/os"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func OSInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process os info command:%s", msg.String())

	sysinfo, err := uos.OS().GetHostInfo()
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
		Data:   sysinfo,
	}
	return c.Send(resp_msg)
}

func CPUInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process cpu info command:%s", msg.String())

	cpuinfo, err := uos.OS().GetCPUInfo()
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
		Data:   cpuinfo,
	}
	return c.Send(resp_msg)
}

func MemoryInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process memory info command:%s", msg.String())

	memoryInfo, err := uos.OS().GetMemoryConfig()
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
		Data:   memoryInfo,
	}
	return c.Send(resp_msg)
}

func AgentInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	IP, err := uos.OS().GetHostIp()
	if err != nil {
		logger.Error("failed to get IP: %s", err.Error())
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  fmt.Sprintf("failed to get IP: %s", err.Error()),
		}
		return c.Send(resp_msg)
	}

	result := struct {
		AgentVersion string `json:"agent_version"`
		IP           string `json:"IP"`
		AgentUUID    string `json:"agent_uuid"`
	}{
		AgentVersion: global.AgentVersion,
		IP:           IP,
		AgentUUID:    localstorage.AgentUUID(),
	}

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   result,
	}
	return c.Send(resp_msg)
}

func AgentOSInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	os, erros := uos.OS().GetHostInfo()
	cpu, errcpu := uos.OS().GetCPUInfo()
	systemAndCPUInfo := common.SystemAndCPUInfo{}

	if erros != nil && errcpu != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  erros.Error(),
			Data:   systemAndCPUInfo,
		}
		return c.Send(resp_msg)
	} else if erros != nil && errcpu == nil {
		systemAndCPUInfo.ModelName = cpu.ModelName
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  erros.Error(),
			Data:   systemAndCPUInfo,
		}
		return c.Send(resp_msg)
	} else if erros == nil && errcpu != nil {
		systemAndCPUInfo.IP = os.IP
		systemAndCPUInfo.Platform = os.Platform
		systemAndCPUInfo.PlatformVersion = os.PlatformVersion
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  errcpu.Error(),
			Data:   systemAndCPUInfo,
		}
		return c.Send(resp_msg)
	}
	systemAndCPUInfo = common.SystemAndCPUInfo{
		IP:              os.IP,
		Platform:        os.Platform,
		PlatformVersion: os.PlatformVersion,
		ModelName:       cpu.ModelName,
	}
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   systemAndCPUInfo,
	}
	return c.Send(resp_msg)
}

func AgentTimeHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	timeinfo, err := uos.OS().GetTime()
	if err != nil {
		logger.Debug(err.Error())
	}

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   timeinfo,
	}
	return c.Send(resp_msg)
}
