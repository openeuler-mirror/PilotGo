/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package register

import (
	"time"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/register/handler"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	uos "gitee.com/openeuler/PilotGo/pkg/utils/os"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/google/uuid"
)

func Send_heartbeat(client *network.SocketClient) {
	for {
		msg := &protocol.Message{
			UUID: uuid.New().String(),
			Type: protocol.Heartbeat,
			Data: "连接正常",
		}

		if err := client.Send(msg); err != nil {
			logger.Debug("send message failed, error:%s", err)
		}
		logger.Debug("send heartbeat message")

		time.Sleep(time.Second * 5)

		// 接受远程指令并执行
		if false {
			break
		}
	}

	out, err := uos.OS().GetTime()
	if err == nil {
		logger.Debug("%s", string(out))
	}
}

func RegitsterHandler(c *network.SocketClient) {
	c.BindHandler(protocol.Heartbeat, func(c *network.SocketClient, msg *protocol.Message) error {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "连接正常",
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.RunCommand, handler.RunCommandHandler)
	c.BindHandler(protocol.RunScript, handler.RunScriptHandler)
	// 针对脚本运维所做的可以执行各种脚本的拓展
	c.BindHandler(protocol.AgentRunScripts, handler.AgentRunScriptsHandler)

	c.BindHandler(protocol.AgentOverview, handler.AgentOverviewHandler)
	c.BindHandler(protocol.AgentInfo, handler.AgentInfoHandler)
	c.BindHandler(protocol.AgentTime, handler.AgentTimeHandler)
	c.BindHandler(protocol.AgentOSInfo, handler.AgentOSInfoHandler)
	c.BindHandler(protocol.OsInfo, handler.OSInfoHandler)
	c.BindHandler(protocol.CPUInfo, handler.CPUInfoHandler)
	c.BindHandler(protocol.MemoryInfo, handler.MemoryInfoHandler)

	c.BindHandler(protocol.SysctlInfo, handler.SysctlInfoHandler)
	c.BindHandler(protocol.SysctlChange, handler.SysctlChangeHandler)
	c.BindHandler(protocol.SysctlView, handler.SysctlViewHandler)

	c.BindHandler(protocol.ServiceList, handler.ServiceListHandler)
	c.BindHandler(protocol.GetService, handler.GetServiceHandler)
	c.BindHandler(protocol.ServiceRestart, handler.ServiceRestartHandler)
	c.BindHandler(protocol.ServiceStart, handler.ServiceStartHandler)
	c.BindHandler(protocol.ServiceStop, handler.ServiceStopHandler)

	c.BindHandler(protocol.AllRpm, handler.AllRpmHandler)
	c.BindHandler(protocol.RpmSource, handler.RpmSourceHandler)
	c.BindHandler(protocol.RpmInfo, handler.RpmInfoHandler)
	c.BindHandler(protocol.InstallRpm, handler.InstallRpmHandler)
	c.BindHandler(protocol.RemoveRpm, handler.RemoveRpmHandler)
	c.BindHandler(protocol.GetRepoSource, handler.GetRepoSourceHandler)

	c.BindHandler(protocol.DiskUsage, handler.DiskUsageHandler)
	c.BindHandler(protocol.DiskInfo, handler.DiskInfoHandler)
	c.BindHandler(protocol.DiskMount, handler.DiskMountHandler)
	c.BindHandler(protocol.DiskUMount, handler.DiskUMountHandler)
	c.BindHandler(protocol.DiskFormat, handler.DiskFormatHandler)

	c.BindHandler(protocol.NetTCP, handler.NetTCPHandler)
	c.BindHandler(protocol.NetUDP, handler.NetUDPHandler)
	c.BindHandler(protocol.NetIOCounter, handler.NetIOCounterHandler)
	c.BindHandler(protocol.NetNICConfig, handler.NetNICConfigHandler)
	c.BindHandler(protocol.GetNetWorkConnectInfo, handler.GetNetWorkConnectInfoHandler)
	c.BindHandler(protocol.GetNetWorkConnInfo, handler.GetNetWorkConnInfoHandler)
	c.BindHandler(protocol.RestartNetWork, handler.RestartNetWorkHandler)
	c.BindHandler(protocol.GetNICName, handler.GetNICNameHandler)

	c.BindHandler(protocol.CurrentUser, handler.CurrentUserHandler)
	c.BindHandler(protocol.AllUser, handler.AllUserHandler)
	c.BindHandler(protocol.AddLinuxUser, handler.AddLinuxUserHandler)
	c.BindHandler(protocol.DelUser, handler.DelUserHandler)
	c.BindHandler(protocol.ChangePermission, handler.ChangePermissionHandler)
	c.BindHandler(protocol.ChangeFileOwner, handler.ChangeFileOwnerHandler)

	c.BindHandler(protocol.FirewalldConfig, handler.FirewalldConfigHandler)
	c.BindHandler(protocol.FirewalldDefaultZone, handler.FirewalldDefaultZoneHandler)
	c.BindHandler(protocol.FirewalldZoneConfig, handler.FirewalldZoneConfigHandler)
	c.BindHandler(protocol.FirewalldServiceAdd, handler.FirewalldServiceAddHandler)
	c.BindHandler(protocol.FirewalldServiceRemove, handler.FirewalldServiceRemoveHandler)
	c.BindHandler(protocol.FirewalldSourceAdd, handler.FirewalldSourceAddHandler)
	c.BindHandler(protocol.FirewalldSourceRemove, handler.FirewalldSourceRemoveHandler)
	c.BindHandler(protocol.FirewalldRestart, handler.FirewalldRestartHandler)
	c.BindHandler(protocol.FirewalldStop, handler.FirewalldStopHandler)
	c.BindHandler(protocol.FirewalldZonePortAdd, handler.FirewalldZonePortAddHandler)
	c.BindHandler(protocol.FirewalldZonePortDel, handler.FirewalldZonePortDelHandler)

	c.BindHandler(protocol.CronStart, handler.CronStartHandler)
	c.BindHandler(protocol.CronStopAndDel, handler.CronStopAndDelHandler)

	c.BindHandler(protocol.ReadFilePattern, handler.ReadFilePatternHandler)
	c.BindHandler(protocol.EditFile, handler.EditFileHandler)
	c.BindHandler(protocol.SaveFile, handler.SaveFileHandler)
	c.BindHandler(protocol.AgentConfig, handler.AgentConfigHandler)
}
