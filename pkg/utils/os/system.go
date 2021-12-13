package os

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/host"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

type SystemInfo struct {
	Platform        string //系统平台
	PlatformVersion string //系统版本
	KernelVersion   string //内核版本
	KernelArch      string //内核支持架构
	HostId          string //系统id
	Uptime          string //系统最新启动时间
}

func GetHostInfo() string {
	SysInfo, _ := host.Info()
	Uptime := fmt.Sprintf("date -d '%v second ago'", SysInfo.Uptime)
	uptime, _ := utils.RunCommand(Uptime)
	uptime = strings.Replace(uptime, "\n", "", -1)
	sysinfo := SystemInfo{
		Platform:        SysInfo.Platform,
		PlatformVersion: SysInfo.PlatformVersion,
		KernelVersion:   SysInfo.KernelVersion,
		KernelArch:      SysInfo.KernelArch,
		HostId:          SysInfo.HostID,
		Uptime:          uptime,
	}
	system, _ := json.Marshal(sysinfo)
	return string(system)
}
