package dao

import (
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
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

func DeleteByRid(rid int) error {
	var urs UserRole
	return mysqlmanager.MySQL().Where("role_id = ?", rid).Unscoped().Delete(urs).Error
}

func DeleteByUid(uid uint) error {
	var urs UserRole
	return mysqlmanager.MySQL().Where("user_id = ?", uid).Unscoped().Delete(urs).Error
}

// 是否有用户绑定某角色
func IsUserBindingRole(roleId int) (bool, error) {
	var Userid []int
	err := mysqlmanager.MySQL().Model(&UserRole{}).Select("user_id").Where("role_id=?", roleId).Find(&Userid).Error
	if len(Userid) > 0 {
		return true, err
	}
	return false, err
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
