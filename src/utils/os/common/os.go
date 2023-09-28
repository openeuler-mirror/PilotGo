package common

import (
	"strings"

	"gitee.com/PilotGo/PilotGo/utils"
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
