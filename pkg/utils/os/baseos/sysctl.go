package baseos

import (
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func (b *BaseOS) GetSysctlConfig() (map[string]string, error) {
	tmp, err := utils.RunCommand("sysctl -a")
	if err != nil {
		logger.Error("Failed to retrieve the kernel configuration file: %s", err.Error())
		return nil, err
	}

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

// sysctl -w net.ipv4.ip_forward=1  临时修改系统参数
func (b *BaseOS) TempModifyPar(arg string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("sudo sysctl -w %s", arg))
	tmp = strings.Replace(tmp, "\n", "", -1)

	if err != nil {
		logger.Error("failed to modify the kernel runtime parameters: %s", err.Error())
	}
	return tmp
}

// sysctl -n net.ipv4.ip_forward  查看某个内核参数的值
func (b *BaseOS) GetVarNameValue(arg string) string {
	tmp, err := utils.RunCommand(fmt.Sprintf("sysctl -n %s", arg))
	tmp = strings.Replace(tmp, "\n", "", -1)
	if err != nil {
		logger.Error("failed to get the value of the parameter: %s", err.Error())
	}
	return tmp
}
