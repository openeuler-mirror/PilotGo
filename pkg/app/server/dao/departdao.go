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
 * LastEditTime: 2022-04-28 14:25:41
 * Description: depart相关数据获取
 ******************************************************************************/
package dao

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

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
