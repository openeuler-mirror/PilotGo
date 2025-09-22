/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package baseos

import (
	"fmt"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// 获取CPU型号
func (b *BaseOS) GetCPUName() (string, error) {
	// 先检测当前架构
	arch, err := getSystemArchitecture()
	if err != nil {
		logger.Error("failed to detect system architecture: %v", err)
		return "", fmt.Errorf("failed to detect system architecture: %v", err)
	}

	// RISC-V 架构特殊处理
	if arch == "riscv64" {
		exitc, cpuname, stde, err := utils.RunCommand("cat /proc/device-tree/compatible")
		if exitc == 0 && len(cpuname) > 0 && stde == "" && err == nil {
			cpuname = strings.TrimSpace(cpuname)
			// 提取第一个兼容的设备名称
			compatibleList := strings.Split(cpuname, "\x00")
			if len(compatibleList) > 0 {
				return compatibleList[0], nil
			}
		}
		logger.Error("failed to get RISC-V CPU info: %d, %s, %s, %v", exitc, cpuname, stde, err)
		return "", fmt.Errorf("failed to get RISC-V CPU info: %d, %s, %s, %v", exitc, cpuname, stde, err)
	}

	// 其他架构保持原有逻辑
	exitc, cpuname, stde, err := utils.RunCommand("lscpu | grep 'Model name' | sort | uniq")
	if exitc == 0 && len(cpuname) > 0 && stde == "" && err == nil {
		cpuname = strings.Replace(cpuname, "\n", "", -1)
		str := strings.Split(cpuname, ":")
		if len(str) >= 2 {
			return strings.TrimSpace(str[1]), nil
		}
	}
	logger.Error("failed to get cpu model name: %d, %s, %s, %v", exitc, cpuname, stde, err)
	return "", fmt.Errorf("failed to get cpu model name: %d, %s, %s, %v", exitc, cpuname, stde, err)
}

// 获取系统架构
func getSystemArchitecture() (string, error) {
	exitc, arch, stde, err := utils.RunCommand("uname -m")
	if exitc != 0 || stde != "" || err != nil {
		return "", fmt.Errorf("failed to get architecture: %d, %s, %v", exitc, stde, err)
	}
	return strings.TrimSpace(arch), nil
}

// 获取物理CPU个数
func (b *BaseOS) GetPhysicalCPU() (int, error) {
	exitc, num, stde, err := utils.RunCommand("cat /proc/cpuinfo| grep 'processor'| sort| uniq| wc -l")
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
		return nil, fmt.Errorf(err.Error())
	}
	cpunum, err := b.GetPhysicalCPU()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	cpuinfo := &common.CPUInfo{
		ModelName: cpuname,
		CpuNum:    cpunum,
	}
	return cpuinfo, nil
}
