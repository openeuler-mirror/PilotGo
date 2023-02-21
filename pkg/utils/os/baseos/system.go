package baseos

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/host"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func (b *BaseOS) GetHostInfo() *common.SystemInfo {
	//获取IP
	IP, err := utils.RunCommand("hostname -I")
	if err != nil {
		fmt.Println("获取IP失败!")
	}
	str := strings.Split(IP, " ")
	IP = str[0]
	SysInfo, _ := host.Info()
	Uptime := fmt.Sprintf("date -d '%v second ago'", SysInfo.Uptime)
	uptime, _ := utils.RunCommand(Uptime)
	uptime = strings.Replace(uptime, "\n", "", -1)
	sysinfo := &common.SystemInfo{
		IP:              IP,
		Platform:        SysInfo.Platform,
		PlatformVersion: SysInfo.PlatformVersion,
		KernelVersion:   SysInfo.KernelVersion,
		KernelArch:      SysInfo.KernelArch,
		HostId:          SysInfo.HostID,
		Uptime:          uptime,
	}
	return sysinfo
}
