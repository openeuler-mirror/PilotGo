package dao

import (
	"strings"

	"gorm.io/gorm"
	"openeuler.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type Batch struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(100)" json:"description"`
	Manager     string `gorm:"type:varchar(100)" json:"manager"`
	Machinelist string `json:"machinelist"`
	Depart      string `gorm:"type:varchar(100)"`
	DepartName  string `gorm:"type:varchar(100)"`
}

func (b *Batch) ReturnBatch() (list *[]Batch, tx *gorm.DB) {
	list = &[]Batch{}
	tx = mysqlmanager.MySQL().Order("created_at desc").Find(&list)
	return
}

func IsExistName(name string) (bool, error) {
	var batch Batch
	err := mysqlmanager.MySQL().Where("name=?", name).Find(&batch).Error
	return batch.ID != 0, err
}

func IsExistID(id int) (bool, error) {
	var batch Batch
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&batch).Error
	return batch.ID != 0, err
}

func GetBatchID(name string) (uint, error) {
	var batch Batch
	err := mysqlmanager.MySQL().Where("name=?", name).Find(&batch).Error
	return batch.ID, err
}

func DeleteBatch(departid int) error {
	var batch Batch
	return mysqlmanager.MySQL().Where("id=?", departid).Unscoped().Delete(&batch).Error
}

func UpdateBatch(BatchID int, BatchName string, Descrip string) error {
	var batch Batch
	BatchNew := Batch{
		Name:        BatchName,
		Description: Descrip,
	}
	return mysqlmanager.MySQL().Model(&batch).Where("id=?", BatchID).Updates(&BatchNew).Error
}

func GetMachineID(BatchID int) ([]string, error) {
	var Batch Batch
	err := mysqlmanager.MySQL().Where("id=?", BatchID).Find(&Batch).Error
	str := strings.Split(Batch.Machinelist, ",")
	return str, err
}

// 创建批次
func CreateBatchMessage(batch Batch) error {
	return mysqlmanager.MySQL().Create(&batch).Error
}

// 根据批次id获取所属的所有uuids
func BatchIds2UUIDs(batchIds []int) (uuids []string) {
	for _, batchId := range batchIds {
		var batch Batch
		err := mysqlmanager.MySQL().Where("id=?", batchId).Find(&batch).Error
		if err != nil {
			logger.Error(err.Error())
		}
		str := strings.Split(batch.Machinelist, ",")
		macIds := utils.String2Int(str)

		for _, macId := range macIds {
			var machine MachineNode
			err = mysqlmanager.MySQL().Where("id=?", macId).Find(&machine).Error
			if err != nil {
				logger.Error(err.Error())
			}
			uuids = append(uuids, machine.MachineUUID)
		}
	}
	return
}

func GetBatch() ([]Batch, error) {
	var batch []Batch
	err := mysqlmanager.MySQL().Find(&batch).Error
	return batch, err
}
