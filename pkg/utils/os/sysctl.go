package os

import (
	"fmt"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

func GetSysctlConfig() ([]map[string]string, error) {
	tmp, err := utils.RunCommand("sysctl -a")
	if err != nil {
		logger.Error("获取内核配置文件失败!%s", err.Error())
		return nil, err
	}
	// TODO: 修正数据结构
	var sysConfig []map[string]string
	lines := strings.Split(tmp, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		strSlice := strings.Split(line, " =")
		key := strSlice[0]
		value := strings.TrimLeft(line[len(key)+2:], " ")
		sysPars := map[string]string{
			key: value,
		}
		sysConfig = append(sysConfig, sysPars)
	}
	return sysConfig, nil
}

// sysctl -w net.ipv4.ip_forward=1  临时修改系统参数
func TempModifyPar(arg string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("sudo sysctl -w %s", arg))
	tmp = strings.Replace(tmp, "\n", "", -1)

	if err != nil {
		logger.Error("修改内核运行时参数失败!%s", err.Error())
	}
	return tmp
}

// sysctl -n net.ipv4.ip_forward  查看某个内核参数的值
func GetVarNameValue(arg string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("sysctl -n %s", arg))
	tmp = strings.Replace(tmp, "\n", "", -1)
	if err != nil {
		logger.Error("获取该参数的值失败!%s", err.Error())
	}
	return tmp
}
