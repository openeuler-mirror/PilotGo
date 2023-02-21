package baseos

import (
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

// 获取机器时间
func (b *BaseOS) GetTime() (string, error) {
	nowtime, err := utils.RunCommand("date +%s")
	if err != nil {
		return "", err
	}
	return nowtime, nil
}
