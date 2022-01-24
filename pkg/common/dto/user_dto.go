package dto

import "openeluer.org/PilotGo/PilotGo/pkg/app/server/model"

/**
 * @Author: zhang han
 * @Date: 2021/11/1 14:44
 * @Description:
 */

type UserDto struct {
	Name     string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Enable   string `json:"enable,omitempty"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:     user.Username,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Enable:   user.Enable,
	}
}
