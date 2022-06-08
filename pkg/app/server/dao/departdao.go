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
 * LastEditTime: 2022-00-08 14:25:41
 * Description: depart相关数据获取
 ******************************************************************************/
package dao

import (
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

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

func IsRootExist() bool {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("node_locate=?", 0).Find(&Depart)
	return Depart.ID != 0
}

func DepartStore() []model.DepartNode {
	var Depart []model.DepartNode
	mysqlmanager.DB.Find(&Depart)
	return Depart
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

func Pid2Depart(pid int) []model.DepartNode {
	var DepartInfo []model.DepartNode
	mysqlmanager.DB.Where("p_id=?", pid).Find(&DepartInfo)
	return DepartInfo
}

func Deletedepartdata(needdelete []int) {
	var DepartInfo []model.DepartNode
	mysqlmanager.DB.Where("id=?", needdelete[0]).Delete(&DepartInfo)
}

//向需要删除的depart的组内增加需要删除的子节点
func Insertdepartlist(needdelete []int, str string) []int {
	var DepartInfo []model.DepartNode

	mysqlmanager.DB.Where("p_id=?", str).Find(&DepartInfo)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}
	return needdelete
}

// 根据部门名字查询id和pid
func GetPidAndId(depart string) (pid, id int) {
	var dep model.DepartNode
	mysqlmanager.DB.Where("depart=?", depart).Find(&dep)
	return dep.PID, dep.ID
}

//添加部门
func AddDepart(db *gorm.DB, depart *model.DepartNode) error {
	err := db.Create(depart).Error
	return err
}

// 根据部门id获取部门名称
func DepartIdToGetDepartName(id int) string {
	var departNames model.DepartNode
	mysqlmanager.DB.Where("id =?", id).Find(&departNames)
	return departNames.Depart
}

// 根据部门ids查询所属部门
func DepartIdsToGetDepartNames(ids []int) (names []string) {
	for _, id := range ids {
		var depart model.DepartNode
		mysqlmanager.DB.Where("id = ?", id).Find(&depart)
		names = append(names, depart.Depart)
	}
	return
}

// 获取下级部门id
func SubDepartId(id int) []int {
	var depart []model.DepartNode
	mysqlmanager.DB.Where("p_id=?", id).Find(&depart)

	res := make([]int, 0)
	for _, value := range depart {
		res = append(res, value.ID)
	}
	return res
}

// 获取所有的一级部门id
func FirstDepartId() (departIds []int) {
	departs := []model.DepartNode{}
	mysqlmanager.DB.Where("p_id = ?", model.UncateloguedDepartId).Find(&departs)
	for _, depart := range departs {
		departIds = append(departIds, depart.ID)
	}
	return
}

// 创建公司组织
func CreateOrganization() {
	var Depart model.DepartNode
	mysqlmanager.DB.Where("p_id=?", model.Departroot).Find(&Depart)
	if Depart.ID == 0 {
		Depart = model.DepartNode{
			PID:          model.Departroot,
			ParentDepart: "",
			Depart:       "组织名",
			NodeLocate:   model.Departroot,
		}
		mysqlmanager.DB.Save(&Depart)
	}
}
