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

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

// 根据角色名称返回角色id和用户类型
func GetRoleIdAndUserType(role string) (roleId string, user_type int, err error) {
	var Role model.UserRole
	err = global.PILOTGO_DB.Where("role = ?", role).Find(&Role).Error
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
func RoleIdToGetAllInfo(roleid int) (model.UserRole, error) {
	var role model.UserRole
	err := global.PILOTGO_DB.Where("id=?", roleid).Find(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

// 登录用户的权限按钮
func PermissionButtons(button string) (interface{}, error) {
	var buttons []string
	if len(button) == 0 {
		return []interface{}{}, nil
	}
	IDs := strings.Split(button, ",")

	for _, id := range IDs {
		var SubButton model.RoleButton
		i, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		err = global.PILOTGO_DB.Where("id = ?", i).Find(&SubButton).Error
		if err != nil {
			return buttons, err
		}
		button := SubButton.Button
		buttons = append(buttons, button)
	}
	return buttons, nil
}

// 获取所有的用户角色
func GetAllRoles() ([]model.ReturnUserRole, int, error) {
	var roles []model.UserRole
	var getRole []model.ReturnUserRole
	err := global.PILOTGO_DB.Order("id desc").Find(&roles).Error
	if err != nil {
		return getRole, 0, err
	}
	total := len(roles)

	for _, role := range roles {
		var buts []string

		if len(role.ButtonID) == 0 {
			r := model.ReturnUserRole{
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
				var but model.RoleButton
				i, _ := strconv.Atoi(button)
				err := global.PILOTGO_DB.Where("id=?", i).Find(&but).Error
				if err != nil {
					return getRole, total, err
				}
				buts = append(buts, but.Button)
			}
			r := model.ReturnUserRole{
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
func AddRole(r model.UserRole) error {
	role := r.Role
	if len(role) == 0 {
		return fmt.Errorf("用户角色不能为空")
	}
	userRole := model.UserRole{
		Role:        role,
		Type:        r.Type,
		Description: r.Description,
	}
	return global.PILOTGO_DB.Save(&userRole).Error
}

// 是否有用户绑定某角色
func IsUserBindingRole(roleId int) (bool, error) {
	var users []model.User
	err := global.PILOTGO_DB.Find(&users).Error
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
	var UserRole model.UserRole
	err := global.PILOTGO_DB.Where("id = ?", roleId).Unscoped().Delete(UserRole).Error
	if err != nil {
		return err
	}
	return nil
}

// 修改角色名称
func UpdateRoleName(roleId int, name string) error {
	var UserRole model.UserRole
	err := global.PILOTGO_DB.Model(&UserRole).Where("id = ?", roleId).Update("role", name).Error
	if err != nil {
		return err
	}
	return nil
}

// 修改角色描述
func UpdateRoleDescription(roleId int, desc string) error {
	var UserRole model.UserRole
	err := global.PILOTGO_DB.Model(&UserRole).Where("id = ?", roleId).Update("description", desc).Error
	if err != nil {
		return err
	}
	return nil
}

// 变更用户角色权限
func UpdateRolePermission(permission model.RolePermissionChange) (model.UserRole, error) {
	var userRole model.UserRole
	// 数组切片转为string
	menus := strings.Replace(strings.Trim(fmt.Sprint(permission.Menus), "[]"), " ", ",", -1)
	buttonId := strings.Replace(strings.Trim(fmt.Sprint(permission.ButtonId), "[]"), " ", ",", -1)

	r := model.UserRole{
		Menus:    menus,
		ButtonID: buttonId,
	}
	err := global.PILOTGO_DB.Model(&userRole).Where("id = ?", permission.RoleID).Updates(&r).Error
	if err != nil {
		return userRole, err
	}
	return userRole, nil
}

// 创建超级管理员账户
func CreateSuperAdministratorUser() {
	var user model.User
	var role model.UserRole
	var roleButton model.RoleButton
	global.PILOTGO_DB.Where("type =?", global.AdminUserType).Find(&role)
	if role.ID == 0 {
		role = model.UserRole{
			Role:        "超级用户",
			Type:        global.AdminUserType,
			Description: "超级管理员",
			Menus:       global.PILOTGO_MENUS,
			ButtonID:    global.PILOTGO_BUTTONID,
		}
		global.PILOTGO_DB.Create(&role)
		user = model.User{
			CreatedAt:    time.Time{},
			DepartFirst:  global.Departroot,
			DepartSecond: global.UncateloguedDepartId,
			DepartName:   "超级用户",
			Username:     "admin",
			Password:     global.DefaultUserPassword,
			Email:        "admin@123.com",
			UserType:     global.AdminUserType,
			RoleID:       strconv.Itoa(role.ID),
		}
		global.PILOTGO_DB.Create(&user)
	}
	global.PILOTGO_DB.First(&roleButton)
	if roleButton.ID == 0 {
		global.PILOTGO_DB.Raw("INSERT INTO role_button(id, button)" +
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
}
