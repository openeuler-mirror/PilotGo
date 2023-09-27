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
 * Date: 2022-03-07 15:32:38
 * LastEditTime: 2022-04-12 14:10:09
 * Description: 权限控制
 ******************************************************************************/
package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auditlog"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/jwt"
	roleservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/role"
	"openeuler.org/PilotGo/PilotGo/sdk/response"
)

// 获取登录用户权限
func GetLoginUserPermissionHandler(c *gin.Context) {
	var RoleId roleservice.RoleID
	if err := c.Bind(&RoleId); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	menu, buttons, err := roleservice.GetLoginUserPermission(RoleId)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"menu": menu, "button": buttons}, "用户权限列表")
}

func GetRolesHandler(c *gin.Context) {
	data, err := roleservice.GetRoles()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	response.Success(c, data, "角色权限列表")
}

func AddRoleHandler(c *gin.Context) {
	params := &struct {
		Role        string `json:"role"`
		Description string `json:"description"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypePermission, "添加角色", "", user)
	auditlog.Add(log)

	userRole := &roleservice.UserRole{
		Role:        params.Role,
		Description: params.Description,
	}

	err = roleservice.AddRole(userRole)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, params.Role, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, gin.H{"error": err.Error()}, "角色添加失败")
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "新增角色成功", log.Module, params.Role, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, nil, "新增角色成功")
}

func DeleteRoleHandler(c *gin.Context) {
	params := &struct {
		Role string `json:"role"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypePermission, "删除角色", "", user)
	auditlog.Add(log)

	err = roleservice.DeleteRole(params.Role)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, params.Role, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, "有用户绑定此角色，不可删除")
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, params.Role, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, nil, "角色删除成功")
}

func UpdateRoleInfoHandler(c *gin.Context) {
	params := &struct {
		Role        string `json:"role"`
		Description string `json:"description"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypePermission, "修改角色信息", "", user)
	auditlog.Add(log)

	err = roleservice.UpdateRoleInfo(params.Role, params.Description)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, params.Role, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, params.Role, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
	response.Success(c, nil, "角色信息修改成功")
}

func RolePermissionChangeHandler(c *gin.Context) {
	params := &struct {
		ButtonId []string `json:"buttonId"`
		Menus    []string `json:"menus"`
		Role     string   `json:"role"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypePermission, "修改角色权限", "", user)
	auditlog.Add(log)

	err = roleservice.UpdateRolePermissions(params.Role, params.ButtonId, params.Menus)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, params.Role, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, params.Role, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)

	response.Success(c, nil, "角色权限变更成功")
}
