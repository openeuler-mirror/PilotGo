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
)

type DepartNode struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID          int    `gorm:"type:int(100);not null" json:"pid"`
	ParentDepart string `gorm:"type:varchar(100);not null" json:"parentdepart"`
	Depart       string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate   int    `gorm:"type:int(100);not null" json:"nodelocate"`
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

// 查询子部门
func Pid2Depart(pid int) ([]DepartNode, error) {
	var DepartInfo []DepartNode
	err := mysqlmanager.MySQL().Where("p_id=?", pid).Find(&DepartInfo).Error
	return DepartInfo, err
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

// 未定
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

// 根据部门名字查询id和pid,批量操作需完善
func GetPidAndId(depart string) (pid, id int, err error) {
	var dep DepartNode
	err = mysqlmanager.MySQL().Where("depart=?", depart).Find(&dep).Error
	return dep.PID, dep.ID, err
}

// 根据部门ids查询所属部门，待确定
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

// 创建公司组织
func CreateOrganization() error {
	var Depart DepartNode
	err := mysqlmanager.MySQL().Where("p_id=?", global.Departroot).Find(&Depart).Error
	if err != nil {
		return err
	}
	if Depart.ID == 0 {
		Depart = DepartNode{
			PID: global.Departroot,
			//ParentDepart: "",
			Depart:     "组织名",
			NodeLocate: global.Departroot,
		}
		return mysqlmanager.MySQL().Save(&Depart).Error
	}
	return nil
}
