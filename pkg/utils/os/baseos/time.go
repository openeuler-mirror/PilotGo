package baseos

import (
	"fmt"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

// 获取机器时间
func (b *BaseOS) GetTime() (string, error) {
	exitc, nowtime, stde, err := utils.RunCommand("date +%s")
	if exitc == 0 && nowtime != "" && stde == "" && err == nil {
		return nowtime, nil
	}
	return "", fmt.Errorf("failed to get unix time: %d, %s, %s, %v", exitc, nowtime, stde, err)
}
