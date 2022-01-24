package dao

/**
 * @Author: zhang han
 * @Date: 2021/10/30 15:22
 * @Description:
 */

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func IsEmailExist(email string) bool {
	var user model.User
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
