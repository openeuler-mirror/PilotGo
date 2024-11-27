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
	"regexp"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
)

func moduleMatch(name string, value int64, memconf *common.MemoryConfig) {
	if name == "MemTotal" {
		memconf.MemTotal = value
	} else if name == "MemFree" {
		memconf.MemFree = value
	} else if name == "MemAvailable" {
		memconf.MemAvailable = value
	} else if name == "Buffers" {
		memconf.Buffers = value
	} else if name == "Cached" {
		memconf.Cached = value
	} else if name == "SwapCached" {
		memconf.SwapCached = value
	}
}

func (b *BaseOS) GetMemoryConfig() (*common.MemoryConfig, error) {
	exitc, output, stde, err := utils.RunCommand("cat /proc/meminfo")
	if exitc == 0 && output != "" && stde == "" && err == nil {
		outputlines := strings.Split(output, "\n")
		m := &common.MemoryConfig{}
		reg := regexp.MustCompile(`[a-zA-Z\s]+`)
		for _, line := range outputlines {
			//一次获取一行,_ 获取当前行是否被读完
			if line == "" {
				continue
			}
			k := strings.Split(line, ":")[0]
			v := reg.ReplaceAllString(strings.Split(line, ":")[1], "")
			vint64, _ := strconv.ParseInt(v, 10, 64)
			moduleMatch(k, vint64, m)
		}
		return m, nil
	}
	return nil, fmt.Errorf("failed to get memory config: %d, %s, %s, %v", exitc, output, stde, err)

}
