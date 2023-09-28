package baseos

import (
	"fmt"
	"net"
	"strings"

	aconfig "gitee.com/PilotGo/PilotGo/app/agent/config"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/utils"
	"gitee.com/PilotGo/PilotGo/utils/os/common"
	"github.com/shirou/gopsutil/host"
)

func (b *BaseOS) GetHostInfo() (*common.SystemInfo, error) {
	//获取IP
	conn, err := net.Dial("udp", aconfig.Config().Server.Addr)
	if err != nil {
		fmt.Println("failed to get IP")
	}
	defer conn.Close()
	IP := strings.Split(conn.LocalAddr().String(), ":")[0]
	SysInfo, _ := host.Info()
	exitc, uptime, stde, err := utils.RunCommand(fmt.Sprintf("date -d '%v second ago'", SysInfo.Uptime))
	if exitc == 0 && uptime != "" && stde == "" && err == nil {
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
		return sysinfo, nil
	}
	logger.Error("failed to get host info: %d, %s, %s, %v", exitc, uptime, stde, err)
	return nil, fmt.Errorf("failed to get host info: %d, %s, %s, %v", exitc, uptime, stde, err)
}
