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
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service/middleware"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// 删除过滤策略
func PolicyDelete(c *gin.Context) {
	var Rule model.CasbinRule
	c.Bind(&Rule)

	if ok := middleware.PolicyRemove(Rule); !ok {
		response.Response(c, http.StatusOK, http.StatusBadRequest, nil, "Pilocy不存在")
	} else {
		response.Success(c, gin.H{"code": http.StatusOK}, "Pilocy删除成功")
	}
}

// 增加过滤策略
func PolicyAdd(c *gin.Context) {
	var Rule model.CasbinRule
	c.Bind(&Rule)

	if ok := middleware.PolicyAdd(Rule); !ok {
		response.Response(c, http.StatusOK, http.StatusBadRequest, nil, "Pilocy已存在")
	} else {
		response.Success(c, gin.H{"code": http.StatusOK}, "Pilocy添加成功")
	}
}

// 获取所有过滤策略
func GetPolicy(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	policy, total := middleware.AllPolicy()

	data, err := DataPaging(query, policy, total)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	JsonPagination(c, data, total, query)
}

// 获取登录用户权限
func GetLoginUserPermission(c *gin.Context) {
	var RoleId model.RoleID
	c.Bind(&RoleId)

	roleId := service.RoleId(RoleId) //用户的最高权限

	userRole := dao.RoleIdToGetAllInfo(roleId)
	buttons := dao.PermissionButtons(userRole.ButtonID)

	response.Response(c, http.StatusOK, http.StatusOK, gin.H{"userType": userRole.Type, "menu": userRole.Menus, "button": buttons}, "用户权限列表")
}

func GetRoles(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	roles, total := dao.GetAllRoles()

	data, err := DataPaging(query, roles, total)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	JsonPagination(c, data, total, query)
}

func AddUserRole(c *gin.Context) {
	var userRole model.UserRole
	c.Bind(&userRole)

	err := dao.AddRole(userRole)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "角色添加失败")
	}
	response.Success(c, nil, "新增角色成功")
}

func DeleteUserRole(c *gin.Context) {
	var UserRole model.UserRole
	c.Bind(&UserRole)

	if ok := dao.IsUserBindingRole(UserRole.ID); !ok {
		dao.DeleteRole(UserRole.ID)
		response.Success(c, nil, "角色删除成功")
	} else {
		response.Fail(c, nil, "有用户绑定此角色，不可删除")
	}
}

func UpdateUserRole(c *gin.Context) {
	var UserRole model.UserRole
	c.Bind(&UserRole)
	id := UserRole.ID
	role := UserRole.Role
	description := UserRole.Description

	userRole := dao.RoleIdToGetAllInfo(id)
	if userRole.Role != role && userRole.Description != description {
		dao.UpdateRoleName(id, role)
		dao.UpdateRoleDescription(id, description)
		response.Success(c, gin.H{"data": UserRole}, "角色信息修改成功")
		return
	}
	if userRole.Role == role && userRole.Description != description {
		dao.UpdateRoleDescription(id, description)
		response.Success(c, gin.H{"data": UserRole}, "角色信息修改成功")
		return
	}
	if userRole.Role != role && userRole.Description == description {
		dao.UpdateRoleName(id, role)
		response.Success(c, gin.H{"data": UserRole}, "角色信息修改成功")
		return
	}
}

func RolePermissionChange(c *gin.Context) {
	var roleChange model.RolePermissionChange
	c.Bind(&roleChange)

	userRole := dao.UpdateRolePermission(roleChange)
	response.Success(c, gin.H{"data": userRole}, "角色权限变更成功")
}
