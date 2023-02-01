package dao

import (
	"fmt"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

func AddConfigFile(cf model.ConfigFile) error {
	UUIDExistbool, err := IsUUIDExist(cf.MachineUUID)
	if err != nil {
		return err
	}
	if UUIDExistbool {
		return global.PILOTGO_DB.Save(&cf).Error
	}
	return fmt.Errorf("机器不存在")
}
