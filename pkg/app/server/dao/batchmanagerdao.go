package dao

import (
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func IsExistName(name string) bool {
	var batch model.Batch
	global.PILOTGO_DB.Where("name=?", name).Find(&batch)
	return batch.ID != 0
}
func IsExistID(id int) bool {
	var batch model.Batch
	global.PILOTGO_DB.Where("id=?", id).Find(&batch)
	return batch.ID != 0
}
func GetBatchID(name string) uint {
	var batch model.Batch
	global.PILOTGO_DB.Where("name=?", name).Find(&batch)
	return batch.ID
}

func DeleteBatch(departid int) {
	var batch model.Batch
	global.PILOTGO_DB.Where("id=?", departid).Unscoped().Delete(&batch)
}

func UpdateBatch(BatchID int, BatchName string, Descrip string) {
	var Batch model.Batch
	BatchNew := model.Batch{
		Name:        BatchName,
		Description: Descrip,
	}
	global.PILOTGO_DB.Model(&Batch).Where("id=?", BatchID).Updates(&BatchNew)
}

func GetMachineID(BatchID int) []string {
	var Batch model.Batch
	global.PILOTGO_DB.Where("id=?", BatchID).Find(&Batch)
	str := strings.Split(Batch.Machinelist, ",")
	return str
}

// 创建批次
func CreateBatch(batch model.Batch) {
	global.PILOTGO_DB.Create(&batch)
}

// 根据批次id获取所属的所有uuids
func BatchIds2UUIDs(batchIds []int) (uuids []string) {
	for _, batchId := range batchIds {
		var batch model.Batch
		global.PILOTGO_DB.Where("id=?", batchId).Find(&batch)

		str := strings.Split(batch.Machinelist, ",")
		macIds := utils.String2Int(str)

		for _, macId := range macIds {
			var machine model.MachineNode
			global.PILOTGO_DB.Where("id=?", macId).Find(&machine)
			uuids = append(uuids, machine.MachineUUID)
		}
	}
	return
}

func GetBatch() []model.Batch {
	var batch []model.Batch
	global.PILOTGO_DB.Find(&batch)
	return batch
}
