package agentmanager

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	"github.com/google/uuid"
)

func (a *Agent) GetRepoConfig() (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RepoConfig,
		Data: "",
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to get repo config on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to get repo config on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := ""
	err = resp_message.BindData(&info)
	if err != nil {
		logger.Error("bind repo cnfig data error: %s", err)
		return "", err
	}
	return info, nil
}
