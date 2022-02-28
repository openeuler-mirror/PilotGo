package dao

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

/**
 * @Author: wang hao
 * @Date: 2021/12/28 9:34
 * @Description:
 */
func IsParentDepartExist(parent string) bool {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("depart=? ", parent).Find(&Depart)
	return Depart.ID != 0
}
func IsDepartNodeExist(parent string, depart string) bool {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("depart=? and parent_depart=?", depart, parent).Find(&Depart)
	// mysqlmanager.DB.Where("", parent).Find(&Depart)
	return Depart.ID != 0
}
func IsDepartIDExist(ID int) bool {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("id=?", ID).Find(&Depart)
	return Depart.ID != 0
}
func DepartStore() []model.DepartNode {
	var Depart []model.DepartNode
	mysqlmanager.DB.Find(&Depart)
	logger.Info("%v", Depart)
	return Depart
}
func IsRootExist() bool {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("node_locate=?", 0).Find(&Depart)
	return Depart.ID != 0
}
func IsUUIDExist(uuid string) bool {
	var Machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&Machine)
	return Machine.DepartId != 0
}
func IsIPExist(ip string) bool {
	var Machine model.MachineNode
	mysqlmanager.DB.Where("ip=?", ip).Find(&Machine)
	return Machine.DepartId != 0
}
func Deleteuuid(uuid string) {
	var Machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Delete(Machine)
}
func MachineStore(departid int) []model.MachineNode {
	var Machineinfo []model.MachineNode
	mysqlmanager.DB.Where("depart_id=?", departid).Find(&Machineinfo)
	logger.Info("%v", Machineinfo)
	return Machineinfo
}

func GetPid(departid string) []model.DepartNode {
	var DepartInfo []model.DepartNode
	mysqlmanager.DB.Where("p_id=?", departid).Find(&DepartInfo)
	// logger.Info("%v", DepartInfo)
	return DepartInfo
}

func Deletedepartdata(needdelete []int) {
	var DepartInfo []model.DepartNode
	mysqlmanager.DB.Where("id=?", needdelete[0]).Delete(&DepartInfo)
}

//向需要删除的depart的组内增加需要删除的子节点
func Insertdepartlist(needdelete []int, str string) []int {
	var DepartInfo []model.DepartNode

	// needdelete = append(needdelete[:0], needdelete[1:]...)
	mysqlmanager.DB.Where("p_id=?", str).Find(&DepartInfo)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}
	return needdelete
}

func UpdateDepart(DepartID int, DepartName string) {
	var DepartInfo model.DepartNode
	Depart := model.DepartNode{
		Depart: DepartName,
	}
	mysqlmanager.DB.Model(&DepartInfo).Where("id=?", DepartID).Update(&Depart)
}
func UpdateParentDepart(DepartID int, DepartName string) {
	var DepartInfo model.DepartNode
	Depart := model.DepartNode{
		ParentDepart: DepartName,
	}
	mysqlmanager.DB.Model(&DepartInfo).Where("p_id=?", DepartID).Update(&Depart)
}
func MachineData(MachineID int) model.MachineNode {
	var m model.MachineNode
	mysqlmanager.DB.Where("id=?", MachineID).Find(&m)
	return m
}
