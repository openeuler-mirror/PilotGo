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
package service

import (
	"errors"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
)

// 获取用户最高权限的角色id
func RoleId(R model.RoleID) int {
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

func GetLoginUserPermission(Roleid model.RoleID) (model.UserRole, interface{}, error) {
	roleId := RoleId(Roleid) //用户的最高权限
	userRole := dao.RoleIdToGetAllInfo(roleId)
	buttons := dao.PermissionButtons(userRole.ButtonID)
	return userRole, buttons, nil
}

func GetRoles(query *model.PaginationQ) (int, interface{}, error) {
	roles, total := dao.GetAllRoles()
	data, err := DataPaging(query, roles, total)
	if err != nil {
		return total, data, err
	}
	return total, data, nil
}

func AddUserRole(userRole *model.UserRole) error {
	err := dao.AddRole(*userRole)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserRole(ID int) error {
	if ok := dao.IsUserBindingRole(ID); !ok {
		dao.DeleteRole(ID)
		return nil
	} else {
		return errors.New("有用户绑定此角色，不可删除")
	}
}

func UpdateUserRole(UserRole *model.UserRole) error {
	id := UserRole.ID
	role := UserRole.Role
	description := UserRole.Description
	userRole := dao.RoleIdToGetAllInfo(id)
	if userRole.Role != role && userRole.Description != description {
		dao.UpdateRoleName(id, role)
		dao.UpdateRoleDescription(id, description)
		return nil
	}
	if userRole.Role == role && userRole.Description != description {
		dao.UpdateRoleDescription(id, description)
		return nil
	}
	if userRole.Role != role && userRole.Description == description {
		dao.UpdateRoleName(id, role)
		return nil
	}
	return errors.New("没有修改信息")
}

func RolePermissionChange(roleChange model.RolePermissionChange) (*model.UserRole, error) {
	userRole := dao.UpdateRolePermission(roleChange)
	return &userRole, nil
}
