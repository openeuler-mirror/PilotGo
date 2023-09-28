package agentmanager

import (
	"fmt"

	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils/message/protocol"
	"gitee.com/PilotGo/PilotGo/utils/os/common"
	"github.com/google/uuid"
)

// 查看配置文件内容
func (a *Agent) ReadConfigFile(filepath string) (string, string, error) {
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
func (a *Agent) UpdateConfigFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		Path: filepath,
		Name: filename,
		Text: text,
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
		logger.Error("bind UpdateFile data error:%s", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}
