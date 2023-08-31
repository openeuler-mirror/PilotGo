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

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auditlog"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auth"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	roleservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/role"
	userservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/user"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 删除过滤策略
func PolicyDelete(c *gin.Context) {
	var Rule auth.CasbinRule
	if err := c.Bind(&Rule); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if ok := auth.PolicyRemove(Rule); !ok {
		response.Fail(c, nil, "Pilocy不存在")
	} else {
		response.Success(c, gin.H{"code": http.StatusOK}, "Pilocy删除成功")
	}
}

// 增加过滤策略
func PolicyAdd(c *gin.Context) {
	var Rule auth.CasbinRule
	if err := c.Bind(&Rule); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if ok := auth.PolicyAdd(Rule); !ok {
		response.Fail(c, nil, "Pilocy已存在")
	} else {
		response.Success(c, gin.H{"code": http.StatusOK}, "Pilocy添加成功")
	}
}

// 获取所有过滤策略
func GetPolicy(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	policy, total := auth.AllPolicy()

	data, err := common.DataPaging(query, policy, total)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	common.JsonPagination(c, data, int64(total), query)
}

// 获取登录用户权限
func GetLoginUserPermissionHandler(c *gin.Context) {
	var RoleId roleservice.RoleID
	if err := c.Bind(&RoleId); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	userRole, buttons, err := roleservice.GetLoginUserPermission(RoleId)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"userType": userRole.Type, "menu": userRole.Menus, "button": buttons}, "用户权限列表")
}

func GetRolesHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, data, err := roleservice.GetRoles(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, int64(total), query)
}

func AddUserRoleHandler(c *gin.Context) {
	var userRole roleservice.UserRole
	if err := c.Bind(&userRole); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user *userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "添加角色", "", user)
	auditlog.Add(log)

	err := userservice.AddUserRole(&userRole)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, gin.H{"error": err.Error()}, "角色添加失败")
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "新增角色成功")
}

func DeleteUserRoleHandler(c *gin.Context) {
	var UserRole roleservice.UserRole
	if err := c.Bind(&UserRole); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user *userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "删除角色", "", user)
	auditlog.Add(log)

	err := userservice.DeleteUserRole(UserRole.ID)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, "有用户绑定此角色，不可删除")
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "角色删除成功")
}

func UpdateUserRoleHandler(c *gin.Context) {
	var UserRole roleservice.UserRole
	if err := c.Bind(&UserRole); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user *userservice.User
	log := auditlog.New(auditlog.LogTypeUser, "修改角色", "", user)
	auditlog.Add(log)

	err := userservice.UpdateUserRole(&UserRole)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"data": UserRole}, "角色信息修改成功")
}

func RolePermissionChangeHandler(c *gin.Context) {
	var roleChange roleservice.RolePermissionChange
	if err := c.Bind(&roleChange); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	//TODO:
	var user *userservice.User
	log := auditlog.New(auditlog.LogTypePermission, "修改角色权限", "", user)
	auditlog.Add(log)

	userRole, err := roleservice.RolePermissionChangeMethod(roleChange)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, gin.H{"data": userRole}, "角色权限变更成功")
}
