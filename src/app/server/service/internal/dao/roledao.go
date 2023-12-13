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
 * Date: 2021-01-24 15:08:08
 * LastEditTime: 2023-09-04 14:03:47
 * Description: 角色模块相关数据获取
 ******************************************************************************/
package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/global"
	"gitee.com/openeuler/PilotGo/utils"
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

// 登录用户的权限按钮
func PermissionButtons(button string) ([]string, error) {
	var buttons []string
	if button == "" {
		return []string{}, nil
	}
	IDs := strings.Split(button, ",")

	for _, id := range IDs {
		var SubButton RoleButton
		i, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		err = mysqlmanager.MySQL().Where("id = ?", i).Find(&SubButton).Error
		if err != nil {
			return buttons, err
		}
		button := SubButton.Button
		buttons = append(buttons, button)
	}
	return buttons, nil
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
	if r.Name == "" {
		return fmt.Errorf("角色名不能为空")
	}

	return mysqlmanager.MySQL().Save(r).Error
}

// 删除用户角色
func DeleteRole(roleId int) error {
	var UserRole Role
	return mysqlmanager.MySQL().Where("id = ?", roleId).Unscoped().Delete(UserRole).Error
}

// 修改角色名称
func UpdateRoleName(roleId int, name string) error {
	var UserRole Role
	return mysqlmanager.MySQL().Model(&UserRole).Where("id = ?", roleId).Update("role", name).Error
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
	mysqlmanager.MySQL().Where("role =?", "admin").Find(&role)
	if role.ID == 0 {
		role = Role{
			Name:        "admin",
			Description: "超级管理员",
		}
		mysqlmanager.MySQL().Create(&role)
		bs, err := utils.CryptoPassword(SuperUserPasswd)
		if err != nil {
			return err
		}

		user := User{
			CreatedAt:    time.Time{},
			DepartFirst:  global.Departroot,
			DepartSecond: global.UncateloguedDepartId,
			DepartName:   "超级用户",
			Username:     strings.Split(SuperUser, "@")[0],
			Password:     string(bs),
			Email:        SuperUser,
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
