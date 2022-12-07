package dao

import (
	"fmt"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

func AddConfigFile(cf model.ConfigFile) error {
	if IsUUIDExist(cf.MachineUUID) {
		global.PILOTGO_DB.Save(&cf)
		return nil
	}
	return fmt.Errorf("机器不存在")
}
