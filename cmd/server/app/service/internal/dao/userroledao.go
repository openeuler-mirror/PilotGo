/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
)

type UserRole struct {
	ID     int  `gorm:"primary_key;AUTO_INCREMENT"`
	User   User `gorm:"Foreignkey:UserID"`
	UserID uint
	Role   Role `gorm:"Foreignkey:RoleID"`
	RoleID int
}

func (ur *UserRole) Add() error {
	return mysqlmanager.MySQL().Create(ur).Error
}

func GetRolesByUid(uid uint) ([]int, error) {
	var Roleid []int
	err := mysqlmanager.MySQL().Model(&UserRole{}).Select("role_id").Where("user_id=?", uid).Find(&Roleid).Error
	return Roleid, err
}

func DeleteByUid(uid uint) error {
	var urs UserRole
	return mysqlmanager.MySQL().Where("user_id = ?", uid).Unscoped().Delete(urs).Error
}

func UpdateU2R(uid uint, rids []int) error {
	//更新用户权限先删除旧权限再添加新权限
	err := DeleteByUid(uid)
	if err != nil {
		return err
	}
	for _, v := range rids {
		u2r := &UserRole{
			UserID: uid,
			RoleID: v,
		}
		err = u2r.Add()
	}
	return err
}
