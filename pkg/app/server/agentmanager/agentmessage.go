package agentmanager

import (
	"fmt"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

type AgentInfo struct {
	AgentVersion string `mapstructure:"agent_version"`
	AgentUUID    string `mapstructure:"agent_uuid"`
	IP           string `mapstructure:"IP"`
}

// 远程获取agent端的系统信息
func (a *Agent) AgentInfo() (*AgentInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentInfo,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &AgentInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind AgentInfo data error: %v", err)
		return nil, err
	}

	return info, nil
}

// 远程获取agent端的系统信息
func (a *Agent) GetOSInfo() (*common.SystemInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.OsInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.SystemInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetOSInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 远程获取agent端的CPU信息
func (a *Agent) GetCPUInfo() (*common.CPUInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CPUInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.CPUInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetCPUInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 远程获取agent端的内存信息
func (a *Agent) GetMemoryInfo() (*common.MemoryConfig, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.MemoryInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent: %s", err.Error())
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.MemoryConfig{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetMemoryInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetSysctlInfo() (*map[string]string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &map[string]string{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetSysctlInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 查看某个内核参数的值
func (a *Agent) SysctlView(args string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlView,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// 获取磁盘的使用情况
func (a *Agent) DiskUsage() ([]*common.DiskUsageINfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskUsage,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &[]*common.DiskUsageINfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind DiskUsage data error: %v", err)
		return nil, err
	}
	return *info, nil
}

// 获取磁盘的IO信息
func (a *Agent) DiskInfo() (*common.DiskIOInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.DiskIOInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind DiskInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

/*
挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘
*/
func (a *Agent) DiskMount(sourceDisk, destPath string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskMount,
		Data: sourceDisk + "," + destPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return err.Error(), err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

func (a *Agent) DiskUMount(diskPath string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskUMount,
		Data: diskPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return err.Error(), err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

func (a *Agent) DiskFormat(fileType, diskPath string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskFormat,
		Data: fileType + "," + diskPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// 获取当前TCP网络连接信息
func (a *Agent) NetTCP() (*common.NetConnect, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetTCP,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.NetConnect{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind NetTCP data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 获取当前UDP网络连接信息
func (a *Agent) NetUDP() (*common.NetConnect, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetUDP,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.NetConnect{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind NetUDP data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 获取网络读写字节／包的个数
func (a *Agent) NetIOCounter() (*common.IOCnt, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetIOCounter,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.IOCnt{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind NetIOCounter data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 获取网卡配置
func (a *Agent) NetNICConfig() (*common.NetInterfaceCard, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetNICConfig,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.NetInterfaceCard{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind NetNICConfig data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetAgentOSInfo() (*common.SystemAndCPUInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentOSInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, fmt.Errorf(resp_message.Error)
	}

	info := &common.SystemAndCPUInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetAgentOSInfo data error: %v", err)
		return nil, err
	}
	return info, nil
}

// 远程获取agent端的repo文件
func (a *Agent) GetRepoSource() ([]*common.RepoSource, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.GetRepoSource,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := &[]*common.RepoSource{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind data error: %v", err)
		return nil, resp_message.Error, err
	}
	return *info, resp_message.Error, nil
}

// 远程获取agent端的时间信息
func (a *Agent) GetTimeInfo() (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentTime,
		Data: struct{}{},
	}
	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to get time on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to get time on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}
