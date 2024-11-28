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
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func (b *BaseOS) GetSysctlConfig() (map[string]string, error) {
	exitc, tmp, stde, err := utils.RunCommand("sysctl -a")
	if exitc == 0 && tmp != "" && stde == "" && err == nil {
		sysConfig := make(map[string]string)
		lines := strings.Split(tmp, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			strSlice := strings.Split(line, " =")
			key := strSlice[0]
			value := strings.TrimLeft(line[len(key)+2:], " ")
			sysConfig[key] = value
		}
		return sysConfig, nil
	}
	logger.Error("failed to retrieve the kernel configuration file: %d, %s, %s, %v", exitc, tmp, stde, err)
	return nil, fmt.Errorf("failed to retrieve the kernel configuration file: %d, %s, %s, %v", exitc, tmp, stde, err)

}

// sysctl -w net.ipv4.ip_forward=1  临时修改系统参数
func (b *BaseOS) TempModifyPar(arg string) (string, error) {
	exitc, tmp, stde, err := utils.RunCommand(fmt.Sprintf("sudo sysctl -w %s", arg))
	if exitc == 0 && tmp != "" && stde == "" && err == nil {
		tmp = strings.Replace(tmp, "\n", "", -1)
		return tmp, nil
	}
	logger.Error("failed to modify the kernel runtime parameters: %d, %s, %s, %v", exitc, tmp, stde, err)
	return "", fmt.Errorf("failed to modify the kernel runtime parameters: %d, %s, %s, %v", exitc, tmp, stde, err)

}

// sysctl -n net.ipv4.ip_forward  查看某个内核参数的值
func (b *BaseOS) GetVarNameValue(arg string) (string, error) {
	exitc, tmp, stde, err := utils.RunCommand(fmt.Sprintf("sysctl -n %s", arg))
	if exitc == 0 && tmp != "" && stde == "" && err == nil {
		tmp = strings.Replace(tmp, "\n", "", -1)
		return tmp, nil
	}
	logger.Error("failed to get the value of the parameter: %d, %s, %s, %v", exitc, tmp, stde, err)
	return "", fmt.Errorf("failed to get the value of the parameter: %d, %s, %s, %v", exitc, tmp, stde, err)

}
