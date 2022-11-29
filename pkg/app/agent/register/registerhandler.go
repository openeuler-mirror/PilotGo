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
 * LastEditTime: 2022-07-05 14:10:23
 * Description: socket client register
 ******************************************************************************/
package register

import (
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

	out, err := utils.RunCommand("date")
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
		IP, err := utils.RunCommand("hostname -I")
		if err != nil {
			logger.Debug("获取IP失败!")
		}
		str := strings.Split(IP, " ")
		IP = str[0]
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data: map[string]string{
				"agent_version": agent_version,
				"IP":            IP,
				"agent_uuid":    localstorage.AgentUUID(),
			},
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.OsInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		sysinfo := uos.GetHostInfo()

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

		cpuinfo := uos.GetCPUInfo()

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

		memoryinfo := uos.GetMemoryConfig()

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
		sysctlinfo, _ := uos.GetSysctlConfig()

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
		sysctlchange := uos.TempModifyPar(args)

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
		sysctlview := uos.GetVarNameValue(args)

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

		servicelist, _ := uos.GetServiceList()

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
		servicestatus, _ := uos.GetServiceStatus(service)

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
		err := uos.RestartService(service)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		err := uos.StartService(service)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		err := uos.StopService(service)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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

		allrpm := uos.GetAllRpm()

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
		rpmsource, _ := uos.GetRpmSource(rpmname)

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
		rpminfo, Err, err := uos.GetRpmInfo(rpmname)
		if Err != nil && err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   rpminfo,
				Error:  Err.Error(),
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

		err := uos.InstallRpm(rpmname)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		err := uos.RemoveRpm(rpmname)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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

		diskusage := uos.GetDiskUsageInfo()

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
		diskinfo := uos.GetDiskInfo()

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   diskinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.CreateDiskPath, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		mountpath := msg.Data.(string)
		creatdiskpath := uos.CreateDiskPath(mountpath)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   creatdiskpath,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.DiskMount, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		disk := msg.Data.(string)
		disks := strings.Split(disk, ",")
		source := disks[0]
		dest := disks[1]
		mountpath := uos.DiskMount(source, dest)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   mountpath,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.DiskUMount, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		disk := msg.Data.(string)
		diskPath := uos.DiskUMount(disk)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   diskPath,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.DiskFormat, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		disk := msg.Data.(string)
		disks := strings.Split(disk, ",")
		fileType := disks[0]
		diskPath := disks[1]
		formatpath := uos.DiskFormat(fileType, diskPath)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   formatpath,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.NetTCP, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		nettcp, err := uos.GetTCP()
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   err,
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   nettcp,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.NetUDP, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netudp, err := uos.GetUDP()

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   err,
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   netudp,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.NetIOCounter, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netio, err := uos.GetIOCounter()

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   err,
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   netio,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.NetNICConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		netnic, err := uos.GetNICConfig()

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   err,
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   netnic,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.CurrentUser, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		user_info := uos.GetCurrentUserInfo()

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

		user_all := uos.GetAllUserInfo()

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
		err := uos.AddLinuxUser(username, password)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		user_del, err := uos.DelUser(username)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		user_per, err := uos.ChangePermission(permission, file)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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
		user_ower, err := uos.ChangeFileOwner(fileType, diskPath)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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

		os := uos.GetHostInfo()
		cpu := uos.GetCPUInfo()
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   os.IP + ";" + os.Platform + ";" + os.PlatformVersion + ";" + cpu.ModelName,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.FirewalldConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		config, err := uos.Config()
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   config,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldDefaultZone, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zone := msg.Data.(string)
		default_zone, err := uos.FirewalldSetDefaultZone(zone)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   default_zone,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldZoneConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zone := msg.Data.(string)
		default_zone, err := uos.FirewalldZoneConfig(zone)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   default_zone,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldServiceAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		service := zps[1]
		err := uos.FirewalldServiceAdd(zone, service)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   struct{}{},
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldServiceRemove, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		service := zps[1]
		err := uos.FirewalldServiceRemove(zone, service)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   struct{}{},
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldSourceAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		source := zps[1]
		err := uos.FirewalldSourceAdd(zone, source)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   struct{}{},
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldSourceRemove, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		source := zps[1]
		err := uos.FirewalldSourceRemove(zone, source)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   struct{}{},
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldRestart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		Restart := uos.Restart()
		if !Restart {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  "重启防火墙失败",
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   Restart,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldStop, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		Stop := uos.Stop()
		if !Stop {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  "关闭防火墙失败",
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   Stop,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldZonePortAdd, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		proto := zps[2]
		add, err := uos.AddZonePort(zone, port, proto)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   add,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.FirewalldZonePortDel, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		proto := zps[2]
		del, err := uos.DelZonePort(zone, port, proto)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   del,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.CronStart, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		msgg := msg.Data.(string)
		message := strings.Split(msgg, ",")
		id, _ := strconv.Atoi(message[0])
		spec := message[1]
		command := message[2]

		err := uos.CronStart(id, spec, command)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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

		err := uos.StopAndDel(id)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
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

		repo, err := uos.GetRepoSource()

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   repo,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.GetNetWorkConnectInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		network, err := uos.ConfigNetworkConnect()
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {

			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   network,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.GetNetWorkConnInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		network, err := uos.GetNetworkConnInfo()
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   network,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.RestartNetWork, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		msgg := msg.Data.(string)
		err := uos.RestartNetwork(msgg)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   struct{}{},
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.GetNICName, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		nic_name, err := uos.GetNICName()
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   nic_name,
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
				Status: 0,
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

		file := msg.Data.(map[string]interface{})
		filepath := file["path"]
		filename := file["name"]
		text := file["text"]

		LastVersion, err := uos.UpdateFile(filepath, filename, text)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   LastVersion,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.AgentTime, func(c *network.SocketClient, msg *protocol.Message) error {
		logger.Debug("process agent info command:%s", msg.String())

		timeinfo := uos.GetTime()

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
			ConMess.ConfigName = p["ConfigName"].(string)
			err := global.Configfsnotify(ConMess, c)
			if err != nil {
				resp_msg := &protocol.Message{
					UUID:   msg.UUID,
					Type:   msg.Type,
					Status: 0,
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
				Status: 0,
				Error:  "监控文件有误",
			}
			return c.Send(resp_msg)
		}
	})
}
