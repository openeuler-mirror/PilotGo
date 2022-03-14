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
 * LastEditTime: 2022-03-14 14:51:15
 * Description: 权限控制
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
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
	casbin := []map[string]string{}

	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	list := common.E.GetPolicy()
	for _, vlist := range list {
		policy := map[string]string{}
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
