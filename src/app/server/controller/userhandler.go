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
 * LastEditTime: 2023-09-04 14:06:46
 * Description: 用户登录、增删改查
 ******************************************************************************/
package controller

import (
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/app/server/service"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/app/server/service/jwt"
	userservice "gitee.com/openeuler/PilotGo/app/server/service/user"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func GetRoleListHandler(c *gin.Context) {
	roles, err := userservice.GetUserRole()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"role": roles}, "获取用户角色")
}

func RegisterHandler(c *gin.Context) {
	fd := &userservice.Frontdata{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "添加用户", user.ID)
	auditlog.Add(log)

	err = userservice.Register(fd)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "添加用户成功!") //Return result
}

func LoginHandler(c *gin.Context) {
	fd := &userservice.Frontdata{}
	if c.Bind(fd) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := dao.UserInfo(fd.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	fd.Username_creaate = u.Email
	fd.Departname_create = u.DepartName
	log := auditlog.New(auditlog.LogTypeUser, "用户登录", u.ID)
	auditlog.Add(log)

	token, departName, departId, roleId, err := userservice.Login(*fd)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "roleId": roleId}, "登录成功!")
}

// 退出
func Logout(c *gin.Context) {
	fd := &userservice.Frontdata{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "用户注销", u.ID)
	auditlog.Add(log)
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "退出成功!")

}

func Info(c *gin.Context) {
	user, _ := c.Get("x-user")
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		// TODO: fix type assertion
		"data": gin.H{"user": dao.ToUserDto(user.(common.User))},
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
	var user common.User
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	user.DepartName = u.DepartName
	log := auditlog.New(auditlog.LogTypeUser, "修改密码", u.ID)
	auditlog.Add(log)

	user, err = userservice.UpdatePassword(user.Email, user.Password)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"data": user}, "密码修改成功!")
}

// 重置密码
func ResetPasswordHandler(c *gin.Context) {
	var user userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "重置密码", u.ID)
	auditlog.Add(log)

	user, err = userservice.ResetPassword(user.Email)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"data": user}, "密码重置成功!")

}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	statuscodes := []string{}
	fd := &userservice.Frontdata{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "删除用户", u.ID)
	auditlog.Add(log)

	for _, ps := range fd.Deldatas {
		err := userservice.DeleteUser(strings.Split(ps, "/")[0])
		if err != nil {
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		statuscodes = append(statuscodes, strconv.Itoa(http.StatusOK))
	}

	status := service.BatchActionStatus(statuscodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error(err.Error())
	}

	switch strings.Split(status, ",")[2] {
	case "0.00":
		response.Fail(c, nil, "用户删除失败")
		return
	case "1.00":
		response.Success(c, nil, "用户删除成功")
	default:
		response.Success(c, nil, "用户删除部分成功")
	}
}

// 修改用户信息
func UpdateUserHandler(c *gin.Context) {
	fd := &userservice.Frontdata{}
	if c.Bind(fd) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "修改用户信息", u.ID)
	auditlog.Add(log)

	user, err := userservice.UpdateUser(*fd)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"data": user}, "用户信息修改成功")

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

	// TODO:
	fd := &userservice.Frontdata{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypeUser, "批量导入用户", u.ID)
	auditlog.Add(log)

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
