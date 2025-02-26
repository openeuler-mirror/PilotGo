/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package baseos

import (
	"errors"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"github.com/duke-git/lancet/datetime"
	"github.com/duke-git/lancet/fileutil"
	"github.com/shirou/gopsutil/v3/host"
)

func (b *BaseOS) GetHostInfo() (*common.SystemInfo, error) {
	IP, err := b.GetHostIp()
	if err != nil {
		return nil, errors.New("get host ip failed:" + err.Error())
	}

	hostInfo, err := host.Info()
	if err != nil {
		return nil, errors.New("get host info failed:" + err.Error())
	}

	osReleaseInfo, err := b.OSReleaseInfo()
	if err != nil {
		return nil, errors.New("get os-release failed:" + err.Error())
	}

	bootTime := time.Unix(int64(hostInfo.BootTime), 0)
	bootTimeStr := datetime.FormatTimeToStr(bootTime, "yyyy-mm-dd hh:mm:ss")
	sysinfo := &common.SystemInfo{
		IP:              IP,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		PrettyName:      osReleaseInfo.PrettyName,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		HostId:          hostInfo.HostID,
		Uptime:          bootTimeStr,
	}
	return sysinfo, nil
}

// 读取os-release文件信息
func (b *BaseOS) OSReleaseInfo() (*common.OSReleaseInfo, error) {
	lines, err := fileutil.ReadFileByLine("/etc/os-release")
	if err != nil {
		return nil, err
	}

	info := &common.OSReleaseInfo{}
	for _, line := range lines {
		if line == "" {
			continue
		}

		words := strings.Split(line, "=")
		if len(words) == 2 {
			k := words[0]
			v := strings.Trim(words[1], "\"")

			switch k {
			case "NAME":
				info.Name = v
			case "VERSION":
				info.Version = v
			case "ID":
				info.ID = v
			case "VERSION_ID":
				info.VersionID = v
			case "PRETTY_NAME":
				info.PrettyName = v
			}
		} else {
			return nil, errors.New("invalid os-release format:" + line)
		}
	}

	return info, nil
}
