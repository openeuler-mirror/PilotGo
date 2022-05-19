package dao

import (
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

func IsExistName(name string) bool {
	var batch model.Batch
	mysqlmanager.DB.Where("name=?", name).Find(&batch)
	return batch.ID != 0
}
func IsExistID(id int) bool {
	var batch model.Batch
	mysqlmanager.DB.Where("id=?", id).Find(&batch)
	return batch.ID != 0
}
func GetBatchID(name string) uint {
	var batch model.Batch
	mysqlmanager.DB.Where("name=?", name).Find(&batch)
	return batch.ID
}

func DeleteBatch(departid int) {
	var batch model.Batch
	mysqlmanager.DB.Where("id=?", departid).Unscoped().Delete(&batch)
}

func UpdateBatch(BatchID int, BatchName string, Descrip string) {
	var Batch model.Batch
	BatchNew := model.Batch{
		Name:        BatchName,
		Description: Descrip,
	}
	mysqlmanager.DB.Model(&Batch).Where("id=?", BatchID).Update(&BatchNew)
}

func GetMachineID(BatchID int) []string {
	var Batch model.Batch
	mysqlmanager.DB.Where("id=?", BatchID).Find(&Batch)
	str := strings.Split(Batch.Machinelist, ",")
	return str
}

// 创建批次
func CreateBatch(batch model.Batch) {
	mysqlmanager.DB.Create(&batch)
}
