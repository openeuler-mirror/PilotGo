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
 * Date: 2022-01-04 12:56:59
 * LastEditTime: 2023-07-11 19:34:18
 * Description: 部门管理数据库相关函数
 ******************************************************************************/
package dao

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
)

type MachineNode struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	DepartId    int    `gorm:"type:int(100);not null" json:"departid"`
	IP          string `gorm:"type:varchar(100)" json:"ip"`
	MachineUUID string `gorm:"type:varchar(100);not null" json:"machineuuid"`
	CPU         string `gorm:"type:varchar(100)" json:"CPU"`
	State       int    `gorm:"type:varchar(100)" json:"state"`
	Systeminfo  string `gorm:"type:varchar(100)" json:"sysinfo"`
}

type Res struct {
	ID         int    `json:"id"`
	Departid   int    `json:"departid"`
	Departname string `json:"departname"`
	IP         string `json:"ip"`
	UUID       string `json:"uuid"`
	CPU        string `json:"cpu"`
	State      int    `json:"state"`
	Systeminfo string `json:"systeminfo"`
}

func ReturnMachinePaged(departid, offset, size int) (int64, []Res, error) {
	var count int64
	var machinelist []Res
	list := &[]Res{}
	err := mysqlmanager.MySQL().Model(MachineNode{}).Where("depart_id=?", departid).Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
		"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Offset(offset).Limit(size).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	for _, value1 := range *list {
		if value1.Departid == departid {
			machinelist = append(machinelist, value1)
		}
	}
	return count, machinelist, err
}

func IsUUIDExist(uuid string) (bool, error) {
	var Machine MachineNode
	err := mysqlmanager.MySQL().Where("machine_uuid=?", uuid).Find(&Machine).Error
	return Machine.DepartId != 0, err
}

// 根据uuid获取部门id
func UUIDForDepartId(uuid string) (int, error) {
	var Machine MachineNode
	err := mysqlmanager.MySQL().Where("machine_uuid=?", uuid).Find(&Machine).Error
	return Machine.DepartId, err
}

// 根据机器id获取机器的uuid
func MachineIdToUUID(id int) (string, error) {
	var Machine MachineNode
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&Machine).Error
	return Machine.MachineUUID, err
}

// 更新机器状态
func UpdateMachineState(uuid string, state int) error {
	return mysqlmanager.MySQL().Model(&MachineNode{}).Where("machine_uuid=?", uuid).UpdateColumn("state", state).Error
}

// 更新机器IP及状态
func UpdateMachineIPState(uuid, ip string, state int) error {
	Ma := &MachineNode{
		State: state,
		IP:    ip,
	}
	return mysqlmanager.MySQL().Model(&MachineNode{}).Where("machine_uuid=?", uuid).Updates(Ma).Error
}

// 新增agent机器
func AddNewMachine(Machine MachineNode) error {
	return mysqlmanager.MySQL().Save(&Machine).Error
}

// 分页获取该部门下的所有机器
func GetMachinePaged(departId []int, offset, size int) (int64, []Res, error) {
	var count int64 = 0
	var machinelist []Res
	var err error
	for _, value := range departId {
		num, data, errtemp := ReturnMachinePaged(value, offset, size)
		count = count + num
		machinelist = append(machinelist, data...)
		err = errtemp
	}
	return count, machinelist, err
}

// 获取该部门下的所有机器
func MachineList(departId []int) (machinelist []Res, err error) {
	for _, value := range departId {
		list := &[]Res{}
		err = mysqlmanager.MySQL().Table("machine_node").Where("depart_id=?", value).Select("machine_node.id as id,machine_node.depart_id as departid," +
			"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
			"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&list).Error
		if err != nil {
			return
		}
		for _, value1 := range *list {
			if value1.Departid == value {
				machinelist = append(machinelist, value1)
			}
		}
	}
	return
}

func MachineStore(departid int) ([]MachineNode, error) {
	var Machineinfo []MachineNode
	err := mysqlmanager.MySQL().Where("depart_id=?", departid).Find(&Machineinfo).Error
	return Machineinfo, err
}

func MachineInfo(id int) (*MachineNode, error) {
	var machine = &MachineNode{}
	if err := mysqlmanager.MySQL().Where("id=?", id).Find(machine).Error; err != nil {
		return nil, err
	}
	return machine, nil
}

func UpdateMachineDepartState(uuid string, DeptId int) error {
	Ma := &MachineNode{
		DepartId: DeptId,
	}
	return mysqlmanager.MySQL().Model(&MachineNode{}).Where("machineuuid=?", uuid).Updates(Ma).Error
}

// 根据机器id获取机器信息
func MachineData(MacId int) (MachineNode, error) {
	var m MachineNode
	err := mysqlmanager.MySQL().Where("id=?", MacId).Find(&m).Error
	return m, err
}

// 获取所有机器
func AllMachine() ([]MachineNode, error) {
	var m []MachineNode
	err := mysqlmanager.MySQL().Find(&m).Error
	return m, err
}

func MachineAllData() ([]Res, error) {
	var mch []Res
	err := mysqlmanager.MySQL().Table("machine_node").Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
		"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&mch).Error
	return mch, err
}

// 获取某一级部门下的所有机器
func SomeDepartMachine(Departids []int) (lists []MachineNode, err error) {
	for _, id := range Departids {
		list := []MachineNode{}
		err = mysqlmanager.MySQL().Where("depart_id = ?", id).Find(&list).Error
		if err != nil {
			return
		}
		lists = append(lists, list...)
	}
	return
}

// 根据uuid获取机器的ip、状态和部门
func MachineBasic(uuid string) (ip string, state int, dept string, err error) {
	var machine MachineNode
	var depart DepartNode
	err = mysqlmanager.MySQL().Where("machine_uuid = ?", uuid).Find(&machine).Error
	if err != nil {
		return machine.IP, machine.State, "", err
	}
	err = mysqlmanager.MySQL().Where("id = ?", machine.DepartId).Find(&depart).Error
	return machine.IP, machine.State, depart.Depart, err
}

// 使用uuid删除机器
func DeleteMachine(machinedeluuid string) (err error) {
	var machine MachineNode
	UUIDExistbool, err := IsUUIDExist(machinedeluuid)
	if err != nil {
		return err
	}
	if UUIDExistbool {
		if err := mysqlmanager.MySQL().Where("machine_uuid=?", machinedeluuid).Unscoped().Delete(machine).Error; err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("该机器不存在")
}
