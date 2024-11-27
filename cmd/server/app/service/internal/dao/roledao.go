/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/pkg/utils"
)

type Role struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name        string `gorm:"not null;unique" json:"name"` // 超管和部门等级
	Description string `json:"description"`
	Menus       string `json:"menus"`
	ButtonID    string `json:"buttonId"`
}

type RoleButton struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Button string `json:"button"`
}

// 根据角色名称返回角色id和用户类型
func GetRoleId(name string) (roleId int, err error) {
	var role Role
	err = mysqlmanager.MySQL().Where("name = ?", name).Find(&role).Error
	if err != nil {
		return 0, err
	}
	return role.ID, nil
}

// 根据id获取该角色的所有信息
func GetRoleById(roleid int) (Role, error) {
	var role Role
	err := mysqlmanager.MySQL().Where("id=?", roleid).Find(&role).Error
	return role, err
}

// 分页查询
func GetRolePaged(offset, size int) (int64, []Role, error) {
	var count int64
	var userRoles []Role
	err := mysqlmanager.MySQL().Model(Role{}).Order("id desc").Offset(offset).Limit(size).Find(&userRoles).Offset(-1).Limit(-1).Count(&count).Error
	return count, userRoles, err
}

// 获取所有的用户角色
func GetRoles() ([]Role, error) {
	var roles []Role
	err := mysqlmanager.MySQL().Order("id").Find(&roles).Error
	return roles, err
}

// 新增角色
func AddRole(r *Role) error {
	return mysqlmanager.MySQL().Save(r).Error
}

// 删除用户角色
func DeleteRole(roleId int) error {
	var UserRole Role
	return mysqlmanager.MySQL().Where("id = ?", roleId).Unscoped().Delete(UserRole).Error
}

// 修改角色描述
func UpdateRoleDescription(name, desc string) error {
	var UserRole Role
	return mysqlmanager.MySQL().Model(&UserRole).Where("name = ?", name).Update("description", desc).Error
}

const SuperUser = "admin"
const SuperUserPasswd = "admin"

// 创建管理员账户
func CreateAdministratorUser() error {
	var role Role
	mysqlmanager.MySQL().Where("name =?", "admin").Find(&role)
	if role.ID == 0 {
		role = Role{
			Name:        "admin",
			Description: "超级管理员",
		}
		mysqlmanager.MySQL().Create(&role)
	}

	var user User
	mysqlmanager.MySQL().Where("username =?", "admin").Find(&user)
	if user.ID == 0 {
		bs, err := utils.CryptoPassword(SuperUserPasswd)
		if err != nil {
			return err
		}

		user := User{
			CreatedAt: time.Time{},
			DepartId:  global.UncateloguedDepartId,
			Username:  strings.Split(SuperUser, "@")[0],
			Password:  string(bs),
			Email:     SuperUser,
		}
		mysqlmanager.MySQL().Create(&user)

		ur := UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
		mysqlmanager.MySQL().Create(&ur)
	}
	return nil
}

func GetNamesByRoleIds(roleids []int) ([]string, error) {
	var RoleNames []string
	for _, v := range roleids {
		role, err := GetRoleById(v)
		if err != nil {
			return nil, err
		}
		RoleNames = append(RoleNames, role.Name)
	}
	return RoleNames, nil
}
