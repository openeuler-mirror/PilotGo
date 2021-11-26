package service

/**
 * @Author: zhang han
 * @Date: 2021/11/12 17:22
 * @Description:
 */

import (
	"openeluer.org/PilotGo/PilotGo/pkg/db"
	"openeluer.org/PilotGo/PilotGo/pkg/model"
)

func GetUserInfoList(info model.PageInfo) (err error, list interface{}, total int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var userList []model.User
	err = db.DB.Find(&userList).Count(&total).Error
	err = db.DB.Limit(limit).Offset(offset).Find(&userList).Error
	return err, userList, total
}
