package baseos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

// 通过 /proc/cpuinfo来获取CPU型号
type CPUInfo struct {
	ModelName string
	CpuNum    int
}

func (cpu *CPUInfo) String() string {
	b, err := json.Marshal(*cpu)
	if err != nil {
		return fmt.Sprintf("%+v", *cpu)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *cpu)
	}
	return out.String()
}

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

func (b *BaseOS) GetCPUInfo() *CPUInfo {
	cpuinfo := &CPUInfo{
		ModelName: b.GetCPUName(),
		CpuNum:    b.GetPhysicalCPU(),
	}
	return cpuinfo
}
