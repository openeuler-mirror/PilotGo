package handler

import (
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/global"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func ReadFileHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	file := msg.Data.(string)
	data, err := utils.FileReadString(file)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   data,
		}
		return c.Send(resp_msg)
	}
}

func EditFileHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	result := &common.UpdateFile{}
	err := msg.BindData(result)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  err.Error(),
		}
		return c.Send(resp_msg)
	}

	LastVersion, err := utils.UpdateFile(result.Path, result.Name, result.Text)
	result.FileLastVersion = LastVersion
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
		Data:   result,
	}
	return c.Send(resp_msg)
}

func AgentConfigHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	p, ok := msg.Data.(map[string]interface{})
	if ok {
		var ConMess global.ConfigMessage
		ConMess.Machine_uuid = p["Machine_uuid"].(string)
		ConMess.ConfigName = p["ConfigName"].(string)
		err := global.Configfsnotify(ConMess, c)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   "正常监听文件",
			}
			return c.Send(resp_msg)
		}
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "监控文件有误",
		}
		return c.Send(resp_msg)
	}
}
