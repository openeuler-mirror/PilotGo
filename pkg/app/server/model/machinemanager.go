package model

import "openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"

/**
 * @Author: wang hao
 * @Date: 2021/12/23 17:00
 * @Description:机器管理树形结构体
 */

type DepartNode struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID          int    `gorm:"type:int(100);not null" json:"pid"`
	ParentDepart string `gorm:"type:varchar(100);not null" json:"parentdepart"`
	Depart       string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate   int    `gorm:"type:int(100);not null" json:"nodelocate"`
	//根节点为0,普通节点为1
}
type MachineNode struct {
	DepartId    int    `gorm:"type:int(100);not null" json:"departid"`
	IP          string `gorm:"type:varchar(100)" json:"ip"`
	MachineUUID string `gorm:"type:varchar(100);not null" json:"machineuuid"`
	CPU         string `gorm:"type:varchar(100)" json:"CPU"`
	State       string `gorm:"type:varchar(100)" json:"state"`
	Systeminfo  string `gorm:"type:varchar(100)" json:"sysinfo"`
	BatchInfo   string `gorm:"type:varchar(100)"`
} // kylin/serve/opensource/ops/wanghao/dfsagdasgs

type MachineTreeNode struct {
	Label    string             `json:"label"`
	Id       int                `json:"id"`
	Pid      int                `json:"pid"`
	Children []*MachineTreeNode `json:"children"`
}

type MachineInfo struct {
	Uuid []string `json:"label"`
}

func (m *MachineNode) ReturnMachine(q *PaginationQ, departid int) (list *[]MachineNode, total uint, err error) {
	list = &[]MachineNode{}
	tx := mysqlmanager.DB.Where("depart_id=?", departid).Find(&list)
	total, err = CrudAll(q, tx, list)
	return
}
