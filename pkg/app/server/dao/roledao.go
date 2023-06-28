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

type ReturnUserRole struct {
	ID          int      `json:"id"`
	Role        string   `json:"role"`
	Type        int      `json:"type"`
	Description string   `json:"description"`
	Menus       string   `json:"menus"`
	Buttons     []string `json:"buttons"`
}

type RoleButton struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Button string `json:"button"`
}

type RolePermissionChange struct {
	RoleID   int      `json:"id"`
	Menus    []string `json:"menus"`
	ButtonId []string `json:"buttonId"`
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
func PermissionButtons(button string) (interface{}, error) {
	var buttons []string
	if len(button) == 0 {
		return []interface{}{}, nil
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
func GetAllRoles() ([]ReturnUserRole, int, error) {
	var roles []UserRole
	var getRole []ReturnUserRole
	err := mysqlmanager.MySQL().Order("id desc").Find(&roles).Error
	if err != nil {
		return getRole, 0, err
	}
	total := len(roles)

	for _, role := range roles {
		var buts []string

		if len(role.ButtonID) == 0 {
			r := ReturnUserRole{
				ID:          role.ID,
				Role:        role.Role,
				Type:        role.Type,
				Description: role.Description,
				Menus:       role.Menus,
				Buttons:     []string{},
			}
			getRole = append(getRole, r)
			continue
		} else {
			buttonss := strings.Split(role.ButtonID, ",")

			for _, button := range buttonss {
				var but RoleButton
				i, _ := strconv.Atoi(button)
				err := mysqlmanager.MySQL().Where("id=?", i).Find(&but).Error
				if err != nil {
					return getRole, total, err
				}
				buts = append(buts, but.Button)
			}
			r := ReturnUserRole{
				ID:          role.ID,
				Role:        role.Role,
				Type:        role.Type,
				Description: role.Description,
				Menus:       role.Menus,
				Buttons:     buts,
			}
			getRole = append(getRole, r)
		}
	}
	return getRole, total, nil
}

// 新增角色
func AddRole(r UserRole) error {
	role := r.Role
	if len(role) == 0 {
		return fmt.Errorf("用户角色不能为空")
	}
	userRole := UserRole{
		Role:        role,
		Type:        r.Type,
		Description: r.Description,
	}
	return mysqlmanager.MySQL().Save(&userRole).Error
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
func DeleteRole(roleId int) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Where("id = ?", roleId).Unscoped().Delete(UserRole).Error
}

// 修改角色名称
func UpdateRoleName(roleId int, name string) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Model(&UserRole).Where("id = ?", roleId).Update("role", name).Error
}

// 修改角色描述
func UpdateRoleDescription(roleId int, desc string) error {
	var UserRole UserRole
	return mysqlmanager.MySQL().Model(&UserRole).Where("id = ?", roleId).Update("description", desc).Error
}

// 变更用户角色权限
func UpdateRolePermission(permission RolePermissionChange) (UserRole, error) {
	var userRole UserRole
	// 数组切片转为string
	menus := strings.Replace(strings.Trim(fmt.Sprint(permission.Menus), "[]"), " ", ",", -1)
	buttonId := strings.Replace(strings.Trim(fmt.Sprint(permission.ButtonId), "[]"), " ", ",", -1)

	r := UserRole{
		Menus:    menus,
		ButtonID: buttonId,
	}
	err := mysqlmanager.MySQL().Model(&userRole).Where("id = ?", permission.RoleID).Updates(&r).Error
	return userRole, err
}

// 创建管理员账户
func CreateAdministratorUser() error {
	var role UserRole
	mysqlmanager.MySQL().Where("type =?", global.AdminUserType).Find(&role)
	if role.ID == 0 {
		role = UserRole{
			Role:        "超级用户",
			Type:        global.AdminUserType,
			Description: "超级管理员",
			Menus:       global.PILOTGO_MENUS,
			ButtonID:    global.PILOTGO_BUTTONID,
		}
		mysqlmanager.MySQL().Create(&role)
		bs, err := utils.CryptoPassword(global.DefaultUserPassword)
		if err != nil {
			return err
		}

		user := User{
			CreatedAt:    time.Time{},
			DepartFirst:  global.Departroot,
			DepartSecond: global.UncateloguedDepartId,
			DepartName:   "超级用户",
			Username:     "admin",
			Password:     string(bs),
			Email:        "admin@123.com",
			UserType:     global.AdminUserType,
			RoleID:       strconv.Itoa(role.ID),
		}
		mysqlmanager.MySQL().Create(&user)
	}

	var roleButton RoleButton
	mysqlmanager.MySQL().First(&roleButton)
	if roleButton.ID == 0 {
		mysqlmanager.MySQL().Raw("INSERT INTO role_button(id, button)" +
			"VALUES" +
			"('1', 'rpm_install')," +
			"('2', 'rpm_uninstall')," +
			"('3', 'batch_update')," +
			"('4', 'batch_delete')," +
			"('5', 'user_add')," +
			"('6', 'user_import')," +
			"('7', 'user_edit')," +
			"('8', 'user_reset')," +
			"('9', 'user_del')," +
			"('10', 'role_add')," +
			"('11', 'role_update')," +
			"('12', 'role_delete')," +
			"('13', 'role_modify')," +
			"('14', 'config_install')," +
			"('15', 'dept_change')").Scan(&roleButton)
	}
	return nil
}
