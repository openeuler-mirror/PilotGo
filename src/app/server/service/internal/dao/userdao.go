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
	ID        uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time
	DepartId  int    `gorm:"size:25" json:"departId,omitempty"`
	Username  string `json:"username,omitempty" `
	Password  string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone     string `gorm:"size:11" json:"phone,omitempty"`
	Email     string `gorm:"type:varchar(30);not null;unique" json:"email,omitempty"`
}

// 邮箱账户是否存在
func IsEmailExist(email string) (bool, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	return user.ID != 0, err
}

// 查询某用户信息
func GetUserByEmail(email string) (User, error) {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	return user, err
}

// 查询某用户信息
func GetUserByID(userID int) (*User, error) {
	user := &User{}
	err := mysqlmanager.MySQL().Where("id=?", strconv.Itoa(userID)).Find(user).Error
	return user, err
}

func GetUserByName(name string) (*User, error) {
	user := &User{}
	err := mysqlmanager.MySQL().Where("username=?", name).Find(user).Error
	return user, err
}

func GetUserBypid(pid int) ([]User, error) {
	var users []User
	err := mysqlmanager.MySQL().Where("depart_id=?", pid).Find(&users).Error
	return users, err
}

// 分页查询所有用户
func GetUserPaged(offset, size int) (int64, []User, error) {
	var users []User
	var count int64
	err := mysqlmanager.MySQL().Model(User{}).Order("id desc").Offset(offset).Limit(size).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	return count, users, err
}

// 修改密码
func UpdatePassword(email, newPWD string) error {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	if err != nil {
		return err
	} else {
		err = utils.ComparePassword(user.Password, newPWD)
		if err == nil {
			return errors.New("新密码和旧密码一致")
		}

		bs, err := utils.CryptoPassword(newPWD)
		if err != nil {
			return err
		}
		err = mysqlmanager.MySQL().Model(&user).Where("email=?", email).Update("password", string(bs)).Error
		return err
	}
}

// 重置密码
func ResetPassword(email string) error {
	var user User
	err := mysqlmanager.MySQL().Where("email=?", email).Find(&user).Error
	if err != nil {
		return err
	} else {
		bs, err := utils.CryptoPassword(strings.Split(email, "@")[0])
		if err != nil {
			return err
		}
		err = mysqlmanager.MySQL().Model(&user).Where("email=?", email).Update("password", string(bs)).Error
		return err
	}
}

// 添加用户
func AddUser(u *User) error {
	return mysqlmanager.MySQL().Save(u).Error
}

// 修改用户信息
func UpdateUser(email string, u User) error {
	var user User
	return mysqlmanager.MySQL().Model(&user).Where("email=?", email).Updates(&u).Error
}

// 删除用户
func DeleteUser(email string) error {
	var user User
	return mysqlmanager.MySQL().Where("email=?", email).Unscoped().Delete(user).Error
}

// 根据用户邮箱模糊查询
func UserSearchPaged(email string, offset, size int) (int64, []User, error) {
	var users []User
	var count int64
	err := mysqlmanager.MySQL().Order("id desc").Where("email like ?", "%"+email+"%").Offset(offset).Limit(size).Find(&users).Offset(-1).Limit(-1).Count(&count).Error
	return count, users, err
}
