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
 * Date: 2022-03-21 15:32:50
 * LastEditTime: 2022-04-21 15:37:48
 * Description: 用户模块逻辑代码
 ******************************************************************************/
package service

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
)

// 随机产生用户名字
func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for k := range result {
		result[k] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 判断用户类型
func UserType(s string) int {
	// 找出字符串中包含的最小数字，例：“1,2,3,4”，最小的是1
	roleIds := strings.Split(s, ",")
	res := make([]int, len(s))

	for k, v := range roleIds {
		res[k], _ = strconv.Atoi(v)
	}

	min := res[0]
	if len(res) > 1 {
		for _, v := range res {
			if v < min {
				min = v
			}
		}
	}

	var user_type int
	if min > global.OrdinaryUserRoleId {
		user_type = global.OtherUserType
	} else {
		user_type = min - 1
	}
	return user_type
}

// 读取xlsx文件
func ReadFile(xlFile *xlsx.File, UserExit []string) []string {
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			//跳过第一行表头信息
			if rowIndex == 0 {
				continue
			}
			userName := row.Cells[0].Value //1:用户名
			phone := row.Cells[1].Value    //2：手机号
			email := row.Cells[2].Value    //3：邮箱
			if dao.IsEmailExist(email) {
				UserExit = append(UserExit, email)
				continue
			}
			departName := row.Cells[3].Value       //4：部门
			pid, id := dao.GetPidAndId(departName) // 部门对应的PId和Id

			userRole := row.Cells[4].Value                          // 5：角色
			roleId, user_type := dao.GetRoleIdAndUserType(userRole) //角色对应id和用户类型
			password := global.DefaultUserPassword                  // 设置默认密码为123456
			u := model.User{
				Username:     userName,
				Phone:        phone,
				Password:     password,
				Email:        email,
				DepartName:   departName, //4：部门
				DepartFirst:  pid,
				DepartSecond: id,
				UserType:     user_type,
				RoleID:       roleId,
			}
			dao.AddUser(u)
		}
	}
	return UserExit
}
