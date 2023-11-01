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
	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	roleservice "gitee.com/openeuler/PilotGo/app/server/service/role"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	userRole := &roleservice.UserRole{
		Role:        params.Role,
		Description: params.Description,
	}

	err = roleservice.AddRole(userRole)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"error": err.Error()}, "角色添加失败")
		return
	}
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

	err = roleservice.DeleteRole(params.Role)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "有用户绑定此角色，不可删除")
		return
	}
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

	err = roleservice.UpdateRolePermissions(params.Role, params.ButtonId, params.Menus)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "角色权限变更成功")
}
