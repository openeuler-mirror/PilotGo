/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentmanager

import (
	mc "gitee.com/openeuler/PilotGo/pkg/utils/message/common"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
)

// 远程获取agent端的主机的概览信息
func (a *Agent) AgentOverview() (*mc.AgentOverview, error) {
	info := &mc.AgentOverview{}
	_, err := a.SendMessageWrapper(protocol.AgentOverview, nil, "failed to send agent overview message", -1, info, "AgentOverview")
	return info, err
}

type AgentInfo struct {
	AgentVersion string `mapstructure:"agent_version"`
	AgentUUID    string `mapstructure:"agent_uuid"`
	IP           string `mapstructure:"IP"`
}

// 远程获取agent端的系统信息
func (a *Agent) AgentInfo() (*AgentInfo, error) {
	info := &AgentInfo{}
	_, err := a.SendMessageWrapper(protocol.AgentInfo, nil, "failed to run script on agent", -1, info, "AgentInfo")
	return info, err
}

// 远程获取agent端的系统信息
func (a *Agent) GetOSInfo() (*common.SystemInfo, error) {
	info := &common.SystemInfo{}
	_, err := a.SendMessageWrapper(protocol.OsInfo, struct{}{}, "failed to run script on agent", -1, info, "GetOSInfo")
	return info, err
}

// 远程获取agent端的CPU信息
func (a *Agent) GetCPUInfo() (*common.CPUInfo, error) {
	info := &common.CPUInfo{}
	_, err := a.SendMessageWrapper(protocol.CPUInfo, struct{}{}, "failed to run script on agent", -1, info, "GetCPUInfo")
	return info, err
}

// 远程获取agent端的内存信息
func (a *Agent) GetMemoryInfo() (*common.MemoryConfig, error) {
	info := &common.MemoryConfig{}
	_, _ = a.SendMessageWrapper(protocol.MemoryInfo, struct{}{}, "failed to run script on agent", -1, info, "GetMemoryInfo")
	return info, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetSysctlInfo() (*map[string]string, error) {
	info := &map[string]string{}
	_, _ = a.SendMessageWrapper(protocol.SysctlInfo, struct{}{}, "failed to run script on agent", -1, info, "GetSysctlInfo")
	return info, nil
}

// 查看某个内核参数的值
func (a *Agent) SysctlView(args string) (string, error) {
	responseMessage, _ := a.SendMessageWrapper(protocol.SysctlView, args, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), nil
}

// 获取磁盘的使用情况
func (a *Agent) DiskUsage() ([]*common.DiskUsageINfo, error) {
	info := &[]*common.DiskUsageINfo{}
	_, _ = a.SendMessageWrapper(protocol.DiskUsage, struct{}{}, "failed to run script on agent", -1, info, "DiskUsage")
	return *info, nil
}

// 获取磁盘的IO信息
func (a *Agent) DiskInfo() (*common.DiskIOInfo, error) {
	info := &common.DiskIOInfo{}
	_, _ = a.SendMessageWrapper(protocol.DiskInfo, struct{}{}, "failed to run script on agent", -1, info, "DiskInfo")
	return info, nil
}

/*
挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘
*/
func (a *Agent) DiskMount(sourceDisk, destPath string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.DiskMount, sourceDisk+","+destPath, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}

func (a *Agent) DiskUMount(diskPath string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.DiskUMount, diskPath, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}

func (a *Agent) DiskFormat(fileType, diskPath string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.DiskFormat, fileType+","+diskPath, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}

// 获取当前TCP网络连接信息
func (a *Agent) NetTCP() (*common.NetConnect, error) {
	info := &common.NetConnect{}
	_, err := a.SendMessageWrapper(protocol.NetTCP, struct{}{}, "failed to run script on agent", -1, info, "NetTCP")
	return info, err
}

// 获取当前UDP网络连接信息
func (a *Agent) NetUDP() (*common.NetConnect, error) {
	info := &common.NetConnect{}
	_, err := a.SendMessageWrapper(protocol.NetUDP, struct{}{}, "failed to run script on agent", -1, info, "NetUDP")
	return info, err
}

// 获取网络读写字节／包的个数
func (a *Agent) NetIOCounter() (*common.IOCnt, error) {
	info := &common.IOCnt{}
	_, err := a.SendMessageWrapper(protocol.NetIOCounter, struct{}{}, "failed to run script on agent", -1, info, "NetIOCounter")
	return info, err
}

// 获取网卡配置
func (a *Agent) NetNICConfig() (*common.NetInterfaceCard, error) {
	info := &common.NetInterfaceCard{}
	_, err := a.SendMessageWrapper(protocol.NetNICConfig, struct{}{}, "failed to run script on agent", -1, info, "NetNICConfig")
	return info, err
}

// 远程获取agent端的内核信息
func (a *Agent) GetAgentOSInfo() (*common.SystemAndCPUInfo, error) {
	info := &common.SystemAndCPUInfo{}
	_, err := a.SendMessageWrapper(protocol.AgentOSInfo, struct{}{}, "failed to run script on agent", -1, info, "GetAgentOSInfo")
	return info, err
}

// 远程获取agent端的repo文件
func (a *Agent) GetRepoSource() ([]*common.RepoSource, error) {
	info := &[]*common.RepoSource{}
	_, err := a.SendMessageWrapper(protocol.GetRepoSource, struct{}{}, "failed to run script on agent", -1, info, "GetRepoSource")
	return *info, err
}

// 远程获取agent端的时间信息
func (a *Agent) GetTimeInfo() (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.AgentTime, struct{}{}, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}
