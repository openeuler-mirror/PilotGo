package os

import (
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/kylin"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/nestos"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/openeuler"
)

const (
	OpenEuler = "openEuler"
	Kylin     = "kylin"
	NestOS    = "NestOS"
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
			return &nestos.NestOS{}
		case "NestOS For Virt":
			return &openeuler.OpenEuler{}
		}
	}
	return nil
}
