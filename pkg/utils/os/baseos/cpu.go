package baseos

import (
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取CPU型号
func (b *BaseOS) GetCPUName() string {
	cpuname, _ := utils.RunCommand("lscpu | grep '型号名称'| sort| uniq")
	if cpuname == "" {
		cpuname, _ = utils.RunCommand("lscpu | grep 'Model name'| sort| uniq")
		if cpuname == "" {
			logger.Error("获取cpu型号失败!")
		}
	}

	cpuname = strings.Replace(cpuname, "\n", "", -1)
	str := strings.Split(cpuname, "：")
	if len(str) == 1 {
		str = strings.Split(cpuname, ":")
		cpuname = strings.TrimLeft(str[1], " ")
	} else {
		cpuname = strings.TrimLeft(str[1], " ")
	}
	return cpuname
}

// 获取物理CPU个数
func (b *BaseOS) GetPhysicalCPU() int {
	num, _ := utils.RunCommand("cat /proc/cpuinfo| grep 'processor'| sort| uniq| wc -l")
	num = strings.Replace(num, "\n", "", -1)
	cpunum, err := strconv.Atoi(num)
	if err != nil {
		logger.Error("获取cpu个数失败!")
	}
	return cpunum
}

func (b *BaseOS) GetCPUInfo() *common.CPUInfo {
	cpuinfo := &common.CPUInfo{
		ModelName: b.GetCPUName(),
		CpuNum:    b.GetPhysicalCPU(),
	}
	return cpuinfo
}
