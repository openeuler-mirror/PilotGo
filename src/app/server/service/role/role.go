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
	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
)

type UserRole = dao.UserRole

type RoleID struct {
	RoleId []int `json:"roleId"`
}

// 获取用户最高权限的角色id
func RoleId(R RoleID) int {
	min := R.RoleId[0]
	if len(R.RoleId) > 1 {
		for _, v := range R.RoleId {
			if v < min {
				min = v
			}
		}
	}
	return min
}

// return menu, button, error
func GetLoginUserPermission(Roleid RoleID) (string, []string, error) {
	// TODO: multi role case
	roleId := RoleId(Roleid) //用户的最高权限
	role, err := dao.RoleIdToGetAllInfo(roleId)
	if err != nil {
		return "", nil, err
	}

	menu, buttons := getRoleMenuButtons(role.Role)

	return menu, buttons, nil
}

type ReturnRole struct {
	ID          int      `json:"id"`
	Role        string   `json:"role"`
	Type        int      `json:"type"`
	Description string   `json:"description"`
	Menus       string   `json:"menus"`
	Buttons     []string `json:"buttons"`
}

func GetRoles() ([]*ReturnRole, error) {
	result := []*ReturnRole{}
	roles, err := dao.GetRoleList()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		result = append(result, &ReturnRole{
			ID:          role.ID,
			Role:        role.Role,
			Description: role.Description,
			Menus:       "",
			Buttons:     []string{},
		})
	}

	policies := auth.GetAllPolicies()

	for _, item := range result {
		for _, p := range policies {
			if item.Role == p.Role {
				switch p.Action {
				case "menu":
					item.Menus = item.Menus + "," + p.Resource
				case "button":
					item.Buttons = append(item.Buttons, p.Resource)
				}
			}
		}
		if len(item.Menus) > 0 {
			item.Menus = item.Menus[1:]
		}
	}

	return result, nil
}

func getRoleMenuButtons(role string) (string, []string) {
	menu := ""
	buttons := []string{}

	policies := auth.GetAllPolicies()
	for _, p := range policies {
		if role == p.Role {
			switch p.Action {
			case "button":
				buttons = append(buttons, p.Resource)
			case "menu":
				menu = menu + "," + p.Resource
			}
		}
	}

	if len(menu) > 0 {
		menu = menu[1:]
	}

	return menu, buttons
}

func AddRole(userRole *UserRole) error {
	err := dao.AddRole(userRole)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRole(role string) error {
	if err := auth.DeleteRole(role); err != nil {
		return err
	}

	err := dao.DeleteRole(role)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRoleInfo(role, description string) error {
	return dao.UpdateRoleDescription(role, description)
}

func UpdateRolePermissions(role string, buttons, menus []string) error {
	return auth.UpdateRolePermissions(role, buttons, menus)
}
