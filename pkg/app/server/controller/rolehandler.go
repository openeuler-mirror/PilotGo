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
 * LastEditTime: 2022-03-17 16:15:54
 * Description: 权限控制
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
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
	if ok := common.E.RemovePolicy(role_type, url, method); !ok {
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

	if ok := common.E.AddPolicy(role_type, url, method); !ok {
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
	list := common.E.GetPolicy()
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

// 给用户添加角色
func AddPermission(c *gin.Context) {
	var role model.AddRole //Data verification
	var user model.User
	c.Bind(&role)
	email := role.Email
	addrole := role.RoleID
	mysqlmanager.DB.Where("email=?", email).Find(&user)
	if !dao.IsEmailExist(email) {
		response.Fail(c, nil, "无此用户!")
		return
	}
	roleid := user.RoleID
	if dao.IsContain(roleid, addrole) {
		response.Response(c,
			http.StatusOK,
			400,
			nil,
			"用户已拥有该权限!")
		return
	} else {
		roleid = roleid + "," + strconv.Itoa(addrole)
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("roleId", roleid)

		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"用户角色添加成功!")
	}
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
