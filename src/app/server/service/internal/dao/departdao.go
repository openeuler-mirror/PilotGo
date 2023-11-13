/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-04-28 13:08:08
 * LastEditTime: 2023-06-28 16:02:24
 * Description: depart相关数据获取
 ******************************************************************************/
package dao

import (
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gorm.io/gorm"
)

type DepartNode struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID          int    `gorm:"type:int(100);not null" json:"pid"`
	ParentDepart string `gorm:"type:varchar(100);not null" json:"parentdepart"`
	Depart       string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate   int    `gorm:"type:int(100);not null" json:"nodelocate"`
	//根节点为0,普通节点为1
}

type DepartTreeNode struct {
	Label    string            `json:"label"`
	Id       int               `json:"id"`
	Pid      int               `json:"pid"`
	Children []*DepartTreeNode `json:"children"`
}

type NewDepart struct {
	DepartID   int    `json:"DepartID" binding:"required" msg:"部门id不能为空"`
	DepartName string `json:"DepartName" binding:"required" msg:"部门名称不能为空"`
}

type DeleteDepart struct {
	DepartID int `json:"DepartID"`
}
type AddDepart struct {
	ParentID     int    `json:"PID"`
	ParentDepart string `json:"ParentDepart"`
	DepartName   string `json:"Depart" binding:"required" msg:"部门名称不能为空"`
}

func IsParentDepartExist(parent string) (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("depart=? ", parent).Find(&Depart).Error
	return Depart.ID != 0, err
}
func IsDepartNodeExist(parent string, depart string) (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("depart=? and parent_depart=?", depart, parent).Find(&Depart).Error
	// mysqlmanager.MySQL().Where("", parent).Find(&Depart)
	return Depart.ID != 0, err
}
func IsDepartIDExist(ID int) (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("id=?", ID).Find(&Depart).Error
	return Depart.ID != 0, err
}

func IsRootExist() (bool, error) {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("node_locate=?", 0).Find(&Depart).Error
	return Depart.ID != 0, err
}

func DepartStore() ([]DepartNode, error) {
	var Depart []DepartNode
	err := mysqlmanager.MySQL().Find(&Depart).Error
	return Depart, err
}

func UpdateDepart(DepartID int, DepartName string) error {
	var DepartInfo DepartNode
	Depart := DepartNode{
		Depart: DepartName,
	}
	return mysqlmanager.MySQL().Model(&DepartInfo).Where("id=?", DepartID).Updates(&Depart).Error
}

func UpdateParentDepart(DepartID int, DepartName string) error {
	var DepartInfo DepartNode
	Depart := DepartNode{
		ParentDepart: DepartName,
	}
	return mysqlmanager.MySQL().Model(&DepartInfo).Where("p_id=?", DepartID).Updates(&Depart).Error
}

func Pid2Depart(pid int) ([]DepartNode, error) {
	var DepartInfo []DepartNode
	err := mysqlmanager.MySQL().Where("p_id=?", pid).Find(&DepartInfo).Error
	return DepartInfo, err
}

func Deletedepartdata(needdelete []int) error {
	var DepartInfo []DepartNode
	return mysqlmanager.MySQL().Where("id=?", needdelete[0]).Delete(&DepartInfo).Error
}

// 向需要删除的depart的组内增加需要删除的子节点
func Insertdepartlist(needdelete []int, str string) ([]int, error) {
	var DepartInfo []DepartNode

	err := mysqlmanager.MySQL().Where("p_id=?", str).Find(&DepartInfo).Error
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}
	return needdelete, err
}

// 根据部门名字查询id和pid
func GetPidAndId(depart string) (pid, id int, err error) {
	var dep DepartNode
	err = mysqlmanager.MySQL().Where("depart=?", depart).Find(&dep).Error
	return dep.PID, dep.ID, err
}

// 添加部门
func AddDepartMessage(db *gorm.DB, depart *DepartNode) error {
	return db.Create(depart).Error
}

// 根据部门id获取部门名称
func DepartIdToGetDepartName(id int) (string, error) {
	var departNames DepartNode
	err := mysqlmanager.MySQL().Where("id =?", id).Find(&departNames).Error
	return departNames.Depart, err
}

// 根据部门ids查询所属部门
func DepartIdsToGetDepartNames(ids []int) (names []string) {
	for _, id := range ids {
		var depart DepartNode
		err := mysqlmanager.MySQL().Where("id = ?", id).Find(&depart).Error
		if err != nil {
			logger.Error(err.Error())
		}
		names = append(names, depart.Depart)
	}
	return
}

// 获取下级部门id
func SubDepartId(id int) ([]int, error) {
	var depart []DepartNode
	err := mysqlmanager.MySQL().Where("p_id=?", id).Find(&depart).Error

	res := make([]int, 0)
	for _, value := range depart {
		res = append(res, value.ID)
	}
	return res, err
}

// 获取所有的一级部门id
func FirstDepartId() (departIds []int, err error) {
	departs := []DepartNode{}
	err = mysqlmanager.MySQL().Where("p_id = ?", global.UncateloguedDepartId).Find(&departs).Error
	for _, depart := range departs {
		departIds = append(departIds, depart.ID)
	}
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
			PID:          global.Departroot,
			ParentDepart: "",
			Depart:       "组织名",
			NodeLocate:   global.Departroot,
		}
		return mysqlmanager.MySQL().Save(&Depart).Error
	}
	return nil
}
