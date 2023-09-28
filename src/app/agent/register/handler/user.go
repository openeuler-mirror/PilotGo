package handler

import (
	"strings"

	"gitee.com/PilotGo/PilotGo/app/agent/network"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils/message/protocol"
	uos "gitee.com/PilotGo/PilotGo/utils/os"
)

func CurrentUserHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	user_info := uos.OS().GetCurrentUserInfo()

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data:   user_info,
	}
	return c.Send(resp_msg)
}

func AllUserHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	user_all, err := uos.OS().GetAllUserInfo()
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
		Data:   user_all,
	}
	return c.Send(resp_msg)
}

func AddLinuxUserHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	user := msg.Data.(string)
	users := strings.Split(user, ",")
	username := users[0]
	password := users[1]
	err := uos.OS().AddLinuxUser(username, password)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Data:   err,
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "新增用户成功!",
		}
		return c.Send(resp_msg)
	}
}

func DelUserHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	username := msg.Data.(string)
	user_del, err := uos.OS().DelUser(username)
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
			Data:   user_del,
		}
		return c.Send(resp_msg)
	}
}

func ChangePermissionHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	data := msg.Data.(string)
	datas := strings.Split(data, ",")
	permission := datas[0]
	file := datas[1]
	user_per, err := uos.OS().ChangePermission(permission, file)

	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Data:   err,
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   user_per,
		}
		return c.Send(resp_msg)
	}
}

func ChangeFileOwnerHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	disk := msg.Data.(string)
	disks := strings.Split(disk, ",")
	fileType := disks[0]
	diskPath := disks[1]
	user_ower, err := uos.OS().ChangeFileOwner(fileType, diskPath)

	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Data:   err,
		}
		return c.Send(resp_msg)
	} else {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   user_ower,
		}
		return c.Send(resp_msg)
	}
}
