package handler

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/global"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/localstorage"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	mc "gitee.com/openeuler/PilotGo/pkg/utils/message/common"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/pkg/utils/os"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
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

func AgentOverviewHandler(c *network.SocketClient, msg *protocol.Message) error {
	var ip string
	var sysinfo *common.SystemInfo
	var diskusage []common.DiskUsageINfo
	var memoryInfo *common.MemoryConfig
	var cpuinfo *common.CPUInfo

	var resp_msg *protocol.Message
	var result mc.AgentOverview
	var err error
	logger.Debug("process agent overview command:%s", msg.String())
	// 获取IP信息
	ip, err = uos.OS().GetHostIp()
	if err != nil {
		logger.Error("failed to get IP: %s", err.Error())
		goto OnError
	}

	// 获取host主机信息
	sysinfo, err = uos.OS().GetHostInfo()
	if err != nil {
		logger.Error("failed to get host info: %s", err.Error())
		goto OnError
	}

	// 获取host cpu信息
	cpuinfo, err = uos.OS().GetCPUInfo()
	if err != nil {
		logger.Error("failed to get host info: %s", err.Error())
		goto OnError
	}

	// 获取内存信息
	memoryInfo, err = uos.OS().GetMemoryConfig()
	if err != nil {
		logger.Error("failed to get memory info: %s", err.Error())
		goto OnError
	}

	// 获取磁盘信息
	diskusage, err = uos.OS().GetDiskUsageInfo()
	if err != nil {
		logger.Error("failed to get dist usage info: %s", err.Error())
		goto OnError
	}

	result = mc.AgentOverview{
		IP:         ip,
		SysInfo:    sysinfo,
		DiskUsage:  diskusage,
		MemoryInfo: memoryInfo,
		CpuInfo:    cpuinfo,
	}

	if sysinfo.Platform == "nestos-for-container" {
		result.IsImmutable = true
	}

	resp_msg = &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   result,
	}
	return c.Send(resp_msg)

OnError:
	resp_msg = &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: -1,
		Error:  err.Error(),
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

	if erros != nil || errcpu != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  erros.Error(),
			Data:   systemAndCPUInfo,
		}
		return c.Send(resp_msg)
	}
	systemAndCPUInfo = common.SystemAndCPUInfo{
		IP:              os.IP,
		Platform:        os.Platform,
		PlatformVersion: os.PlatformVersion,
		PrettyName:      os.PrettyName,
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
