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
 * LastEditTime: 2022-04-18 16:02:48
 * Description: agent main
 ******************************************************************************/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	aconfig "openeluer.org/PilotGo/PilotGo/pkg/app/agent/config"
	"openeluer.org/PilotGo/PilotGo/pkg/app/agent/localstorage"
	"openeluer.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	uos "openeluer.org/PilotGo/PilotGo/pkg/utils/os"
)

var agent_version = "v0.0.1"
var RESP_MSG = make(chan interface{})

func main() {
	fmt.Println("Start PilotGo agent.")

	// 加载系统配置
	err := aconfig.Init()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	// 初始化日志
	if err := logger.Init(&aconfig.Config().Logopts); err != nil {
		fmt.Println("logger init failed, please check the config file")
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo!")

	// 定时任务初始化
	if err := uos.CronInit(); err != nil {
		fmt.Println("cron init failed")
		os.Exit(-1)
	}

	// init agent info
	if err := localstorage.Init(); err != nil {
		fmt.Println("local storage init failed")
		os.Exit(-1)
	}
	logger.Info("agent uuid is:%s", localstorage.AgentUUID())

	go func(conf *aconfig.Server) {
		// 与server握手
		client := network.NewSocketClient()
		regitsterHandler(client)
		go FileMonitor(client)

		for {
			logger.Info("start to connect server")
			err = client.Run(&aconfig.Config().Server)
			if err != nil {
				logger.Error("socket client exit, error:%s", err.Error())
			}

			// 延迟5s+5s内随机时间
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			delayTime := time.Second*5 + time.Duration(r.Uint32()%5000*uint32(time.Millisecond))
			time.Sleep(delayTime)
		}
	}(&aconfig.Config().Server)

	// 文件监控初始化
	RESP_MSG = make(chan interface{})
	if err := FileMonitorInit(); err != nil {
		fmt.Println("config file monitor init failed")
		os.Exit(-1)
	}

	// 信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}

func FileMonitorInit() error {
	//获取IP
	IP, err := utils.RunCommand("hostname -I")
	if err != nil {
		return fmt.Errorf("can not to get IP")
	}
	str := strings.Split(IP, " ")
	IP = str[0]

	// 1、NewWatcher 初始化一个 watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 2、使用 watcher 的 Add 方法增加需要监听的文件或目录到监听队列中
	go func() {
		err = watcher.Add(uos.RepoPath)
		if err != nil {
			fmt.Println("failed to monitor repo")
		}
		logger.Info("start to monitor repo")
	}()

	go func() {
		err = watcher.Add(uos.NetWorkPath)
		if err != nil {
			fmt.Println("failed to monitor network")
		}
		logger.Info("start to monitor network")
	}()

	//3、创建新的 goroutine，等待管道中的事件或错误
	done := make(chan bool)
	go func() {
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}
				fileExt := path.Ext(e.Name)
				if strings.Contains(fileExt, ".sw") || strings.Contains(fileExt, "~") || strings.Contains(e.Name, "~") {
					continue
				}

				if e.Op&fsnotify.Write == fsnotify.Write {
					RESP_MSG <- fmt.Sprintf("机器 %s 上的文件已被修改 : %s", IP, e.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("error:", err)
			}
		}
	}()
	<-done
	return nil
}

func FileMonitor(client *network.SocketClient) {
	for data := range RESP_MSG {
		if data == nil {
			continue
		}

		msg := &protocol.Message{
			UUID:   uuid.New().String(),
			Type:   protocol.FileMonitor,
			Status: 0,
			Data:   data,
		}

		if err := client.Send(msg); err != nil {
			fmt.Println("send message failed, error:", err)
		}

	}
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
		IP, err := utils.RunCommand("hostname -I")
		if err != nil {
			fmt.Println("获取IP失败!")
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
		fmt.Println("process agent info command:", msg.String())
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
	c.BindHandler(protocol.FirewalldConfig, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

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
	c.BindHandler(protocol.FirewalldRestart, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		add, err := uos.AddZonePort(zone, port)

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
		fmt.Println("process agent info command:", msg.String())
		zp := msg.Data.(string)
		zps := strings.Split(zp, ",")
		zone := zps[0]
		port := zps[1]
		del, err := uos.DelZonePort(zone, port)

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
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
	c.BindHandler(protocol.GetRepoFile, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

		repos, err := uos.GetFiles(uos.RepoPath)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			datas := make([]map[string]string, 0)
			for _, repo := range repos {
				if ok := strings.Contains(repo, ".repo"); !ok {
					continue
				}
				datas = append(datas, map[string]string{"path": uos.RepoPath, "type": "repo", "name": repo})
			}
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   datas,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.GetNetWorkFile, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

		network, err := uos.GetFiles(uos.NetWorkPath)
		if err != nil {
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Error:  err.Error(),
			}
			return c.Send(resp_msg)
		} else {
			data := make(map[string]string, 0)
			for _, n := range network {
				if ok := strings.Contains(n, "ifcfg-e"); !ok {
					continue
				}
				data = map[string]string{"path": uos.NetWorkPath, "type": "network", "name": n}
			}
			resp_msg := &protocol.Message{
				UUID:   msg.UUID,
				Type:   msg.Type,
				Status: 0,
				Data:   data,
			}
			return c.Send(resp_msg)
		}
	})
	c.BindHandler(protocol.ReadFile, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())

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
		fmt.Println("process agent info command:", msg.String())

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
}
