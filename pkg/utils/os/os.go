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

type OSOperator interface {
	common.NetworkOperator
}

func OS() OSOperator {
	switch detectOS() {
	case OpenEuler:
		return &openeuler.OpenEuler{}
	case Kylin:
		return &baseos.BaseOS{}
	}
	return nil
}
