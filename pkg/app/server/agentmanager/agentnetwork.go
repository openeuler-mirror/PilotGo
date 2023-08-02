package agentmanager

import (
	"fmt"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 远程获取agent端的网络连接信息
func (a *Agent) GetNetWorkConnectInfo() (*map[string]string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.GetNetWorkConnectInfo,
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

	info := &map[string]string{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetSysctlInfo data error:", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}

// 获取agent的基础网络配置
func (a *Agent) GetNetWorkConnInfo() (*common.NetworkConfig, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.GetNetWorkConnInfo,
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

	info := &common.NetworkConfig{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind GetNetWorkConnInfo data error:", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}

// 获取网卡名字
func (a *Agent) GetNICName() (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.GetNICName,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), resp_message.Error, nil
}

// 重启网卡配置
func (a *Agent) RestartNetWork(NIC string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RestartNetWork,
		Data: NIC,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Error, nil
}

// 查看配置文件内容
func (a *Agent) ReadFile(filepath string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ReadFile,
		Data: filepath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), resp_message.Error, nil
}

// 更新配置文件
func (a *Agent) UpdateFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		FilePath: filepath,
		FileName: filename,
		FileText: text,
	}
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.EditFile,
		Data: updatefile,
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

	info := &common.UpdateFile{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind UpdateFile data error:", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}
