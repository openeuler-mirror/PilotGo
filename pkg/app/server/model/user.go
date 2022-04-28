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
 * LastEditTime: 2022-04-27 09:31:56
 * Description: 用户数据存储结构体
 ******************************************************************************/
package model

import (
	"time"
)

type User struct {
	ID           uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time
	DepartFirst  int    `gorm:"size:25" json:"departPid,omitempty"`
	DepartSecond int    `gorm:"size:25" json:"departId,omitempty"`
	DepartName   string `gorm:"size:25" json:"departName,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone        string `gorm:"size:11" json:"phone,omitempty"`
	Email        string `gorm:"type:varchar(30);not null" json:"email,omitempty"`
	UserType     int    `json:"userType,omitempty"`
	RoleID       string `json:"role,omitempty"`
}
type ReturnUser struct {
	ID           uint     `json:"id"`
	DepartFirst  int      `json:"departPId"`
	DepartSecond int      `json:"departid"`
	DepartName   string   `json:"departName"`
	Username     string   `json:"username"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	UserType     int      `json:"userType"`
	Roles        []string `json:"role"`
}
type UserDto struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Name:     user.Username,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	}
}

type Userdel struct {
	Emails []string `json:"email"`
}

const (
	// 超级管理员
	AdminUserType = 0
	// 部门管理员
	DepartManagerType = 1
	// 普通用户
	OrdinaryUserType = 2
	// 其他用户，如实习生
	OtherUserType = 3
	//普通用户角色id
	OrdinaryUserRoleId = 3
	// 默认用户密码
	DefaultUserPassword = "123456"
)
