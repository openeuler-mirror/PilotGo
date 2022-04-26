/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: wanghao
 * Date: 2022-02-18 10:25:44
 * LastEditTime: 2022-02-24 14:31:19
 * Description: provide machine manager functions.
 ******************************************************************************/
package model

import (
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

type DepartNode struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID          int    `gorm:"type:int(100);not null" json:"pid"`
	ParentDepart string `gorm:"type:varchar(100);not null" json:"parentdepart"`
	Depart       string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate   int    `gorm:"type:int(100);not null" json:"nodelocate"`
	//根节点为0,普通节点为1
}
type MachineNode struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	DepartId    int    `gorm:"type:int(100);not null" json:"departid"`
	IP          string `gorm:"type:varchar(100)" json:"ip"`
	MachineUUID string `gorm:"type:varchar(100);not null" json:"machineuuid"`
	CPU         string `gorm:"type:varchar(100)" json:"CPU"`
	State       int    `gorm:"type:varchar(100)" json:"state"`
	Systeminfo  string `gorm:"type:varchar(100)" json:"sysinfo"`
}

type MachineTreeNode struct {
	Label    string             `json:"label"`
	Id       int                `json:"id"`
	Pid      int                `json:"pid"`
	Children []*MachineTreeNode `json:"children"`
}

func GetMachList() (mi []MachineNode, e error) {
	mysqlmanager.DB.Find(&mi)
	return mi, nil
}

const (
	// 机器运行
	Normal = 1
	// 脱机
	OffLine = 2
	// 空闲
	Free = 3
)

type MachineInfo struct {
	RPM        string `json:"rpm"`
	Service    string `json:"service"`
	SysctlArgs string `json:"args"`
	MountPath  string `json:"mountpath"`
	SourceDisk string `json:"source"`
	DestPath   string `json:"dest"`
	File       string `json:"file"`
	FileType   string `json:"type"`
	DiskPath   string `json:"path"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	Permission string `json:"per"`
}

type Res struct {
	ID         int    `json:"id"`
	Departid   int    `json:"departid"`
	Departname string `json:"departname"`
	IP         string `json:"ip"`
	UUID       string `json:"uuid"`
	CPU        string `json:"cpu"`
	State      string `json:"state"`
	Systeminfo string `json:"systeminfo"`
}

func (m *MachineNode) ReturnMachine(q *PaginationQ, departid int) (list *[]Res, tx *gorm.DB, res []Res) {
	list = &[]Res{}
	// tx := mysqlmanager.DB.Where("depart_id=?", departid).Find(&list)
	tx = mysqlmanager.DB.Table("machine_node").Where("depart_id=?", departid).Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
		"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&list)
	logger.Debug("%+v", list)
	res = make([]Res, 0)
	for _, value := range *list {
		if value.Departid == departid {
			res = append(res, value)
		}
	}
	return
}
func MachineAllData() []Res {
	var mch []Res
	mysqlmanager.DB.Table("machine_node").Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
		"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&mch)
	return mch
}

type DeleteDepart struct {
	DepartID int `json:"DepartID"`
}
type Depart struct {
	Page       int  `form:"page"`
	Size       int  `form:"size"`
	ID         int  `form:"DepartId"`
	ShowSelect bool `form:"ShowSelect"`
}
