package handler

import (
	"gitee.com/openeuler/PilotGo/app/agent/global"
	"gitee.com/openeuler/PilotGo/app/agent/network"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/utils/os/common"
)

func ReadFilePatternHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	file, ok := msg.Data.(map[string]interface{})
	if !ok {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "get msg data error",
		}
		return c.Send(resp_msg)
	}

	data, err := utils.ReadFilePattern(file["path"].(string), file["name"].(string))
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

func SaveFileHandler(c *network.SocketClient, msg *protocol.Message) error {
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

	err = utils.SaveFile(result.Path, result.Name, result.Text)
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
