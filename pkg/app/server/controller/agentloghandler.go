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
 * Date: 2022-02-23 17:44:00
 * LastEditTime: 2022-03-18 17:22:53
 * Description: provide agent log manager functions.
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

// 查询所有父日志
func LogAll(c *gin.Context) {
	var logparents []model.AgentLogParent
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	mysqlmanager.DB.Find(&logparents)
	datas := make([]map[string]interface{}, 0)
	for _, logparent := range logparents {
		var user model.User
		data := make(map[string]interface{})
		data["id"] = logparent.ID
		data["created_at"] = logparent.CreatedAt
		data["userName"] = logparent.UserName
		mysqlmanager.DB.Where("email = ?", logparent.UserName).Find(&user)
		data["departName"] = user.DepartName
		data["type"] = logparent.Type
		data["status"] = logparent.Status
		datas = append(datas, data)
	}
	total, data, err := model.SearchAll(query, datas)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	model.JsonPagination(c, data, total, query)
}

// 查询所有子日志
func AgentLogs(c *gin.Context) {
	ParentId := c.Query("id")
	parentId, err := strconv.Atoi(ParentId)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"父日志ID输入格式有误")
		return
	}
	logs := model.AgentLog{}
	query := &model.PaginationQ{}
	err = c.ShouldBindQuery(query)
	if model.HandleError(c, err) {
		return
	}

	list, total, err := logs.AgentLog(query, parentId)
	if model.HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	model.JsonPagination(c, list, total, query)
}

// 删除机器日志
type AgentLogDel struct {
	IDs []int `json:"ids,omitempty"`
}

func DeleteLog(c *gin.Context) {
	var logparent model.AgentLogParent
	var log model.AgentLog
	var logid AgentLogDel

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &logid)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	for _, id := range logid.IDs {
		mysqlmanager.DB.Where("log_parent_id=?", id).Unscoped().Delete(log)
		mysqlmanager.DB.Where("id=?", id).Unscoped().Delete(logparent)
	}
	response.Response(c, http.StatusOK,
		200,
		nil,
		"日志删除成功!")
}
