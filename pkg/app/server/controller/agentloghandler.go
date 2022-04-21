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
 * LastEditTime: 2022-03-24 00:18:14
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
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// 查询所有父日志
func LogAll(c *gin.Context) {
	Id := c.Query("departId")
	departId, err := strconv.Atoi(Id)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	logParent := model.AgentLogParent{}
	query := &model.PaginationQ{}
	err = c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	Dids := make([]int, 0)
	Dnames := make([]string, 0)
	ReturnSpecifiedDepart(departId, &Dids)
	Dids = append(Dids, departId)
	for _, id := range Dids {
		var departNames model.DepartNode
		mysqlmanager.DB.Where("id =?", id).Find(&departNames)

		Dnames = append(Dnames, departNames.Depart)
	}
	list, total, err := logParent.LogAll(query, Dnames)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	JsonPagination(c, list, total, query)
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
	if HandleError(c, err) {
		return
	}

	list, tx := logs.AgentLog(query, parentId)

	total, err := CrudAll(query, tx, list)
	if HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	JsonPagination(c, list, total, query)
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
