package os

import (
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/baseos"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/openeuler"
)

const (
	OpenEuler = "openEuler"
	Kylin     = "kylin"
)

func detectOS() string {
	// TODO:
	return OpenEuler
}

func OS() common.OSOperator {
	switch detectOS() {
	case OpenEuler:
		return &openeuler.OpenEuler{}
	case Kylin:
		return &baseos.BaseOS{}
	}
	return nil
}
