package model

/**
 * @Author: zhang han
 * @Date: 2021/10/28 14:23
 * @Description: mysql数据结构体
 */

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(25);not null" json:"username,omitempty" form:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone    string `gorm:"size:255" json:"phone,omitempty" form:"phone"`
	Email    string `gorm:"type:varchar(30);not null" json:"email,omitempty" form:"email"`
	Enable   string `gorm:"size:10;not null" json:"enable,omitempty"`
}
