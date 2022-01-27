package dao

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func IsExistName(name string) bool {
	var batch model.Batch
	mysqlmanager.DB.Where("name=?", name).Find(&batch)
	return batch.ID != 0
}
func GetBatchID(name string) uint {
	var machine model.Batch
	mysqlmanager.DB.Where("name=?", name).Find(&machine)
	logger.Info("%d", machine.ID)
	return machine.ID
}
func GetmachineBatch(uuid string) string {
	var machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machine)
	logger.Info("%s", machine.BatchInfo)
	return machine.BatchInfo
}

func UpdatemachineBatch(s string, b string) {
	var machine model.MachineNode
	mysqlmanager.DB.Find(&machine)
	mysqlmanager.DB.Model(&machine).Where("machine_uuid=?", s).Update("batch_info", b)
	logger.Info("%+v", machine)
}

// func UpdatemachineBatch(m model.MachineNode) {

// }
