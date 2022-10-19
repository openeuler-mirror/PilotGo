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

import "openeuler.org/PilotGo/PilotGo/pkg/app/server/model"

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
