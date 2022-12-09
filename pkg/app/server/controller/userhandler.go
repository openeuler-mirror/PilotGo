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
 * Date: 2021-12-18 02:33:55
 * LastEditTime: 2022-04-11 16:15:51
 * Description: 用户登录、增删改查
 ******************************************************************************/
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

var ServerAddr string

func GetUserRoleHandler(c *gin.Context) {
	roles := service.GetUserRole()
	response.Success(c, gin.H{"role": roles}, "获取用户角色")
}

func RegisterHandler(c *gin.Context) {
	var user model.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.Register(user)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "添加用户成功!") //Return result
}

func LoginHandler(c *gin.Context) {
	var user model.User //Data verification
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	token, departName, departId, userType, roleId, err := service.Login(user)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "userType": userType, "roleId": roleId, "server": ServerAddr}, "登陆成功!")
}

// 退出
func Logout(c *gin.Context) {

	response.Success(c, nil, "退出成功!")

}

func Info(c *gin.Context) {
	user, _ := c.Get("x-user")
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{"user": model.ToUserDto(user.(model.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	users, total := service.UserAll()
	data, err := service.DataPaging(query, users, total)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	service.JsonPagination(c, data, int64(total), query)
}

// 高级搜索
func UserSearchHandler(c *gin.Context) {
	var user model.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	var email = user.Email
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	data, total, err := service.UserSearch(email, query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	service.JsonPagination(c, data, int64(total), query)
}

// 重置密码
func ResetPasswordHandler(c *gin.Context) {
	var user model.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := service.ResetPassword(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
	} else {
		response.Success(c, gin.H{"data": u}, "密码重置成功!")
	}
}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	var userdel model.Userdel
	if c.Bind(&userdel) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.DeleteUser(userdel.Emails)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "用户删除成功!")
}

// 修改用户信息
func UpdateUserHandler(c *gin.Context) {
	var user model.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := service.UpdateUser(user)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"data": u}, "用户信息修改成功")

}

// 一键导入用户数据
func ImportUser(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload"]
	if files == nil {
		response.Fail(c, nil, "请先选择要上传的文件")
		return
	}
	UserExit := make([]string, 0)
	for _, file := range files {
		name := file.Filename
		c.SaveUploadedFile(file, name)
		xlFile, error := xlsx.OpenFile(name)
		if error != nil {
			return
		}
		UserExit = service.ReadFile(xlFile, UserExit)
	}

	if len(UserExit) == 0 {
		response.Success(c, nil, "导入用户信息成功")
	} else {
		response.Fail(c, gin.H{"UserExit": UserExit}, "以上用户已经存在")
	}
}
