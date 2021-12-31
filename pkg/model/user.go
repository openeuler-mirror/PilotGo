package model

/**
 * @Author: zhang han
 * @Date: 2021/10/28 14:23
 * @Description: mysql数据结构体
 */

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(25);not null" json:"username,omitempty" form:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password,omitempty"`
	Phone    string `gorm:"size:255" json:"phone,omitempty" form:"phone"`
	Email    string `gorm:"type:varchar(30);not null" json:"email,omitempty" form:"email"`
	Enable   string `gorm:"size:10;not null" json:"enable,omitempty"`
}

type UserQ struct {
	User
	PaginationQ
}

func (u *User) All(q *PaginationQ) (list *[]User, total uint, err error) {
	list = &[]User{}
	tx := mysqlmanager.DB.Find(list)
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *User) Update() (err error) {
	m.makePassword()
	return mysqlmanager.DB.Model(m).Update(m).Error
}

func (m *User) makePassword() {
	if m.Password != "" {
		if bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost); err != nil {
			logrus.WithError(err).Error("bcrypt making password is failed")
		} else {
			m.Password = string(bytes)
		}
	}
}
