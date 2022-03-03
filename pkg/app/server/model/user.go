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
 * LastEditTime: 2022-03-03 14:15:34
 * Description: 用户数据存储结构体
 ******************************************************************************/
package model

import (
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(25);not null" json:"username,omitempty" form:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password,omitempty" form:"password"`
	Phone    string `gorm:"size:255" json:"phone,omitempty" form:"phone"`
	Email    string `gorm:"type:varchar(30);not null" json:"email,omitempty" form:"email"`
	Enable   string `gorm:"size:10;not null" json:"enable,omitempty"`
}

func (u *User) All(q *PaginationQ) (list *[]User, total uint, err error) {
	list = &[]User{}
	tx := mysqlmanager.DB.Order("ID desc").Find(list)
	total, err = CrudAll(q, tx, list)
	return
}
