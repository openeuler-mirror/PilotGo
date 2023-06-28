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
 * LastEditTime: 2023-06-28 11:21:42
 * Description: socket client register
 ******************************************************************************/
package register

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/global"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/localstorage"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/register/handler"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeuler.org/PilotGo/PilotGo/pkg/utils/os"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

var agent_version = "v0.0.1"

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

	c.BindHandler(protocol.AgentInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		IP, err := uos.OS().GetHostIp()
		if err != nil {
			logger.Error("failed to get IP: %s", err.Error())
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  fmt.Sprintf("failed to get IP: %s", err.Error()),
			}
			return c.Send(resp_msg)
		}

		result := struct {
			AgentVersion string `json:"agent_version"`
			IP           string `json:"IP"`
			AgentUUID    string `json:"agent_uuid"`
		}{
			AgentVersion: agent_version,
			IP:           IP,
			AgentUUID:    localstorage.AgentUUID(),
		}

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   result,
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.OsInfo, handler.OSInfoHandler)
	c.BindHandler(protocol.CPUInfo, handler.CPUInfoHandler)
	c.BindHandler(protocol.MemoryInfo, handler.MemoryInfoHandler)

	c.BindHandler(protocol.SysctlInfo, handler.SysctlInfoHandler)
	c.BindHandler(protocol.SysctlChange, handler.SysctlChangeHandler)
	c.BindHandler(protocol.SysctlView, handler.SysctlViewHandler)

	c.BindHandler(protocol.ServiceList, handler.ServiceListHandler)
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

	c.BindHandler(protocol.AgentOSInfo, handler.AgentOSInfoHandler)

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

	c.BindHandler(protocol.CronStart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		msgg := msg.Data.(string)
		message := strings.Split(msgg, ",")
		id, _ := strconv.Atoi(message[0])
		spec := message[1]
		command := message[2]

		err := common.CronStart(id, spec, command)
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
				Data:   "任务已开始",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.CronStopAndDel, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		msgg := msg.Data.(string)
		message := strings.Split(msgg, ",")
		id, _ := strconv.Atoi(message[0])

		err := common.StopAndDel(id)
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
				Data:   "任务已暂停",
			}
			return c.Send(resp_msg)
		}
	})

	c.BindHandler(protocol.ReadFile, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.EditFile, func(c *network.SocketClient, msg *protocol.Message) error {
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

		LastVersion, err := utils.UpdateFile(result.FilePath, result.FileName, result.FileText)
		result.FileVersion = LastVersion
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
				Data:   result,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.AgentTime, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		timeinfo, err := uos.OS().GetTime()
		if err != nil {
			logger.Debug(err.Error())
		}

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   timeinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.AgentConfig, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
}
