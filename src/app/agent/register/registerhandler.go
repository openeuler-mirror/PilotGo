/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-07-05 13:03:16
 * LastEditTime: 2023-08-30 16:00:51
 * Description: socket client register
 ******************************************************************************/
package register

import (
	"time"

	"gitee.com/PilotGo/PilotGo/app/agent/network"
	"gitee.com/PilotGo/PilotGo/app/agent/register/handler"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils/message/protocol"
	uos "gitee.com/PilotGo/PilotGo/utils/os"
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
		logger.Debug(string(out))
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
	c.BindHandler(protocol.ServiceStatus, handler.ServiceStatusHandler)
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

	c.BindHandler(protocol.ReadFile, handler.ReadFileHandler)
	c.BindHandler(protocol.EditFile, handler.EditFileHandler)
	c.BindHandler(protocol.AgentConfig, handler.AgentConfigHandler)
}
