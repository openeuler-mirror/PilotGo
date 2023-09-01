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
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
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

func GetLoginUserPermission(Roleid RoleID) (dao.UserRole, interface{}, error) {
	roleId := RoleId(Roleid) //用户的最高权限
	userRole, err := dao.RoleIdToGetAllInfo(roleId)
	if err != nil {
		return userRole, nil, err
	}
	buttons, err := dao.PermissionButtons(userRole.ButtonID)
	if err != nil {
		return userRole, buttons, err
	}
	return userRole, buttons, nil
}

func GetRoles(query *common.PaginationQ) (int, interface{}, error) {
	roles, total, err := dao.GetAllRoles()
	if err != nil {
		return total, nil, err
	}
	data, err := common.DataPaging(query, roles, total)
	if err != nil {
		return total, data, err
	}
	return total, data, nil
}

func RolePermissionChangeMethod(roleChange *dao.Frontdata) (*dao.UserRole, error) {
	userRole, err := dao.UpdateRolePermission(roleChange)
	if err != nil {
		return &userRole, err
	}
	return &userRole, nil
}
