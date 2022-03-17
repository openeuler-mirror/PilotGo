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
 * Date: 2021-12-18 02:33:55
 * LastEditTime: 2022-03-17 15:06:37
 * Description: 用户数据存储结构体
 ******************************************************************************/
package model

import (
	"time"

	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

type User struct {
	ID           uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time
	DepartFirst  int    `gorm:"size:25" json:"departPId,omitempty"`
	DepartSecond int    `gorm:"size:25" json:"departid,omitempty"`
	DepartName   string `gorm:"size:25" json:"departName,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone        string `gorm:"size:11" json:"phone,omitempty"`
	Email        string `gorm:"type:varchar(30);not null" json:"email,omitempty"`
	UserType     int    `json:"userType,omitempty"`
	RoleID       string `json:"roleId,omitempty"`
	Enable       string `gorm:"size:10;not null" json:"enable,omitempty"`
}
type ReturnUser struct {
	ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	DepartFirst  int    `json:"departPId"`
	DepartSecond int    `json:"departid"`
	DepartName   string `json:"departName"`
	Username     string `json:"username"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	UserType     int    `json:"userType"`
	RoleID       string `json:"roleId"`
	Role         string `json:"role"`
}

func (u *User) All(q *PaginationQ) (list *[]ReturnUser, total uint, err error) {
	list = &[]ReturnUser{}
	// tx := mysqlmanager.DB.Order("ID desc").Find(list)
	tx := mysqlmanager.DB.Order("ID desc").Table("user").Select("user.id,user.created_at,user.depart_first," +
		"user.depart_second,user.depart_name,user.username,user.phone,user.email," +
		"user.user_type,user.role_id,user_role.role").Joins("left join user_role on user.role_id = user_role.id").Scan(list)
	total, err = CrudAll(q, tx, list)
	return
}
