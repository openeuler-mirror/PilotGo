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
 * LastEditTime: 2023-09-04 14:02:00
 * Description: 用户模块相关数据获取
 ******************************************************************************/
package dao

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/utils"
)

type User struct {
	ID           uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt    time.Time
	DepartFirst  int    `gorm:"size:25" json:"departPid,omitempty"`
	DepartSecond int    `gorm:"size:25" json:"departId,omitempty"`
	DepartName   string `gorm:"size:25" json:"departName,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone        string `gorm:"size:11" json:"phone,omitempty"`
	Email        string `gorm:"type:varchar(30);not null" json:"email,omitempty"`
	RoleID       string `json:"roleid,omitempty"`
	// deprecated
	// UserType int    `json:"userType,omitempty"`
}
type ReturnUser struct {
	ID           uint     `json:"id"`
	DepartFirst  int      `json:"departPId"`
	DepartSecond int      `json:"departid"`
	DepartName   string   `json:"departName"`
	Username     string   `json:"username"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Roles        []string `json:"role"`
	// deprecated
	// UserType int      `json:"userType"`
}
type UserDto struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// 邮箱账户是否存在
func IsEmailExist(email string) (bool, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	return user.ID != 0, err
}

// 查询某用户信息
func UserInfo(email string) (User, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	return user, err
}

// 查询某用户信息
func QueryUserByID(userID int) (*User, error) {
	user := &User{}
	err := mysqlmanager.MySQL().Where("id=?", strconv.Itoa(userID)).Find(user).Error
	return user, err
}

func QueryUserByName(name string) (*User, error) {
	user := &User{}
	err := mysqlmanager.MySQL().Where("username=?", name).Find(user).Error
	return user, err
}

// 分页查询所有用户
func GetUserPaged(offset, size int) (int64, []ReturnUser, error) {
	var users []User
	var returnUsers []ReturnUser
	var count int64
	err := mysqlmanager.MySQL().Model(User{}).Order("id desc").Offset(offset).Limit(size).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return count, returnUsers, err
	}
	for _, user := range users {
		var roles []string
		// 查找角色
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		for _, id := range roleId {
			userRole := UserRole{}
			i, _ := strconv.Atoi(id)
			err := mysqlmanager.MySQL().Where("id = ?", i).Find(&userRole).Error
			if err != nil {
				return count, returnUsers, err
			}
			role := userRole.Role
			roles = append(roles, role)
		}
		u := ReturnUser{
			ID:           user.ID,
			DepartFirst:  user.DepartFirst,
			DepartSecond: user.DepartSecond,
			DepartName:   user.DepartName,
			Username:     user.Username,
			Phone:        user.Phone,
			Email:        user.Email,
			Roles:        roles,
		}
		returnUsers = append(returnUsers, u)
	}
	return count, returnUsers, err
}

// 查询所有的用户
func UserAll() ([]ReturnUser, int, error) {
	var users []User
	var redisUser []ReturnUser

	err := mysqlmanager.MySQL().Order("id desc").Find(&users).Error
	if err != nil {
		return redisUser, 0, err
	}
	totals := len(users)
	for _, user := range users {
		var roles []string
		// 查找角色
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		for _, id := range roleId {
			userRole := UserRole{}
			i, _ := strconv.Atoi(id)
			err := mysqlmanager.MySQL().Where("id = ?", i).Find(&userRole).Error
			if err != nil {
				return redisUser, totals, err
			}
			role := userRole.Role
			roles = append(roles, role)
		}
		u := ReturnUser{
			ID:           user.ID,
			DepartFirst:  user.DepartFirst,
			DepartSecond: user.DepartSecond,
			DepartName:   user.DepartName,
			Username:     user.Username,
			Phone:        user.Phone,
			Email:        user.Email,
			Roles:        roles,
		}
		redisUser = append(redisUser, u)
	}
	return redisUser, totals, nil
}

// 根据用户邮箱模糊查询
func UserSearchPaged(email string, offset, size int) (int64, []ReturnUser, error) {
	var users []User
	var redisUser []ReturnUser
	var count int64
	err := mysqlmanager.MySQL().Order("id desc").Where("email LIKE ?", "%"+email+"%").Offset(offset).Limit(size).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return 0, redisUser, err
	}
	for _, user := range users {
		var roles []string
		// 查找角色
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		for _, id := range roleId {
			userRole := UserRole{}
			i, _ := strconv.Atoi(id)
			err := mysqlmanager.MySQL().Where("id = ?", i).Find(&userRole).Error
			if err != nil {
				return count, redisUser, err
			}
			role := userRole.Role
			roles = append(roles, role)
		}
		u := ReturnUser{
			ID:           user.ID,
			DepartFirst:  user.DepartFirst,
			DepartSecond: user.DepartSecond,
			DepartName:   user.DepartName,
			Username:     user.Username,
			Phone:        user.Phone,
			Email:        user.Email,
			Roles:        roles,
		}
		redisUser = append(redisUser, u)
	}
	return count, redisUser, nil
}

// 修改密码
func UpdatePassword(email, newPWD string) (User, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	if err != nil {
		return user, err
	} else {
		err = utils.ComparePassword(user.Password, newPWD)
		if err == nil {
			return user, errors.New("新密码和旧密码一致")
		}

		bs, err := utils.CryptoPassword(newPWD)
		if err != nil {
			return user, err
		}
		err = mysqlmanager.MySQL().Model(&user).Where("email=?", email).Update("password", string(bs)).Error
		return user, err
	}
}

// 重置密码
func ResetPassword(email string) (User, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	if err != nil {
		return user, err
	} else {
		bs, err := utils.CryptoPassword(strings.Split(email, "@")[0])
		if err != nil {
			return user, err
		}
		err = mysqlmanager.MySQL().Model(&user).Where("email=?", email).Update("password", string(bs)).Error
		return user, err
	}
}

// 删除用户
func DeleteUser(email string) error {
	var user User
	return mysqlmanager.MySQL().Where("email=?", email).Unscoped().Delete(user).Error
}

// 修改用户的部门信息
func UpdateUserDepart(email, departName string, Pid, id int) error {
	var user User
	u := User{
		DepartFirst:  Pid,
		DepartSecond: id,
		DepartName:   departName,
	}
	return mysqlmanager.MySQL().Model(&user).Where("email=?", email).Updates(&u).Error
}

// 添加用户
func AddUser(u *User) error {
	return mysqlmanager.MySQL().Save(u).Error
}

// 修改手机号
func UpdateUserPhone(email, phone string) error {
	var user User
	return mysqlmanager.MySQL().Model(&user).Where("email=?", email).Update("phone", phone).Error
}

func DelUser(deptId int) error {
	var user User
	return mysqlmanager.MySQL().Where("depart_second=?", deptId).Unscoped().Delete(user).Error
}