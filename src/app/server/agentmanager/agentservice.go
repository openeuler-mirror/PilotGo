package agentmanager

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/utils/os/common"
	"github.com/google/uuid"
)

// 查看服务列表
func (a *Agent) ServiceList() ([]*common.ListService, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceList,
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

	info := &[]*common.ListService{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind ServiceList data error:%s", err)
		return nil, err
	}
	return *info, nil
}

// 查看某个服务
func (a *Agent) GetService(service string) (*common.ServiceInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.GetService,
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
	serviceInfo := &common.ServiceInfo{}
	err = resp_message.BindData(serviceInfo)
	if err != nil {
		logger.Error("bind GetServiceInfo data error:%s", err)
		return nil, err
	}
	return serviceInfo, nil
}

// 重启服务
func (a *Agent) ServiceRestart(service string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceRestart,
		Data: service,
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

// 关闭服务
func (a *Agent) ServiceStop(service string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceStop,
		Data: service,
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

// 启动服务
func (a *Agent) ServiceStart(service string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceStart,
		Data: service,
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