/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package role

import (
	"errors"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/common"
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
func GetLoginUserPermission(Roleid []int) (interface{}, error) {
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
	ID          int                 `json:"id"`
	Role        string              `json:"role"`
	Type        int                 `json:"type"`
	Description string              `json:"description"`
	Permissions map[string][]string `json:"permissions"`
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

func getRoleMenuButtons(role string) map[string][]string {
	permissions := make(map[string][]string)
	//获取pilotgo-server的权限信息
	menu := []string{}
	buttons := []string{}
	policys := auth.GetFilteredPolicy(role, "", "button", auth.DomainPilotGo)
	for _, v := range policys {
		buttons = append(buttons, v[1])
	}
	permissions["button"] = buttons

	policys = auth.GetFilteredPolicy(role, "", "menu", auth.DomainPilotGo)
	for _, v := range policys {
		menu = append(menu, v[1])
	}
	permissions["menu"] = menu

	//遍历查询插件权限
	p := plugin.GetRolePluginPermission(role)
	for k, v := range p {
		permissions[k] = append(permissions[k], v...)
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
func UpdateRolePermissions(role string, buttons, menus []string) error {
	if role == "admin" {
		return errors.New("admin角色权限不可修改")
	}

	if err := auth.DeleteRole(role); err != nil {
		return err
	}
	ms, bs, pms := BuildAllPermissions(menus, buttons)
	if err := auth.UpdateRolePermissions(role, bs, ms); err != nil {
		return err
	}

	return plugin.UpdatePluginPermissions(role, pms)
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
func BuildAllPermissions(menus, buttons []string) ([]string, []string, []plugin.PluginPermission) {
	resultMap := make(map[string][]common.Permission)
	serverMenus := []string{}
	serverButtons := []string{}

	addResource := func(res string, isMenu bool) {
		if domain, action := auth.PermissionListMap.FindByInnerKey(res); domain != "" && action != "" {
			if domain == auth.DomainPilotGo {
				if isMenu {
					serverMenus = append(serverMenus, res)
				} else {
					serverButtons = append(serverButtons, res)
				}
			} else {
				resultMap[domain] = append(resultMap[domain], common.Permission{
					Resource: res,
					Operate:  action,
				})
			}

		}
	}

	for _, res := range menus {
		addResource(res, true)
	}

	for _, res := range buttons {
		addResource(res, false)
	}

	var result []plugin.PluginPermission
	for serviceName, perms := range resultMap {
		result = append(result, plugin.PluginPermission{
			ServiceName: serviceName,
			Permissions: perms,
		})
	}

	return serverMenus, serverButtons, result
}
