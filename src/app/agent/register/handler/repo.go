package handler

import (
	"gitee.com/openeuler/PilotGo/app/agent/network"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/utils/os"
)

func RepoConfigHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process repo files command:%s", msg.String())
	repofiles, err := uos.OS().GetRepoConfig()
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
		Data:   repofiles,
	}
	return c.Send(resp_msg)
}
