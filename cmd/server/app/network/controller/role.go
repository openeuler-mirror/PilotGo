/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/common"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	roleservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/role"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 获取所有角色
func GetRolesHandler(c *gin.Context) {
	data, err := roleservice.GetRoles()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, data, "角色权限列表")
}

// 分页获取所有角色
func GetRolesPagedHandler(c *gin.Context) {
	p := &common.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := p.Size * (p.CurrentPageNum - 1)
	total, data, err := roleservice.GetRolePaged(num, p.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, p)
}

// 添加角色
func AddRoleHandler(c *gin.Context) {
	params := &struct {
		Role        string `json:"role"`
		Description string `json:"description"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
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
		Module:     auditlog.ModuleRole,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "添加角色",
	}
	auditlog.Add(log)

	userRole := &roleservice.Role{
		Name:        params.Role,
		Description: params.Description,
	}

	err = roleservice.AddRole(userRole)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "新增角色成功")
}

// 删除角色
func DeleteRoleHandler(c *gin.Context) {
	params := &struct {
		RoleId int `json:"role"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
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
		Module:     auditlog.ModuleRole,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "删除角色",
	}
	auditlog.Add(log)

	err = roleservice.DeleteRole(params.RoleId)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "有用户绑定此角色，不可删除")
		return
	}
	response.Success(c, nil, "角色删除成功")
}

// 更改角色
func UpdateRoleInfoHandler(c *gin.Context) {
	params := &struct {
		Role        string `json:"role"`
		Description string `json:"description"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
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
		Module:     auditlog.ModuleRole,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "修改角色信息",
	}
	auditlog.Add(log)

	err = roleservice.UpdateRoleInfo(params.Role, params.Description)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "角色信息修改成功")
}

// 更改角色权限
func RolePermissionChangeHandler(c *gin.Context) {
	params := &struct {
		Buttons           []string                  `json:"buttons"`
		Menus             []string                  `json:"menus"`
		Role              string                    `json:"role"`
		PluginPermissions []plugin.PluginPermission `json:"pluginpermissions"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error:"+err.Error())
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
		Module:     auditlog.ModuleRole,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "修改角色权限",
	}
	auditlog.Add(log)

	err = roleservice.UpdateRolePermissions(params.Role, params.Buttons, params.Menus, params.PluginPermissions)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "角色权限变更成功")
}
