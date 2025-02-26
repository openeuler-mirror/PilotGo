/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type MachineNode struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	DepartId    int    `gorm:"type:int(100);not null" json:"departid"`
	IP          string `gorm:"type:varchar(100)" json:"ip"`
	MachineUUID string `gorm:"type:varchar(100);not null" json:"machineuuid"`
	CPU         string `gorm:"type:varchar(100)" json:"CPU"`
	RunStatus   string `gorm:"type:varchar(100)" json:"runstatus"`
	MaintStatus string `gorm:"type:varchar(100)" json:"maintstatus"`
	Systeminfo  string `gorm:"type:varchar(100)" json:"sysinfo"`
}

type Res struct {
	ID          int    `json:"id"`
	Departid    int    `json:"departid"`
	Departname  string `json:"departname"`
	IP          string `json:"ip"`
	UUID        string `json:"uuid"`
	CPU         string `json:"cpu"`
	Runstatus   string `json:"runstatus"`
	Maintstatus string `json:"maintstatus"`
	Systeminfo  string `json:"systeminfo"`
}

// 新增agent机器
func (Machine *MachineNode) Add() error {
	return mysqlmanager.MySQL().Save(Machine).Error
}

// 使用uuid删除机器
func DeleteMachine(machinedeluuid string) error {
	var machine MachineNode
	return mysqlmanager.MySQL().Where("machine_uuid=?", machinedeluuid).Unscoped().Delete(machine).Error
}

// 变更机器信息
func UpdateMachine(uuid string, ma *MachineNode) error {
	return mysqlmanager.MySQL().Model(&MachineNode{}).Where("machine_uuid=?", uuid).Updates(ma).Error
}

func ReturnMachinePagedByDepartid(departids []int, offset, size int) (int64, []Res, error) {
	var count int64
	var list []Res
	err := mysqlmanager.MySQL().Model(MachineNode{}).Where("depart_id IN ?", departids).Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, machine_node.cpu as cpu,machine_node.run_status as runstatus," +
		"machine_node.maint_status as maintstatus,machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Offset(offset).Limit(size).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	return count, list, err
}

func ReturnMachinePagedBySearch(search string, offset, size int) (int64, []Res, error) {
	var count int64
	var list []Res

	err := mysqlmanager.MySQL().Model(MachineNode{}).Select("machine_node.id as id,machine_node.depart_id as departid,"+
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, machine_node.cpu as cpu,machine_node.run_status as runstatus,"+
		"machine_node.maint_status as maintstatus,machine_node.systeminfo as systeminfo").Where("ip REGEXP ?", search).Or("machine_uuid REGEXP ?", search).Or("depart_node.depart REGEXP ?", search).Or("cpu REGEXP ?", search).Or("run_status REGEXP ?", search).Or("maint_status REGEXP ?", search).Or("systeminfo REGEXP ?", search).Joins("left join depart_node on machine_node.depart_id = depart_node.id").Offset(offset).Limit(size).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	return count, list, err
}

func ReturnSpecifiedDepart(id int, res *[]int) {
	temp, err := SubDepartId(id)
	if err != nil {
		logger.Error("%v", err.Error())
	}
	if len(temp) == 0 {
		return
	}
	for _, value := range temp {
		*res = append(*res, value)
		ReturnSpecifiedDepart(value, res)
	}
}

// 分页获取该部门下的所有机器
func GetMachinePaged(departId []int, offset, size int, search string) (int64, []Res, error) {
	var count int64 = 0
	var machinelist []Res
	var err error

	if search != "" {
		var temp_depart_id int
		var TheDeptAndSubDeptIds []int
		var isDepart bool

		depart_nodes, err := GetAllDepart()
		if err != nil {
			return 0, []Res{}, err
		}
		for _, dnode := range depart_nodes {
			if strings.Contains(dnode.Depart, search) {
				temp_depart_id = dnode.ID
				isDepart = true
				break
			}
		}

		if isDepart {
			ReturnSpecifiedDepart(temp_depart_id, &TheDeptAndSubDeptIds)
			TheDeptAndSubDeptIds = append([]int{temp_depart_id}, TheDeptAndSubDeptIds...)
			departId = TheDeptAndSubDeptIds
		} else {
			count, machinelist, err = ReturnMachinePagedBySearch(search, offset, size)
			return count, machinelist, err
		}
	}

	count, machinelist, err = ReturnMachinePagedByDepartid(departId, offset, size)
	return count, machinelist, err
}

// 根据机器id获取机器的uuid
func MachineIdToUUID(id int) (string, error) {
	var Machine MachineNode
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&Machine).Error
	return Machine.MachineUUID, err
}

func MachineInfoByUUID(uuid string) (MachineNode, error) {
	var machine MachineNode
	err := mysqlmanager.MySQL().Where("machine_uuid=?", uuid).Find(&machine).Error
	return machine, err
}

// 根据机器id获取机器信息
func MachineInfo(MacId int) (MachineNode, error) {
	var m MachineNode
	err := mysqlmanager.MySQL().Where("id=?", MacId).Find(&m).Error
	return m, err
}

// 获取所有机器
func MachineAll() ([]Res, error) {
	var mch []Res
	err := mysqlmanager.MySQL().Table("machine_node").Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, machine_node.cpu as cpu,machine_node.run_status as runstatus," +
		"machine_node.maint_status as maintstatus, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&mch).Error
	return mch, err
}

// 获取该部门下的所有机器
func MachineList(departId []int) ([]Res, error) {
	machinelist := []Res{}
	for _, value := range departId {
		list := []Res{}
		err := mysqlmanager.MySQL().Table("machine_node").Where("depart_id=?", value).Select("machine_node.id as id,machine_node.depart_id as departid," +
			"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid,  machine_node.cpu as cpu,machine_node.run_status as runstatus," +
			"machine_node.maint_status as maintstatus, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&list).Error
		if err != nil {
			return []Res{}, err
		}
		machinelist = append(machinelist, list...)
	}
	return machinelist, nil
}

// 统计总数
func CountMachineNode(departid interface{}) (int, error) {
	var count int64
	var err error
	if departid != nil {
		err = mysqlmanager.MySQL().Model(MachineNode{}).Where("depart_id= ?", departid).Count(&count).Error
		return int(count), err
	}
	err = mysqlmanager.MySQL().Model(MachineNode{}).Distinct("id").Count(&count).Error
	return int(count), err
}

// 统计runstatus
func CountRunStatus(status string, departid interface{}) (int, error) {
	var count int64
	var err error
	if departid != nil {
		err = mysqlmanager.MySQL().Model(&MachineNode{}).Where("run_status=? and depart_id= ?", status, departid).Count(&count).Error
		return int(count), err
	}
	err = mysqlmanager.MySQL().Model(MachineNode{}).Where("run_status=? ", status).Count(&count).Error
	return int(count), err
}

// 统计maintstatus
func CountMaintStatus(status string, departid interface{}) (int, error) {
	var count int64
	var err error
	if departid != nil {
		err = mysqlmanager.MySQL().Model(MachineNode{}).Where("maint_status=? and depart_id= ?", status, departid).Count(&count).Error
		return int(count), err
	}
	err = mysqlmanager.MySQL().Model(MachineNode{}).Where("maint_status=? ", status).Count(&count).Error
	return int(count), err
}
