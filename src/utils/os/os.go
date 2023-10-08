package os

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/os/centos"
	"gitee.com/openeuler/PilotGo/utils/os/common"
	"gitee.com/openeuler/PilotGo/utils/os/kylin"
	"gitee.com/openeuler/PilotGo/utils/os/nestos"
	"gitee.com/openeuler/PilotGo/utils/os/openeuler"
)

const (
	OpenEuler = "openEuler"
	Kylin     = "kylin"
	NestOS    = "NestOS"
	CentOS    = "CentOS"
)

func OS() common.OSOperator {
	osinfo, err := common.InitOSName()
	if err != nil {
		logger.Error("osname init failed: %s", err)
		return nil
	}
	switch osinfo.OSName {
	case OpenEuler:
		return &openeuler.OpenEuler{}
	case Kylin:
		return &kylin.Kylin{}
	case NestOS:
		switch osinfo.ID {
		case "NestOS For Container":
			return &nestos.NestOS4Container{}
		case "NestOS For Virt":
			return &nestos.NestOS4Virt{}
		}
	case CentOS:
		return &centos.CentOS{}
	}
	return nil
}
