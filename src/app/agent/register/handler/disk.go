package handler

import (
	"strings"

	"gitee.com/openeuler/PilotGo/app/agent/network"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/utils/os"
)

func DiskUsageHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())

	diskusage, err := uos.OS().GetDiskUsageInfo()
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
		Data:   diskusage,
	}
	return c.Send(resp_msg)
}

func DiskInfoHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	diskinfo, err := uos.OS().GetDiskInfo()
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
		Data:   diskinfo,
	}
	return c.Send(resp_msg)
}

func DiskMountHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	disk := msg.Data.(string)
	disks := strings.Split(disk, ",")
	source := disks[0]
	dest := disks[1]
	info, err := uos.OS().DiskMount(source, dest)
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
		Data:   info,
	}
	return c.Send(resp_msg)
}

func DiskUMountHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	disk := msg.Data.(string)
	info, err := uos.OS().DiskUMount(disk)
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
		Data:   info,
	}
	return c.Send(resp_msg)
}

func DiskFormatHandler(c *network.SocketClient, msg *protocol.Message) error {
	logger.Debug("process agent info command:%s", msg.String())
	disk := msg.Data.(string)
	disks := strings.Split(disk, ",")
	fileType := disks[0]
	diskPath := disks[1]
	info, err := uos.OS().DiskFormat(fileType, diskPath)
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
		Data:   info,
	}
	return c.Send(resp_msg)

}