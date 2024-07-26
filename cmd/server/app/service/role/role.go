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
 * Date: 2022-04-27 15:32:50
 * LastEditTime: 2022-04-27 17:17:48
 * Description: 用户角色逻辑代码
 ******************************************************************************/
package role

import (
	"errors"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
)

type Role = dao.Role
type UserRole = dao.UserRole

// 获取用户最高权限的角色id
func RoleId(R []int) int {
	min := R[0]
	if len(R) > 1 {
		for _, v := range R {
			if v < min {
				min = v
			}
		}
	}
	return min
}

// return menu, button, error
func GetLoginUserPermission(Roleid []int) (map[string]interface{}, error) {
	// TODO: multi role case
	roleId := RoleId(Roleid) //用户的最高权限
	role, err := dao.GetRoleById(roleId)
	if err != nil {
		return nil, err
	}
	permissions := getRoleMenuButtons(role.Name)
	return permissions, nil
}

type ReturnRole struct {
	ID          int                    `json:"id"`
	Role        string                 `json:"role"`
	Type        int                    `json:"type"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
}

func GetRoles() ([]*ReturnRole, error) {
	result := []*ReturnRole{}
	roles, err := dao.GetRoles()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		permissions := getRoleMenuButtons(role.Name)

		result = append(result, &ReturnRole{
			ID:          role.ID,
			Role:        role.Name,
			Description: role.Description,
			Permissions: permissions,
		})
	}

	return result, nil
}

func getRoleMenuButtons(role string) map[string]interface{} {
	permissions := make(map[string]interface{})
	//获取pilotgo-server的权限信息
	menu := ""
	buttons := []string{}
	policys := auth.GetFilteredPolicy(role, "", "button", "")
	for _, v := range policys {
		buttons = append(buttons, v[1])
	}
	permissions["button"] = buttons

	policys = auth.GetFilteredPolicy(role, "", "menu", "")
	for _, v := range policys {
		menu = menu + "," + v[1]
	}
	if len(menu) > 0 {
		menu = menu[1:]
	}
	permissions["menu"] = menu

	//遍历查询插件权限
	p := plugin.GetRolePluginPermission(role)
	for k, v := range p {
		permissions[k] = v
	}
	return permissions
}

// 添加角色
func AddRole(userRole *Role) error {
	if userRole.Name == "" {
		return errors.New("role name can not be empty")
	}

	id, _ := dao.GetRoleId(userRole.Name)
	if id > 0 {
		return errors.New("role name already exists")
	}
	err := dao.AddRole(userRole)
	if err != nil {
		return err
	}
	return nil
}

// 删除角色
func DeleteRole(roleId int) error {
	role, err := dao.GetRoleById(roleId)
	if err != nil {
		return err
	}
	if err := auth.DeleteRole(role.Name); err != nil {
		return err
	}

	err = dao.DeleteRole(roleId)
	if err != nil {
		return err
	}
	return nil
}

// 更改角色
func UpdateRoleInfo(name, description string) error {
	return dao.UpdateRoleDescription(name, description)
}

// 更改角色权限
func UpdateRolePermissions(role string, buttons, menus []string, PluginPermissions []plugin.PluginPermission) error {
	if role == "admin" {
		return errors.New("admin角色权限不可修改")
	}

	if err := auth.DeleteRole(role); err != nil {
		return err
	}
	if err := auth.UpdateRolePermissions(role, buttons, menus); err != nil {
		return err
	}

	return plugin.UpdatePluginPermissions(role, PluginPermissions)
}

// 分页查询
func GetRolePaged(offset, size int) (int64, []*ReturnRole, error) {
	result := []*ReturnRole{}
	total, data, err := dao.GetRolePaged(offset, size)
	if err != nil {
		return 0, nil, err
	}

	for _, role := range data {
		permissions := getRoleMenuButtons(role.Name)
		result = append(result, &ReturnRole{
			ID:          role.ID,
			Role:        role.Name,
			Description: role.Description,
			Permissions: permissions,
		})
	}

	return total, result, err
}
