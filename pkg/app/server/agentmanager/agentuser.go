package agentmanager

import (
	"fmt"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取当前用户信息
func (a *Agent) CurrentUser() (*common.CurrentUser, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CurrentUser,
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

	info := &common.CurrentUser{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind CurrentUser data error:%s", err)
		return nil, err
	}
	return info, nil
}

// 获取所有用户的信息
func (a *Agent) AllUser() ([]*common.AllUserInfo, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AllUser,
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

	info := &[]*common.AllUserInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind AllUser data error:%s", err)
		return nil, err
	}
	return *info, nil
}

// 创建新的用户，并新建家目录
func (a *Agent) AddLinuxUser(username, password string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AddLinuxUser,
		Data: username + "," + password,
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

// 删除用户
func (a *Agent) DelUser(username string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DelUser,
		Data: username,
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
