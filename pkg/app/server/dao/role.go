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
 * LastEditTime: 2022-03-16 15:25:41
 * Description: 角色模块相关数据获取
 ******************************************************************************/
package dao

import (
	"strconv"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

// 根据角色名称返回角色id和用户类型
func GetRoleIdAndUserType(role string) (roleId string, user_type int) {
	var Role model.UserRole
	mysqlmanager.DB.Where("role = ?", role).Find(&Role)
	roleID := strconv.Itoa(Role.ID)
	var userType int
	if Role.ID > 3 {
		userType = 3
	} else {
		userType = Role.ID - 1
	}
	return roleID, userType
}
