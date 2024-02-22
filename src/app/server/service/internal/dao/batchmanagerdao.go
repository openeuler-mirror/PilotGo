package dao

import (
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gorm.io/gorm"
)

type Batch struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name" binding:"required" msg:"批次名不能为空"`
	Description string `gorm:"type:varchar(100)" json:"description"`
	Manager     string `gorm:"type:varchar(100)" json:"manager" binding:"required" msg:"创建人不能为空"`
	Depart      string `gorm:"type:varchar(100)"`
	DepartName  string `gorm:"type:varchar(100)"`
}

type BatchMachines struct {
	Batch         Batch `gorm:"Foreignkey:BatchID"`
	BatchID       uint
	MachineNode   MachineNode `gorm:"Foreignkey:MachineNodeID"`
	MachineNodeID uint
}

// 创建批次
func CreateBatchMessage(batch *Batch) error {
	return mysqlmanager.MySQL().Create(&batch).Error
}

func GetBatch() ([]Batch, error) {
	var batch []Batch
	err := mysqlmanager.MySQL().Find(&batch).Error
	return batch, err
}

// 分页查询所有批次
func GetBatchrPaged(offset, size int) (int64, []Batch, error) {
	var batchs []Batch
	var count int64
	err := mysqlmanager.MySQL().Model(Batch{}).Order("id desc").Offset(offset).Limit(size).Find(&batchs).Offset(-1).Limit(-1).Count(&count).Error
	return count, batchs, err
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

func DeleteBatch(departid int) error {
	var batch Batch
	return mysqlmanager.MySQL().Where("id=?", departid).Delete(&batch).Error
}

func UpdateBatch(BatchID int, BatchName string, Descrip string) error {
	var batch Batch
	BatchNew := Batch{
		Name:        BatchName,
		Description: Descrip,
	}
	return mysqlmanager.MySQL().Model(&batch).Where("id=?", BatchID).Updates(&BatchNew).Error
}

// 添加批次-机器数据
func AddBatch2Machine(b2m BatchMachines) error {
	return mysqlmanager.MySQL().Create(&b2m).Error
}

// 分页查询机器id
func GetMachineIDPaged(offset, size, batchid int) (int64, []BatchMachines, error) {
	var machineids []BatchMachines
	var count int64
	err := mysqlmanager.MySQL().Model(BatchMachines{}).Order("machine_node_id desc").Where("batch_id=?", batchid).Offset(offset).Limit(size).Find(&machineids).Offset(-1).Limit(-1).Count(&count).Error
	return count, machineids, err
}

func GetMachineID(BatchID int) ([]uint, error) {
	var machineids []uint
	err := mysqlmanager.MySQL().Model(BatchMachines{}).Select("machine_node_id").Where("batch_id=?", BatchID).Find(&machineids).Error
	return machineids, err
}
