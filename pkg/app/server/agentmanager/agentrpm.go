package agentmanager

import (
	"fmt"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
)

// 获取全部安装的rpm包列表
func (a *Agent) AllRpm() ([]string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AllRpm,
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

	if v, ok := resp_message.Data.([]interface{}); ok {
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("failed to convert interface{} in allrpm")
}

// 获取源软件包名以及源
func (a *Agent) RpmSource(rpm string) (*common.RpmSrc, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RpmSource,
		Data: rpm,
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

	info := &common.RpmSrc{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind RpmSource data error:%s", err)
		return nil, err
	}
	return info, nil
}

// 获取软件包信息
func (a *Agent) RpmInfo(rpm string) (*common.RpmInfo, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RpmInfo,
		Data: rpm,
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

	info := &common.RpmInfo{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind RpmInfo data error:%s", err)
		return nil, "", err
	}
	return info, resp_message.Error, nil
}

// 获取源软件包名以及源
func (a *Agent) InstallRpm(rpm string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.InstallRpm,
		Data: rpm,
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

// 获取源软件包名以及源
func (a *Agent) RemoveRpm(rpm string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RemoveRpm,
		Data: rpm,
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
