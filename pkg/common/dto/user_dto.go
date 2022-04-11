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
 * Date: 2021-11-24 15:08:08
 * LastEditTime: 2022-04-11 09:33:49
 * Description: user struct reflect
 ******************************************************************************/
package dto

import "openeluer.org/PilotGo/PilotGo/pkg/app/server/model"

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
	}
}
