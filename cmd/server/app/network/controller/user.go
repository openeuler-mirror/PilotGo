/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	eventSDK "gitee.com/openeuler/PilotGo-plugins/event/sdk"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/role"
	userservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/user"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/net"
	commonSDK "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
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

	// email格式校验
	if user.Email != "admin" {
		patt := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
		if match, err := regexp.MatchString(patt, user.Email); err != nil || !match {
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("email format error: %s, %s", user.Email, err.Error()))
				return
			}
			response.Fail(c, nil, fmt.Sprintf("email format error: %s", user.Email))
			return
		}
	}
	// phone格式校验
	if user.Phone != "" {
		patt := `^[1]([3-9])[0-9]{9}$`
		if match, err := regexp.MatchString(patt, user.Phone); err != nil || !match {
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("phone format error: %s, %s", user.Phone, err.Error()))
				return
			}
			response.Fail(c, nil, fmt.Sprintf("phone format error: %s", user.Phone))
			return
		}
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

// @Summary 用户登录
// @Description 用户登录接口,返回用户信息和token
// @Tags user
// @Accept json
// @Produce json
// @Param user body userservice.UserInfo true "用户登录信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 200 {string} string "登录成功"
// @Router /api/v1/user/login [post]
func LoginHandler(c *gin.Context) {
	user := userservice.UserInfo{}
	if err := c.Bind(&user); err != nil {
		response.Fail(c, nil, net.GetValidMsg(err, &user))
		return
	}

	// 输入email格式校验
	if user.Email != "admin" {
		patt := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
		if match, err := regexp.MatchString(patt, user.Email); err != nil || !match {
			if err != nil {
				response.Fail(c, nil, fmt.Sprintf("email format error: %s, %s", user.Email, err.Error()))
				return
			}
			response.Fail(c, nil, fmt.Sprintf("email format error: %s", user.Email))
			return
		}
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
	// 发布“用户登录”事件
	msgData := commonSDK.MessageData{
		MsgType:     eventSDK.MsgUserLogin,
		MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgUserLogin),
		TimeStamp:   time.Now(),
		Data: eventSDK.MDUserSystemSession{
			UserName: u.Username,
			Email:    u.Email,
		},
	}
	msgDataString, err := msgData.ToMessageDataString()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	ms := commonSDK.EventMessage{
		MessageType: eventSDK.MsgUserLogin,
		MessageData: msgDataString,
	}
	plugin.PublishEvent(ms)

	departName, departId, roleId, err := userservice.Login(&user)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}

	token, err := jwt.GenerateUserToken(*u)
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
	// 发布“用户退出”事件
	msgData := commonSDK.MessageData{
		MsgType:     eventSDK.MsgUserLogout,
		MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgUserLogout),
		TimeStamp:   time.Now(),
		Data: eventSDK.MDUserSystemSession{
			UserName: u.Username,
			Email:    u.Email,
		},
	}
	msgDataString, err := msgData.ToMessageDataString()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	ms := commonSDK.EventMessage{
		MessageType: eventSDK.MsgUserLogout,
		MessageData: msgDataString,
	}
	plugin.PublishEvent(ms)

	response.Success(c, nil, "退出成功!")
}

func Info(c *gin.Context) {
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	d, err := userservice.GetUserByEmail(u.Email)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, d, "用户信息查询成功")
}

// 查询所有用户
func UserAll(c *gin.Context) {
	p := &response.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := p.PageSize * (p.Page - 1)
	total, data, err := userservice.GetUserPaged(num, p.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), p)
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
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	num := query.PageSize * (query.Page - 1)
	total, data, err := userservice.UserSearchPaged(user.Email, num, query.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), query)
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
		response.Fail(c, nil, "参数错误")
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
	response.Success(c, nil, "密码重置成功! 初始密码为邮箱用户名")
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
	_, file, err := c.Request.FormFile("upload")
	if err != nil {
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

	name := file.Filename
	c.SaveUploadedFile(file, name)
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		response.Fail(c, nil, "获取文件错误："+err.Error())
		return
	}
	UserExit, err = userservice.ReadFile(xlFile, UserExit)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		auditlog.UpdateMessage(log, strings.Join(UserExit, ";"))
		println(log)
		response.Fail(c, gin.H{"UserExit": UserExit}, err.Error())
		return
	}

	if len(UserExit) == 0 {
		response.Success(c, nil, "导入用户信息成功")
		return
	} else {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		auditlog.UpdateMessage(log, strings.Join(UserExit, ";"))
		response.Fail(c, gin.H{"UserExit": UserExit}, "以上用户已经存在")
	}
}

// 获取登录用户权限
func GetLoginUserPermissionHandler(c *gin.Context) {
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	roleids, err := userservice.GetRolesByUid(u.ID)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	permissions, err := role.GetLoginUserPermission(roleids)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, permissions, "用户权限列表")
}
