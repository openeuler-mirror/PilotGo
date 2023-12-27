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

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	userservice "gitee.com/openeuler/PilotGo/app/server/service/user"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/utils/message/net"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
)

// 添加用户
func RegisterHandler(c *gin.Context) {
	user := &userservice.UserInfo{}
	if err := c.Bind(user); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "添加用户",
	}
	auditlog.Add(log)

	err = userservice.Register(user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, nil, "添加用户成功!") //Return result
}

// 登录
func LoginHandler(c *gin.Context) {
	user := userservice.UserInfo{}
	if err := c.Bind(&user); err != nil {
		response.Fail(c, nil, net.GetValidMsg(err, &user))
		return
	}

	u, err := userservice.GetUserByEmail(user.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "用户登录",
	}
	auditlog.Add(log)

	departName, departId, roleId, err := userservice.Login(&user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}

	token, err := jwt.ReleaseToken(*u)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "roleId": roleId}, "登录成功!")
}

// 退出
func Logout(c *gin.Context) {
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "用户退出",
	}
	auditlog.Add(log)
	response.Success(c, nil, "退出成功!")
}

func Info(c *gin.Context) {
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	d := &struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}{
		Name:  u.Username,
		Phone: u.Phone,
	}

	response.Success(c, d, "用户信息查询成功")
}

// 查询所有用户
func UserAll(c *gin.Context) {
	p := &common.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := p.Size * (p.CurrentPageNum - 1)
	total, data, err := userservice.GetUserPaged(num, p.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, p)
}

// 高级搜索
func UserSearchHandler(c *gin.Context) {
	user := &struct {
		Email string `json:"email"`
	}{}
	if err := c.Bind(&user); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := userservice.UserSearchPaged(user.Email, num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 修改密码
func UpdatePasswordHandler(c *gin.Context) {
	var user userservice.UserInfo
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "修改密码",
	}
	auditlog.Add(log)

	err = userservice.UpdatePassword(user.Email, user.Password)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "密码修改成功!")
}

// 重置密码
func ResetPasswordHandler(c *gin.Context) {
	var user userservice.UserInfo
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "重置密码",
	}
	auditlog.Add(log)

	err = userservice.ResetPassword(user.Email)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "密码重置成功!")
}

// 删除用户
func DeleteUserHandler(c *gin.Context) {
	statuscodes := []string{}
	fd := struct {
		Deldatas []string `json:"delDatas,omitempty"`
	}{}
	if err := c.Bind(&fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "删除用户",
	}
	auditlog.Add(log)

	for _, ps := range fd.Deldatas {
		log_s := &auditlog.AuditLog{
			LogUUID:    uuid.New().String(),
			ParentUUID: log.LogUUID,
			Module:     auditlog.ModuleUser,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "删除用户",
			Message:    "用户：" + strings.Split(ps, "/")[0],
		}
		auditlog.Add(log_s)

		err := userservice.DeleteUser(strings.Split(ps, "/")[0])
		if err != nil {
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		statuscodes = append(statuscodes, strconv.Itoa(http.StatusOK))
	}

	status := auditlog.BatchActionStatus(statuscodes)
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
	user := userservice.UserInfo{}
	if c.Bind(&user) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "修改用户信息",
	}
	auditlog.Add(log)

	err = userservice.UpdateUser(user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "用户信息修改成功")
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleUser,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "批量导入用户",
	}
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
		response.Success(c, nil, "导入用户信息成功")
	} else {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"UserExit": UserExit}, "以上用户已经存在")
	}
}
