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
 * Description: 用户模块相关数据验证
 ******************************************************************************/
package dao

import (
	"strconv"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func IsEmailExist(email string) bool {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	return user.ID != 0
}
func IsContain(str string, substr int) bool {
	strs := strings.Split(str, ",")
	for _, value := range strs {
		if value == strconv.Itoa(substr) {
			return true
		}
	}
	return false
}
