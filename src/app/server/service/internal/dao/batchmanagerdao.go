package dao

import (
	"strings"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
	"gorm.io/gorm"
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

type Batch2Machine struct {
	Batch         Batch `gorm:"Foreignkey:BatchID"`
	BatchID       uint
	MachineNode   MachineNode `gorm:"Foreignkey:MachineNodeID"`
	MachineNodeID uint
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

func GetBatchName(id int) (string, error) {
	var batch Batch
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&batch).Error
	return batch.Name, err
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
func CreateBatchMessage(batch *Batch) error {
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

// 分页查询所有用户
func GetBatchrPaged(offset, size int) (int64, []Batch, error) {
	var batchs []Batch
	var count int64
	err := mysqlmanager.MySQL().Model(Batch{}).Order("id desc").Offset(offset).Limit(size).Find(&batchs).Offset(-1).Limit(-1).Count(&count).Error
	return count, batchs, err
}

// 添加机器批次数据
func AddBatch2Machine(b2m Batch2Machine) error {
	return mysqlmanager.MySQL().Create(&b2m).Error
}
