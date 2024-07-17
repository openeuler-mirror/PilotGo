package os

import (
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/os/centos"
	"gitee.com/openeuler/PilotGo/utils/os/common"
	"gitee.com/openeuler/PilotGo/utils/os/kylin"
	"gitee.com/openeuler/PilotGo/utils/os/nestos"
	"gitee.com/openeuler/PilotGo/utils/os/openeuler"
)

const (
	OpenEuler = "openeuler"
	Kylin     = "kylin"
	NestOS    = "nestos"
	CentOS    = "centos"
)

func OS() common.OSOperator {
	osinfo, err := common.InitOSName()
	if err != nil {
		logger.Error("osname init failed: %s", err)
		return nil
	}
	switch strings.ToLower(osinfo.OSName) {
	case OpenEuler:
		return &openeuler.OpenEuler{}
	case Kylin:
		return &kylin.Kylin{}
	case NestOS:
		switch strings.ToLower(osinfo.ID) {
		case "nestos for container":
			return &nestos.NestOS4Container{}
		case "nestos for virt":
			return &nestos.NestOS4Virt{}
		}
	case CentOS:
		return &centos.CentOS{}
	}
	return nil
}
