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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func PolicyDelete(c *gin.Context) {
	var Rule model.CasbinRule
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			421,
			nil,
			err.Error())
		return
	}
	err = json.Unmarshal(body, &Rule)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	role_type := Rule.RoleType
	url := Rule.Url
	method := Rule.Method
	if ok := service.E.RemovePolicy(role_type, url, method); !ok {
		response.Response(c, http.StatusOK, 400, nil, "Pilocy不存在")
	} else {
		response.Success(c, gin.H{"code": 200}, "Pilocy删除成功")
	}
}

func PolicyAdd(c *gin.Context) {
	var Rule model.CasbinRule
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			421,
			nil,
			err.Error())
		return
	}
	err = json.Unmarshal(body, &Rule)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	role_type := Rule.RoleType
	url := Rule.Url
	method := Rule.Method

	if ok := service.E.AddPolicy(role_type, url, method); !ok {
		response.Response(c, http.StatusOK, 400, nil, "Pilocy已存在")
	} else {
		response.Success(c, gin.H{"code": 200}, "Pilocy添加成功")
	}
}

func GetPolicy(c *gin.Context) {
	casbin := make([]map[string]interface{}, 0)

	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	list := service.E.GetPolicy()
	for _, vlist := range list {
		policy := make(map[string]interface{})
		policy["role"] = vlist[0]
		policy["url"] = vlist[1]
		policy["method"] = vlist[2]
		casbin = append(casbin, policy)
	}
	total, data, err := model.SearchAll(query, casbin)
	if model.HandleError(c, err) {
		return
	}
	model.JsonPagination(c, data, total, query)
}

// 获取登录用户权限
type RoleID struct {
	RoleId []int `json:"roleId"`
}

func GetLoginUserPermission(c *gin.Context) {
	var role model.UserRole

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var RoleId RoleID
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &RoleId)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	min := RoleId.RoleId[0]
	if len(RoleId.RoleId) > 1 {
		for _, v := range RoleId.RoleId {
			if v < min {
				min = v
			}
		}
	}
	mysqlmanager.DB.Where("id=?", min).Find(&role)
	ID := role.ButtonID
	IDs := strings.Split(ID, ",")
	var buttons []string
	for _, id := range IDs {
		var SubButton model.RoleButton
		i, err := strconv.Atoi(id)
		if err != nil {
			response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
			return
		}
		mysqlmanager.DB.Where("id = ?", i).Find(&SubButton)
		button := SubButton.Button
		buttons = append(buttons, button)
	}
	response.Response(c, http.StatusOK, 200, gin.H{"userType": role.Type, "menu": role.Menus, "button": buttons}, "用户权限列表")
}

func GetRoles(c *gin.Context) {
	var roles []model.UserRole
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	mysqlmanager.DB.Find(&roles)
	datas := make([]map[string]interface{}, 0)

	for _, role := range roles {
		data := make(map[string]interface{})
		data["id"] = role.ID
		data["role"] = role.Role
		data["type"] = role.Type
		data["description"] = role.Description
		data["menus"] = role.Menus
		buttons := role.ButtonID
		if len(buttons) == 0 {
			data["buttons"] = []string{}
			datas = append(datas, data)
			continue
		}
		buttonss := strings.Split(buttons, ",")
		buts := make([]string, 0)
		for _, button := range buttonss {
			var but model.RoleButton
			i, err := strconv.Atoi(button)
			if err != nil {
				response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
				return
			}
			mysqlmanager.DB.Where("id=?", i).Find(&but)
			buts = append(buts, but.Button)
		}
		data["buttons"] = buts
		datas = append(datas, data)
	}
	total, data, err := model.SearchAll(query, datas)
	if model.HandleError(c, err) {
		return
	}
	model.JsonPagination(c, data, total, query)
}

func AddUserType(c *gin.Context) {
	var userRole model.UserRole
	c.Bind(&userRole)
	role := userRole.Role
	if len(role) == 0 {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"用户角色不能为空")
		return
	}
	user_type := userRole.Type
	description := userRole.Description
	userRole = model.UserRole{ //Create user
		Role:        role,
		Type:        user_type,
		Description: description,
	}
	mysqlmanager.DB.Save(&userRole)

	response.Success(c, nil, "新增角色成功")
}

type Roledel struct {
	ID int `json:"id"`
}

func DeleteUserRole(c *gin.Context) {
	var roledel Roledel
	var UserRole model.UserRole
	var users []model.User

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}

	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &roledel)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	mysqlmanager.DB.Find(&users)
	var contain []string
	for _, user := range users {
		id := user.RoleID
		if find := strings.Contains(id, strconv.Itoa(roledel.ID)); find {
			contain = append(contain, user.Email)
			break
		}
	}

	if len(contain) == 0 {
		mysqlmanager.DB.Where("id = ?", roledel.ID).Unscoped().Delete(UserRole)
		response.Response(c, http.StatusOK,
			422,
			nil,
			"角色删除成功")
	} else {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"有用户绑定此角色，不可删除")
	}
}

func UpdateUserRole(c *gin.Context) {
	var UserRole model.UserRole
	c.Bind(&UserRole)
	id := UserRole.ID
	role := UserRole.Role
	description := UserRole.Description
	mysqlmanager.DB.Where("id = ?", id).Find(&UserRole)
	if UserRole.Role != role && UserRole.Description != description {
		r := model.UserRole{
			Role:        role,
			Description: description,
		}
		mysqlmanager.DB.Model(&UserRole).Where("id = ?", id).Update(&r)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": UserRole},
			"角色信息修改成功")
		return
	}
	if UserRole.Role == role && UserRole.Description != description {
		r := model.UserRole{
			Description: description,
		}
		mysqlmanager.DB.Model(&UserRole).Where("id = ?", id).Update(&r)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": UserRole},
			"角色信息修改成功")
		return
	}
	if UserRole.Role != role && UserRole.Description == description {
		r := model.UserRole{
			Role: role,
		}
		mysqlmanager.DB.Model(&UserRole).Where("id = ?", id).Update(&r)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": UserRole},
			"角色信息修改成功")
		return
	}
}

type RoleChange struct {
	RoleID   int      `json:"id"`
	Menus    []string `json:"menus"`
	ButtonId []string `json:"buttonId"`
}

func RolePermissionChange(c *gin.Context) {
	var userRole model.UserRole
	var roleChange RoleChange

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}

	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &roleChange)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422333,
			nil,
			err.Error())
		return
	}
	// 数组切片转为string
	menus := strings.Replace(strings.Trim(fmt.Sprint(roleChange.Menus), "[]"), " ", ",", -1)
	buttonId := strings.Replace(strings.Trim(fmt.Sprint(roleChange.ButtonId), "[]"), " ", ",", -1)

	r := model.UserRole{
		Menus:    menus,
		ButtonID: buttonId,
	}
	mysqlmanager.DB.Model(&userRole).Where("id = ?", roleChange.RoleID).Update(&r)
	response.Response(c, http.StatusOK,
		200,
		gin.H{"data": userRole},
		"角色权限变更成功")
}
