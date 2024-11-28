/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/pkg/global"
)

type DepartNode struct {
	ID         int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID        int    `gorm:"type:int(100);not null" json:"pid"`
	Depart     string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate int    `gorm:"type:int(100);not null" json:"nodelocate"`
	//根节点为0,普通节点为1
}

// 添加部门
func (deport *DepartNode) Add() error {
	return mysqlmanager.MySQL().Save(deport).Error
}

// 查询部门是否存在
func IsDepartNodeExist(pid int, depart string) (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("depart=? and p_id=?", depart, pid).Find(&Depart).Error
	return Depart.ID != 0, err
}

func IsRootExist() (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("node_locate=?", 0).Find(&Depart).Error
	return Depart.ID != 0, err
}

// 根据部门id获取部门名结构体
func GetDepartById(id int) (DepartNode, error) {
	var departName DepartNode
	err := mysqlmanager.MySQL().Where("id =?", id).Find(&departName).Error
	return departName, err
}

// 获取下级部门id
func SubDepartId(pid int) ([]int, error) {
	var departids []int
	err := mysqlmanager.MySQL().Model(&DepartNode{}).Select("id").Where("p_id=?", pid).Find(&departids).Error
	return departids, err
}

// 查询所有部门
func GetAllDepart() ([]DepartNode, error) {
	var Depart []DepartNode
	err := mysqlmanager.MySQL().Find(&Depart).Error
	return Depart, err
}

// 修改部门名字
func UpdateDepart(DepartID int, DepartName string) error {
	return mysqlmanager.MySQL().Model(&DepartNode{}).Where("id=?", DepartID).Update("depart", DepartName).Error
}

// 删除节点
func Deletedepartdata(needdelete []int) error {
	var DepartInfo []DepartNode
	return mysqlmanager.MySQL().Where("id in (?)", needdelete).Delete(&DepartInfo).Error
}

// 根据部门名字查询id和pid,批量操作需完善
func GetPidAndId(depart string) (pid, id int, err error) {
	var dep DepartNode
	err = mysqlmanager.MySQL().Where("depart=?", depart).Find(&dep).Error
	return dep.PID, dep.ID, err
}

// 根据部门ids查询所属部门，待确定
func DepartIdsToGetDepartNames(ids []int) (names []string) {
	mysqlmanager.MySQL().Model(&DepartNode{}).Select("depart").Where("id=?", ids[0]).Find(&names)
	return
}

// 创建公司组织
func CreateOrganization() error {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("p_id=?", global.Departroot).Find(&Depart).Error
	if err != nil {
		return err
	}
	if Depart.ID == 0 {
		Depart = DepartNode{
			PID:        global.Departroot,
			Depart:     "组织名",
			NodeLocate: global.Departroot,
		}
		return mysqlmanager.MySQL().Save(&Depart).Error
	}
	return nil
}
