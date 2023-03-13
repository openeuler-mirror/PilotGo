package dao

import (
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func IsExistName(name string) (bool, error) {
	var batch model.Batch
	err := global.PILOTGO_DB.Where("name=?", name).Find(&batch).Error
	return batch.ID != 0, err
}
func IsExistID(id int) (bool, error) {
	var batch model.Batch
	err := global.PILOTGO_DB.Where("id=?", id).Find(&batch).Error
	return batch.ID != 0, err
}
func GetBatchID(name string) (uint, error) {
	var batch model.Batch
	err := global.PILOTGO_DB.Where("name=?", name).Find(&batch).Error
	return batch.ID, err
}

func DeleteBatch(departid int) error {
	var batch model.Batch
	return global.PILOTGO_DB.Where("id=?", departid).Unscoped().Delete(&batch).Error
}

func UpdateBatch(BatchID int, BatchName string, Descrip string) error {
	var Batch model.Batch
	BatchNew := model.Batch{
		Name:        BatchName,
		Description: Descrip,
	}
	return global.PILOTGO_DB.Model(&Batch).Where("id=?", BatchID).Updates(&BatchNew).Error
}

func GetMachineID(BatchID int) ([]string, error) {
	var Batch model.Batch
	err := global.PILOTGO_DB.Where("id=?", BatchID).Find(&Batch).Error
	str := strings.Split(Batch.Machinelist, ",")
	return str, err
}

// 创建批次
func CreateBatch(batch model.Batch) error {
	return global.PILOTGO_DB.Create(&batch).Error
}

// 根据批次id获取所属的所有uuids
func BatchIds2UUIDs(batchIds []int) (uuids []string) {
	for _, batchId := range batchIds {
		var batch model.Batch
		err := global.PILOTGO_DB.Where("id=?", batchId).Find(&batch).Error
		if err != nil {
			logger.Error(err.Error())
		}
		str := strings.Split(batch.Machinelist, ",")
		macIds := utils.String2Int(str)

		for _, macId := range macIds {
			var machine MachineNode
			err = global.PILOTGO_DB.Where("id=?", macId).Find(&machine).Error
			if err != nil {
				logger.Error(err.Error())
			}
			uuids = append(uuids, machine.MachineUUID)
		}
	}
	return
}

func GetBatch() ([]model.Batch, error) {
	var batch []model.Batch
	err := global.PILOTGO_DB.Find(&batch).Error
	return batch, err
}
