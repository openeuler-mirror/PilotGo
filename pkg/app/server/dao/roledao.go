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
 * LastEditTime: 2022-04-27 17:25:41
 * Description: 角色模块相关数据获取
 ******************************************************************************/
package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"openeuler.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type UserRole struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT"`
	Role        string `json:"role"` // 超管和部门等级
	Type        int    `json:"type"`
	Description string `json:"description"`
	Menus       string `json:"menus"`
	ButtonID    string `json:"buttonId"`
}

type RoleButton struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Button string `json:"button"`
}

// 根据角色名称返回角色id和用户类型
func GetRoleIdAndUserType(role string) (roleId string, user_type int, err error) {
	var Role UserRole
	err = mysqlmanager.MySQL().Where("role = ?", role).Find(&Role).Error
	if err != nil {
		return "", 0, err
	}
	roleID := strconv.Itoa(Role.ID)
	var userType int
	if Role.ID > global.OrdinaryUserRoleId {
		userType = global.OtherUserType
	} else {
		userType = Role.ID - 1
	}
	return roleID, userType, nil
}

// 根据id获取该角色的所有信息
func RoleIdToGetAllInfo(roleid int) (UserRole, error) {
	var role UserRole
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

// 获取所有的用户角色
func GetRoleList() ([]UserRole, error) {
	var roles []UserRole
	err := mysqlmanager.MySQL().Order("id desc").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// 新增角色
func AddRole(r *UserRole) error {
	if r.Role == "" {
		return fmt.Errorf("角色名不能为空")
	}

	return mysqlmanager.MySQL().Save(r).Error
}

// 是否有用户绑定某角色
func IsUserBindingRole(roleId int) (bool, error) {
	var users []User
	err := mysqlmanager.MySQL().Find(&users).Error
	if err != nil {
		return false, err
	}
	for _, user := range users {
		id := user.RoleID
		if find := strings.Contains(id, strconv.Itoa(roleId)); find {
			return true, nil
		}
	}
	return false, nil
}

// 删除用户角色
func DeleteRole(role string) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Where("role = ?", role).Unscoped().Delete(UserRole).Error
}

// 修改角色名称
func UpdateRoleName(roleId int, name string) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Model(&UserRole).Where("id = ?", roleId).Update("role", name).Error
}

// 修改角色描述
func UpdateRoleDescription(role, desc string) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Model(&UserRole).Where("role = ?", role).Update("description", desc).Error
}

// 变更用户角色权限
func UpdateRolePermission(permission *Frontdata) (UserRole, error) {
	var userRole UserRole
	// 数组切片转为string
	menus := strings.Replace(strings.Trim(fmt.Sprint(permission.Menus), "[]"), " ", ",", -1)
	buttonId := strings.Replace(strings.Trim(fmt.Sprint(permission.ButtonId), "[]"), " ", ",", -1)

	r := UserRole{
		Menus:    menus,
		ButtonID: buttonId,
	}
	err := mysqlmanager.MySQL().Model(&userRole).Where("id = ?", permission.Role_roleid).Updates(&r).Error
	return userRole, err
}

// 创建管理员账户
func CreateAdministratorUser() error {
	var role UserRole
	mysqlmanager.MySQL().Where("type =?", global.AdminUserType).Find(&role)
	if role.ID == 0 {
		role = UserRole{
			Role:        "admin",
			Type:        global.AdminUserType,
			Description: "超级管理员",
			Menus:       global.PILOTGO_MENUS,
			ButtonID:    global.PILOTGO_BUTTONID,
		}
		mysqlmanager.MySQL().Create(&role)
		bs, err := utils.CryptoPassword(strings.Split(global.SuperUser, "@")[0])
		if err != nil {
			return err
		}

		user := User{
			CreatedAt:    time.Time{},
			DepartFirst:  global.Departroot,
			DepartSecond: global.UncateloguedDepartId,
			DepartName:   "超级用户",
			Username:     strings.Split(global.SuperUser, "@")[0],
			Password:     string(bs),
			Email:        global.SuperUser,
			UserType:     global.AdminUserType,
			RoleID:       strconv.Itoa(role.ID),
		}
		mysqlmanager.MySQL().Create(&user)
	}

	return nil
}
