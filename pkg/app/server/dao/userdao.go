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
 * Description: 用户模块相关数据获取
 ******************************************************************************/
package dao

import (
	"fmt"
	"strconv"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

// 获取所有的用户角色
func AllUserRole() []model.UserRole {
	var role []model.UserRole
	mysqlmanager.DB.Find(&role)
	return role
}

// 邮箱账户是否存在
func IsEmailExist(email string) bool {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	return user.ID != 0
}

// 查询数据库中账号密码、用户部门、部门ID、用户类型、用户角色
func UserPassword(email string) (s1, s2, s3 string, i1, i2 int) {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	return user.Password, user.DepartName, user.RoleID, user.DepartSecond, user.UserType
}

// 查询某用户信息
func UserInfo(email string) model.User {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	return user
}

// 查询所有的用户
func UserAll() []map[string]interface{} {
	var users []model.User
	mysqlmanager.DB.Find(&users)
	datas := make([]map[string]interface{}, 0)
	for _, user := range users {
		data := make(map[string]interface{})
		data["id"] = user.ID
		data["departPId"] = user.DepartFirst
		data["departid"] = user.DepartSecond
		data["departName"] = user.DepartName
		data["username"] = user.Username
		data["phone"] = user.Phone
		data["email"] = user.Email
		data["userType"] = user.UserType
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		var roles []string
		for _, id := range roleId {
			userRole := model.UserRole{}
			i, _ := strconv.Atoi(id)
			mysqlmanager.DB.Where("id = ?", i).Find(&userRole)
			role := userRole.Role
			roles = append(roles, role)
		}
		data["role"] = roles
		datas = append(datas, data)
	}
	return datas
}

// 根据用户邮箱模糊查询
func UserSearch(email string) []map[string]interface{} {
	var users []model.User
	mysqlmanager.DB.Where("email LIKE ?", "%"+email+"%").Find(&users)

	datas := make([]map[string]interface{}, 0)
	for _, user := range users {
		data := make(map[string]interface{})
		data["id"] = user.ID
		data["departPId"] = user.DepartFirst
		data["departid"] = user.DepartSecond
		data["departName"] = user.DepartName
		data["username"] = user.Username
		data["phone"] = user.Phone
		data["email"] = user.Email
		data["userType"] = user.UserType
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		var roles []string
		for _, id := range roleId {
			userRole := model.UserRole{}
			i, _ := strconv.Atoi(id)
			mysqlmanager.DB.Where("id = ?", i).Find(&userRole)
			role := userRole.Role
			roles = append(roles, role)
		}
		data["role"] = roles
		datas = append(datas, data)
	}
	return datas
}

// 重置密码
func ResetPassword(email string) (model.User, error) {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	if user.ID != 0 {
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("password", "123456")
		return user, nil
	} else {
		return user, fmt.Errorf("无此用户")
	}
}

// 删除用户
func DeleteUser(email string) {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Unscoped().Delete(user)
}

// 修改用户的部门信息
func UpdateUserDepart(email, departName string, Pid, id int) {
	var user model.User
	u := model.User{
		DepartFirst:  Pid,
		DepartSecond: id,
		DepartName:   departName,
	}
	mysqlmanager.DB.Model(&user).Where("email=?", email).Updates(&u)
}

// 添加用户
func AddUser(u model.User) {
	mysqlmanager.DB.Save(&u)
}

// 修改手机号
func UpdateUserPhone(email, phone string) {
	var user model.User
	mysqlmanager.DB.Model(&user).Where("email=?", email).Update("phone", phone)
}
