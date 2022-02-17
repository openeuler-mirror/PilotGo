package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"openeluer.org/PilotGo/PilotGo/pkg/app/agent/network"
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

	// 初始化日志

	// 与server握手
	client := &network.SocketClient{
		MessageProcesser: protocol.NewMessageProcesser(),
	}
	if err := client.Connect("192.168.47.128:8879"); err != nil {
		fmt.Println("connect server failed, error:", err)
		os.Exit(-1)
	}
	regitsterHandler(client)

	// go send_heartbeat()

	select {}

	fmt.Println("Stop PilotGo agent.")
}

func send_heartbeat(client *network.SocketClient) {
	for {
		msg := &protocol.Message{
			UUID: uuid.New().String(),
			Type: protocol.Heartbeat,
			Data: map[string]string{
				"agent_version": agent_version,
				"agent_id":      agent_uuid,
			},
		}

		if err := client.Send(msg); err != nil {
			fmt.Println("send message failed, error:", err)
		}
		fmt.Println("send heartbeat message")

		time.Sleep(time.Second)

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
		fmt.Println(msg.String())
		return nil
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
		servicerestart := uos.RestartService(service)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   servicerestart,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.ServiceStart, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())
		service := msg.Data.(string)
		servicestart := uos.StartService(service)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   servicestart,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.ServiceStop, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())
		service := msg.Data.(string)
		servicestop := uos.StopService(service)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   servicestop,
		}
		return c.Send(resp_msg)
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
		rpminfo, _ := uos.GetRpmInfo(rpmname)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   rpminfo,
		}
		return c.Send(resp_msg)
	})
	c.BindHandler(protocol.InstallRpm, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process agent info command:", msg.String())
		rpmname := msg.Data.(string)
		rpminstall := uos.InstallRpm(rpmname)

		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: 0,
			Data:   rpminstall,
		}
		return c.Send(resp_msg)
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
}
