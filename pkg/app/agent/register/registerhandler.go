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
 * LastEditTime: 2023-04-19 16:21:52
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

	c.BindHandler(protocol.RunCommand, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process run command:%s", msg.String())

		d := &struct {
			Command string
		}{}

		err := msg.BindData(d)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  "parse data error:" + err.Error(),
			}
			return c.Send(resp_msg)
		}

		retCode, stdout, stderr, err := utils.RunCommandnew(d.Command)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  "run command error:" + err.Error(),
			}
			return c.Send(resp_msg)
		}

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data: &utils.CmdResult{
				RetCode: retCode,
				Stdout:  stdout,
				Stderr:  stderr,
			},
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.RunScript, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process run script command:%s", msg.String())
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "run script result",
		}
		return c.Send(resp_msg)
	})

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

	c.BindHandler(protocol.OsInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		sysinfo, err := uos.OS().GetHostInfo()
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
			Data:   sysinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.CPUInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		cpuinfo, err := uos.OS().GetCPUInfo()
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
			Data:   cpuinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.MemoryInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		memoryinfo, err := uos.OS().GetMemoryConfig()
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
			Data:   memoryinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.SysctlInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		// TODO: process error
		sysctlinfo, _ := uos.OS().GetSysctlConfig()

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   sysctlinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.SysctlChange, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		args := msg.Data.(string)
		sysctlchange, _ := uos.OS().TempModifyPar(args)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   sysctlchange,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.SysctlView, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		args := msg.Data.(string)
		sysctlview, _ := uos.OS().GetVarNameValue(args)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   sysctlview,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.ServiceList, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		servicelist, _ := uos.OS().GetServiceList()

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   servicelist,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.ServiceStatus, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		service := msg.Data.(string)
		servicestatus, _ := uos.OS().GetServiceStatus(service)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   servicestatus,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.ServiceRestart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		service := msg.Data.(string)
		err := uos.OS().RestartService(service)

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
				Data:   "重启成功",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.ServiceStart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		service := msg.Data.(string)
		err := uos.OS().StartService(service)
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
				Data:   "启动成功",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.ServiceStop, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		service := msg.Data.(string)
		err := uos.OS().StopService(service)

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
				Data:   "关闭服务成功",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.AllRpm, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		allrpm, err := uos.OS().GetAllRpm()
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
			Data:   allrpm,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.RpmSource, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		rpmname := msg.Data.(string)
		rpmsource, err := uos.OS().GetRpmSource(rpmname)
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
			Data:   rpmsource,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.RpmInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		rpmname := msg.Data.(string)
		rpminfo, err := uos.OS().GetRpmInfo(rpmname)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Data:   rpminfo,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   rpminfo,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.InstallRpm, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		rpmname := msg.Data.(string)

		err := uos.OS().InstallRpm(rpmname)

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
				Data:   "",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.RemoveRpm, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		rpmname := msg.Data.(string)
		err := uos.OS().RemoveRpm(rpmname)

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
				Data:   "",
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.DiskUsage, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.DiskInfo, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.DiskMount, func(c *network.SocketClient, msg *protocol.Message) error {
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

	})
	c.BindHandler(protocol.DiskUMount, func(c *network.SocketClient, msg *protocol.Message) error {
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

	})
	c.BindHandler(protocol.DiskFormat, func(c *network.SocketClient, msg *protocol.Message) error {
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

	})
	c.BindHandler(protocol.NetTCP, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		nettcp, err := uos.OS().GetTCP()
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
			Data:   nettcp,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.NetUDP, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netudp, err := uos.OS().GetUDP()

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
			Data:   netudp,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.NetIOCounter, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netio, err := uos.OS().GetIOCounter()

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
			Data:   netio,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.NetNICConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netnic, err := uos.OS().GetNICConfig()

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
			Data:   netnic,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.CurrentUser, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		user_info := uos.OS().GetCurrentUserInfo()

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   user_info,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.AllUser, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.AddLinuxUser, func(c *network.SocketClient, msg *protocol.Message) error {
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

	})
	c.BindHandler(protocol.DelUser, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.ChangePermission, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.ChangeFileOwner, func(c *network.SocketClient, msg *protocol.Message) error {
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
	})
	c.BindHandler(protocol.AgentOSInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		os, erros := uos.OS().GetHostInfo()
		cpu, errcpu := uos.OS().GetCPUInfo()
		systemAndCPUInfo := common.SystemAndCPUInfo{}

		if erros != nil && errcpu != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  erros.Error(),
				Data:   systemAndCPUInfo,
			}
			return c.Send(resp_msg)
		} else if erros != nil && errcpu == nil {
			systemAndCPUInfo.ModelName = cpu.ModelName
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  erros.Error(),
				Data:   systemAndCPUInfo,
			}
			return c.Send(resp_msg)
		} else if erros == nil && errcpu != nil {
			systemAndCPUInfo.IP = os.IP
			systemAndCPUInfo.Platform = os.Platform
			systemAndCPUInfo.PlatformVersion = os.PlatformVersion
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  errcpu.Error(),
				Data:   systemAndCPUInfo,
			}
			return c.Send(resp_msg)
		}
		systemAndCPUInfo = common.SystemAndCPUInfo{
			IP:              os.IP,
			Platform:        os.Platform,
			PlatformVersion: os.PlatformVersion,
			ModelName:       cpu.ModelName,
		}
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   systemAndCPUInfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.FirewalldConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		config, err := uos.OS().Config()
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
			Data:   config,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldDefaultZone, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zone := msg.Data.(string)
		default_zone, err := uos.OS().FirewalldSetDefaultZone(zone)
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
			Data:   default_zone,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldZoneConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zone := msg.Data.(string)
		default_zone, err := uos.OS().FirewalldZoneConfig(zone)
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
			Data:   default_zone,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldServiceAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		service := zps[1]
		err := uos.OS().FirewalldServiceAdd(zone, service)
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
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldServiceRemove, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		service := zps[1]
		err := uos.OS().FirewalldServiceRemove(zone, service)
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
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldSourceAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		source := zps[1]
		err := uos.OS().FirewalldSourceAdd(zone, source)
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
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldSourceRemove, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		source := zps[1]
		err := uos.OS().FirewalldSourceRemove(zone, source)
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
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldRestart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		Restart := uos.OS().Restart()
		if !Restart {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  "重启防火墙失败",
			}
			return c.Send(resp_msg)
		}
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   Restart,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldStop, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		Stop := uos.OS().Stop()
		if !Stop {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: -1,
				Error:  "关闭防火墙失败",
			}
			return c.Send(resp_msg)
		}
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   Stop,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldZonePortAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		proto := zps[2]
		add, err := uos.OS().AddZonePort(zone, port, proto)

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
			Data:   add,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.FirewalldZonePortDel, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		proto := zps[2]
		del, err := uos.OS().DelZonePort(zone, port, proto)

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
			Data:   del,
		}
		return c.Send(resp_msg)

	})
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
	c.BindHandler(protocol.GetRepoSource, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		repo, err := common.GetRepoSource()

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
			Data:   repo,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.GetNetWorkConnectInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		network, err := uos.OS().ConfigNetworkConnect()
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
			Data:   network,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.GetNetWorkConnInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		network, err := uos.OS().GetNetworkConnInfo()
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
			Data:   network,
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.RestartNetWork, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		msgg := msg.Data.(string)
		err := uos.OS().RestartNetwork(msgg)
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
		}
		return c.Send(resp_msg)

	})
	c.BindHandler(protocol.GetNICName, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		nic_name, err := uos.OS().GetNICName()
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
			Data:   nic_name,
		}
		return c.Send(resp_msg)

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
