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
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auditlog"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	userservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/user"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func GetUserRoleHandler(c *gin.Context) {
	roles, err := userservice.GetUserRole()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"role": roles}, "获取用户角色")
}

func RegisterHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user1 userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "添加用户", "", user1)
	auditlog.Add(log)

	err := userservice.Register(user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "添加用户成功!") //Return result
}

func LoginHandler(c *gin.Context) {
	var user userservice.User //Data verification
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	log := auditlog.New(auditlog.LogTypeUser, "用户登录", "", user)
	auditlog.Add(log)

	token, departName, departId, userType, roleId, err := userservice.Login(user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "userType": userType, "roleId": roleId}, "登陆成功!")
}

// 退出
func Logout(c *gin.Context) {

	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "用户注销", "", user)
	auditlog.Add(log)
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "退出成功!")

}

func Info(c *gin.Context) {
	user, _ := c.Get("x-user")
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{"user": dao.ToUserDto(user.(userservice.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	users, total, err := userservice.UserAll()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	data, err := common.DataPaging(query, users, total)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, int64(total), query)
}

// 高级搜索
func UserSearchHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	var email = user.Email
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	data, total, err := userservice.UserSearch(email, query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, int64(total), query)
}

// 修改密码
func UpdatePasswordHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := userservice.UpdatePassword(user.Email, user.Password)
	if err != nil {
		response.Fail(c, nil, err.Error())
	} else {
		response.Success(c, gin.H{"data": u}, "密码修改成功!")
	}
}

// 重置密码
func ResetPasswordHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := userservice.ResetPassword(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
	} else {
		response.Success(c, gin.H{"data": u}, "密码重置成功!")
	}
}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	var userdel userservice.Userdel
	if c.Bind(&userdel) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "删除用户", "", user)
	auditlog.Add(log)

	err := userservice.DeleteUser(userdel.Emails)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "用户删除成功!")
}

// 修改用户信息
func UpdateUserHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	log := auditlog.New(auditlog.LogTypeUser, "修改用户信息", "", user)
	auditlog.Add(log)

	u, err := userservice.UpdateUser(user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
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

	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "批量导入用户", "", user)
	auditlog.Add(log)

	var err error
	for _, file := range files {
		name := file.Filename
		c.SaveUploadedFile(file, name)
		xlFile, error := xlsx.OpenFile(name)
		if error != nil {
			return
		}
		UserExit, err = userservice.ReadFile(xlFile, UserExit)
		if err != nil {
			return
		}
	}

	if len(UserExit) == 0 {
		auditlog.UpdateStatus(log, auditlog.StatusSuccess)
		response.Success(c, nil, "导入用户信息成功")
	} else {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, gin.H{"UserExit": UserExit}, "以上用户已经存在")
	}
}
