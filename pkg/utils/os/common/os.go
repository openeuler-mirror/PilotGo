/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Sat Sep 9 15:42:10 2023 +0800
 */
package common

import (
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils"
)

type OSInfo struct {
	OSName string
	ID     string
}

// 从系统文件中读取os名字和id
func InitOSName() (osinfo OSInfo, err error) {
	contents, err := utils.FileReadString("/etc/system-release")
	if err != nil {
		return osinfo, err
	}
	name := strings.Split(contents, " ")[0]
	switch name {
	case "NestOS-For-Container":
		osinfo.OSName = "NestOS"
		osinfo.ID = "NestOS For Container"
		return osinfo, nil
	case "NestOS-For-Virt":
		osinfo.OSName = "NestOS"
		osinfo.ID = "NestOS For Virt"
		return osinfo, nil
	default:
		osinfo.OSName = name
	}
	return osinfo, nil
}
