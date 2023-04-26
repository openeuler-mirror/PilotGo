package baseos

import (
	"fmt"
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取CPU型号
func (b *BaseOS) GetCPUName() (string, error) {
	exitc, cpuname, stde, err := utils.RunCommandnew("lscpu | grep 'Model name' | sort | uniq")
	if exitc == 0 && len(cpuname) > 0 && stde == "" && err == nil {
		cpuname = strings.Replace(cpuname, "\n", "", -1)
		str := strings.Split(cpuname, ":")
		if len(str) == 1 {
			str = strings.Split(cpuname, ":")
			cpuname = strings.TrimLeft(str[1], " ")
		} else {
			cpuname = strings.TrimLeft(str[1], " ")
		}
		return cpuname, nil
	}
	logger.Error("failed to get cpu model name: %d, %s, %s, %v", exitc, cpuname, stde, err)
	return "", fmt.Errorf("failed to get cpu model name: %d, %s, %s, %v", exitc, cpuname, stde, err)
}

// 获取物理CPU个数
func (b *BaseOS) GetPhysicalCPU() (int, error) {
	exitc, num, stde, err := utils.RunCommandnew("cat /proc/cpuinfo| grep 'processor'| sort| uniq| wc -l")
	if exitc == 0 && len(num) > 0 && stde == "" && err == nil {
		num = strings.Replace(num, "\n", "", -1)
		cpunum, erratoi := strconv.Atoi(num)
		if erratoi != nil {
			return -1, erratoi
		}
		return cpunum, nil
	}
	logger.Error("failed to get cpu num: %d, %s, %s, %v", exitc, num, stde, err)
	return -1, fmt.Errorf("failed to get cpu num: %d, %s, %s, %v", exitc, num, stde, err)
}

func (b *BaseOS) GetCPUInfo() (*common.CPUInfo, error) {
	cpuname, err := b.GetCPUName()
	if err != nil {
		return nil, fmt.Errorf("failed to get cpu model name")
	}
	cpunum, err := b.GetPhysicalCPU()
	if err != nil {
		return nil, fmt.Errorf("failed to get cpu num")
	}
	cpuinfo := &common.CPUInfo{
		ModelName: cpuname,
		CpuNum:    cpunum,
	}
	return cpuinfo, nil
}
