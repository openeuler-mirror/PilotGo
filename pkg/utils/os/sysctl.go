package os

import (
	"bufio"
	"fmt"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

// sysctl -p /etc/sysctl.conf 载入配置文件
func GetSysConfig() []map[string]string {
	tmp, err := utils.RunCommand("sysctl -p /etc/sysctl.conf")
	if err != nil {
		logger.Error("获取内核配置文件失败!%s", err.Error())
	}
	reader := strings.NewReader(tmp)
	scanner := bufio.NewScanner(reader)

	var sysConfig []map[string]string
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		strSlice := strings.Split(line, " = ")

		sysPars := map[string]string{
			strSlice[0]: strSlice[1],
		}
		sysConfig = append(sysConfig, sysPars)
	}
	return sysConfig
}

// sysctl -w net.ipv4.ip_forward=1  临时修改系统参数
func TempModifyPar(arg string) {
	tmp, err := utils.RunCommand(fmt.Sprintf("sudo sysctl -w %s", arg))
	if err != nil {
		logger.Error("修改内核运行时参数失败!%s", err.Error())
	}
	logger.Info("修改内核运行时参数成功!%s", tmp)
}

// sysctl -n net.ipv4.ip_forward  查看某个内核参数的值
func GetVarNameValue(arg string) {
	tmp, err := utils.RunCommand(fmt.Sprintf("sysctl -n %s", arg))
	if err != nil {
		logger.Error("获取该参数的值失败!%s", err.Error())
	}
	logger.Info("已将该参数修改!%s", tmp)
}
