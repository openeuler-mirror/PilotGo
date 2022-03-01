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
 * Date: 2021-11-18 10:25:52
 * LastEditTime: 2022-03-01 13:13:02
 * Description: agent main
 ******************************************************************************/
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"openeluer.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
	uos "openeluer.org/PilotGo/PilotGo/pkg/utils/os"
)

var agent_uuid = uuid.New().String()
var agent_version = "v0.0.1"

func main() {
	fmt.Println("Start PilotGo agent.")

	// init agent info

	// 加载系统配置
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}
	url := conf.S.ServerIP + ":" + strconv.Itoa(conf.SocketPort)

	// 初始化日志

	// 与server握手
	client := &network.SocketClient{
		MessageProcesser: protocol.NewMessageProcesser(),
	}

	if err := client.Connect(url); err != nil {
		fmt.Println("connect server failed, error:", err)
		os.Exit(-1)
	}
	regitsterHandler(client)

	// go Send_heartbeat(client)
	select {}

}

func Send_heartbeat(client *network.SocketClient) {
	for {
		msg := &protocol.Message{
			UUID: uuid.New().String(),
			Type: protocol.Heartbeat,
			Data: "连接正常",
		}

		if err := client.Send(msg); err != nil {
			fmt.Println("send message failed, error:", err)
		}
		fmt.Println("send heartbeat message")

		time.Sleep(time.Second * 5)

		// 接受远程指令并执行
		if false {
			break
		}
	}

	out, err := utils.RunCommand("date")
	if err == nil {
		fmt.Println(string(out))
	}
}

func regitsterHandler(c *network.SocketClient) {
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
		fmt.Println("process run script command:", msg.String())
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   "run script result",
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.AgentInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data: map[string]string{
				"agent_version": agent_version,
				"agent_uuid":    agent_uuid,
			},
		}
		return c.Send(resp_msg)
	})

	c.BindHandler(protocol.OsInfo, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

		sysctlinfo := uos.GetSysConfig()

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   sysctlinfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.SysctlChange, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
		rpmname := msg.Data.(string)

		err := uos.InstallRpm(rpmname)

		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   err.Error(),
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
		fmt.Println("process agent info command:", msg.String())
		rpmname := msg.Data.(string)
		rpmremove := uos.RemoveRpm(rpmname)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   rpmremove,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.DiskUsage, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())
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
		fmt.Println("process agent info command:", msg.String())

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
}
