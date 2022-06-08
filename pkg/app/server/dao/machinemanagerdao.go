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
 * LastEditTime: 2022-04-07 13:28:41
 * Description: 部门管理数据库相关函数
 ******************************************************************************/
package dao

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

func IsUUIDExist(uuid string) bool {
	var Machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&Machine)
	return Machine.DepartId != 0
}

// 根据uuid获取部门id
func UUIDForDepartId(uuid string) int {
	var Machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&Machine)
	return Machine.DepartId
}

// agent机器断开
func MachineStatusToOffline(uuid string) {
	var Machine model.MachineNode
	Ma := model.MachineNode{
		State: model.OffLine,
	}
	mysqlmanager.DB.Model(&Machine).Where("machine_uuid=?", uuid).Update(&Ma)
}

// agent机器未分配
func MachineStatusToFree(uuid, ip string) {
	var Machine model.MachineNode
	Ma := model.MachineNode{
		State: model.Free,
		IP:    ip,
	}
	mysqlmanager.DB.Model(&Machine).Where("machine_uuid=?", uuid).Update(&Ma)
}

// agent机器连接正常
func MachineStatusToNormal(uuid, ip string) {
	var Machine model.MachineNode
	Ma := model.MachineNode{
		State: model.Normal,
		IP:    ip,
	}
	mysqlmanager.DB.Model(&Machine).Where("machine_uuid=?", uuid).Update(&Ma)
}

// 新增agent机器
func AddNewMachine(Machine model.MachineNode) {
	mysqlmanager.DB.Save(&Machine)
}

// 获取该部门下的所有机器
func MachineList(departId []int) (machinelist []model.Res) {
	for _, value := range departId {
		list := &[]model.Res{}
		mysqlmanager.DB.Table("machine_node").Where("depart_id=?", value).Select("machine_node.id as id,machine_node.depart_id as departid," +
			"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
			"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&list)
		for _, value1 := range *list {
			if value1.Departid == value {
				machinelist = append(machinelist, value1)
			}
		}
	}
	return
}

func MachineStore(departid int) []model.MachineNode {
	var Machineinfo []model.MachineNode
	mysqlmanager.DB.Where("depart_id=?", departid).Find(&Machineinfo)
	return Machineinfo
}

func ModifyMachineDepart(MadId int, DeptId int) {
	var Machine model.MachineNode
	Ma := model.MachineNode{
		DepartId: DeptId,
		State:    model.Normal,
	}
	mysqlmanager.DB.Model(&Machine).Where("id=?", MadId).Update(&Ma)
}
func ModifyMachineDepart2(MadId int, DeptId int) {
	var Machine model.MachineNode
	Ma := model.MachineNode{
		DepartId: DeptId,
		State:    model.Free,
	}
	mysqlmanager.DB.Model(&Machine).Where("id=?", MadId).Update(&Ma)
}

// 根据机器id获取机器信息
func MachineData(MacId int) model.MachineNode {
	var m model.MachineNode
	mysqlmanager.DB.Where("id=?", MacId).Find(&m)
	return m
}

// 获取所有机器
func AllMachine() []model.MachineNode {
	var m []model.MachineNode
	mysqlmanager.DB.Find(&m)

	return m
}

func MachineAllData() []model.Res {
	var mch []model.Res
	mysqlmanager.DB.Table("machine_node").Select("machine_node.id as id,machine_node.depart_id as departid," +
		"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
		"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&mch)
	return mch
}

// 获取某一级部门下的所有机器
func SomeDepartMachine(Departids []int) (lists []model.MachineNode) {
	for _, id := range Departids {
		list := []model.MachineNode{}
		mysqlmanager.DB.Where("depart_id = ?", id).Find(&list)
		lists = append(lists, list...)
	}
	return
}

// 根据uuid获取机器的ip、状态和部门
func MachineBasic(uuid string) (ip string, state int, dept string) {
	var machine model.MachineNode
	var depart model.DepartNode
	mysqlmanager.DB.Where("machine_uuid = ?", uuid).Find(&machine)
	mysqlmanager.DB.Where("id = ?", machine.DepartId).Find(&depart)
	return machine.IP, machine.State, depart.Depart
}

// 根据uuid获取机器的ip
func UUID2MacIP(uuid string) (ip string) {
	var machine model.MachineNode
	mysqlmanager.DB.Where("machine_uuid = ?", uuid).Find(&machine)
	return machine.IP
}
