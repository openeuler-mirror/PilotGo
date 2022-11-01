package os

import (
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

//获取机器时间
func GetTime() string {
	nowtime := ""
	nowtime, _ = utils.RunCommand("date +%s")
	logger.Debug("nowtime:%s", nowtime)
	return nowtime
}
