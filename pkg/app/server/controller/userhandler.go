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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
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
	var user *userservice.User
	if err := c.Bind(&user); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	log := auditlog.New(auditlog.LogTypeUser, "添加用户", "", user)
	auditlog.Add(log)

	err := userservice.Register(user)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, nil, "添加用户成功!") //Return result
}

func LoginHandler(c *gin.Context) {
	var user *userservice.User //Data verification
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := dao.UserInfo(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	user.DepartName = u.DepartName

	log := auditlog.New(auditlog.LogTypeUser, "用户登录", "", user)
	auditlog.Add(log)

	token, departName, departId, userType, roleId, err := userservice.Login(*user)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "userType": userType, "roleId": roleId}, "登陆成功!")
}

// 退出
func Logout(c *gin.Context) {
	frontdata := &userservice.Frontdata{}
	var user userservice.User
	if err := c.Bind(frontdata); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	user.Username = frontdata.Username
	user.DepartName = frontdata.Departname
	log := auditlog.New(auditlog.LogTypeUser, "用户注销", "", &user)
	auditlog.Add(log)
	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
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
	var user *userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := dao.UserInfo(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	user.DepartName = u.DepartName

	log := auditlog.New(auditlog.LogTypeUser, "修改密码", "", user)
	auditlog.Add(log)

	u, err = userservice.UpdatePassword(user.Email, user.Password)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, gin.H{"data": u}, "密码修改成功!")
}

// 重置密码
func ResetPasswordHandler(c *gin.Context) {
	var user *userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := dao.UserInfo(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	user.DepartName = u.DepartName

	log := auditlog.New(auditlog.LogTypeUser, "重置密码", "", user)
	auditlog.Add(log)

	u, err = userservice.ResetPassword(user.Email)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, gin.H{"data": u}, "密码重置成功!")

}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	user := &dao.User{}
	fail_m := map[string]string{}
	statuscodes := []string{}
	userdel := &userservice.Frontdata{}

	if err := c.Bind(&userdel); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	user.Email = userdel.Username
	user.DepartName = userdel.Departname
	log := auditlog.New(auditlog.LogTypeUser, "删除用户", "", user)
	auditlog.Add(log)

	for _, ps := range userdel.Deldatas {
		err := userservice.DeleteUser(strings.Split(ps, "/")[0])
		if err != nil {
			log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
			auditlog.Add(log_s)
			fail_m[user.Email] = err.Error()
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionOK)
	}

	status := service.BatchActionStatus(statuscodes)
	auditlog.UpdateStatus(log, status)

	if len(fail_m) != 0 {
		fail_sl := []string{}
		for k, v := range fail_m {
			fail_sl = append(fail_sl, k+":"+v)
		}
		response.Fail(c, nil, strings.Join(fail_sl, " "))
		return
	}

	response.Success(c, nil, "用户删除成功!")
}

// 修改用户信息
func UpdateUserHandler(c *gin.Context) {
	var user *userservice.User
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	log := auditlog.New(auditlog.LogTypeUser, "修改用户信息", "", user)
	auditlog.Add(log)

	u, err := userservice.UpdateUser(*user)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusFail, log.Module, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, auditlog.StatusSuccess, log.Module, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
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
	log := auditlog.New(auditlog.LogTypeUser, "批量导入用户", "", &user)
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
