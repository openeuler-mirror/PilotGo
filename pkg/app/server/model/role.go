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
 * Date: 2022-03-07 15:56:45
 * LastEditTime: 2022-03-16 14:05:12
 * Description: 用户权限管理
 ******************************************************************************/
package model

import "openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"

type UserRole struct {
	ID    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Role  string `json:"role"` // 超管和部门等级
	Type  int    `json:"type"`
	Menus string `json:"menus"`
}

type RoleButton struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Menu   string `json:"menu"`
	Button string `json:"button"`
}

type AddRole struct {
	Email  string `json:"email"`
	RoleID int    `json:"roleid"`
}

func (u *UserRole) All(q *PaginationQ) (list *[]UserRole, total uint, err error) {
	list = &[]UserRole{}
	tx := mysqlmanager.DB.Find(list)
	total, err = CrudAll(q, tx, list)
	return
}
