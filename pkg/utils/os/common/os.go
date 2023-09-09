package common

import (
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type OSInfo struct {
	OSName string
	ID     string
}

// 从系统文件中读取os名字和id
func InitOSName() (osinfo OSInfo, err error) {
	contents, err := utils.FileReadString("/etc/os-release")
	if err != nil {
		return osinfo, err
	}
	for _, line := range strings.Split(contents, "\n") {
		field := strings.Split(line, "=")
		if len(field) < 2 {
			continue
		}
		switch field[0] {
		case "NAME":
			osinfo.OSName = trimQuotes(field[1])
		case "ID":
			osinfo.ID = trimQuotes(field[1])
		}
	}
	return osinfo, nil
}

// 剪掉引号
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
