package os

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

// 通过 /proc/cpuinfo来获取CPU型号
type CPUInfo struct {
	CpuName string
	CpuNum  int
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
func GetCPUName() string {
	cpuname, _ := utils.RunCommand("cat /proc/cpuinfo | grep name | sort | uniq")
	cpuname = strings.Replace(cpuname, "\n", "", -1)
	str := len("model name	: ")
	cpuname = cpuname[str:]
	return cpuname
}

// 获取物理CPU个数
func GetPhysicalCPU() int {
	num, _ := utils.RunCommand("cat /proc/cpuinfo| grep 'physical id'| sort| uniq| wc -l")
	num = strings.Replace(num, "\n", "", -1)
	cpunum, err := strconv.Atoi(num)
	if err != nil {
		logger.Error("获取cpu个数失败!")
	}
	return cpunum
}
