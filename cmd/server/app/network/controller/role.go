/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	roleservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/role"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
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
	p := &response.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := p.PageSize * (p.Page - 1)
	total, data, err := roleservice.GetRolePaged(num, p.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), p)
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

	userRole := &roleservice.Role{
		Name:        params.Role,
		Description: params.Description,
	}

	err := roleservice.AddRole(userRole)
	if err != nil {
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

	err := roleservice.DeleteRole(params.RoleId)
	if err != nil {
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

	err := roleservice.UpdateRoleInfo(params.Role, params.Description)
	if err != nil {
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

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "角色权限变更",
		Module:     auditlog.RoleChange,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "角色权限变更：" + params.Role,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	err = roleservice.UpdateRolePermissions(params.Role, params.Buttons, params.Menus, params.PluginPermissions)
	if err != nil {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
		auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "角色权限变更失败："+err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")
	response.Success(c, nil, "角色权限变更成功")
}
