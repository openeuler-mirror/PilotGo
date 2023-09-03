package common

import "openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"

type User = dao.User

type SimpleUser struct {
	ID         uint   `json:"id"`
	DepartName string `json:"departName"`
	Username   string `json:"username"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	RoleID     string `json:"role"`
}
