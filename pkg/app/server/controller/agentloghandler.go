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
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// 查询所有父日志
func LogAll(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	Id := c.Query("departId")
	departId, err := strconv.Atoi(Id)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	departIds := make([]int, 0)
	departIds = append(departIds, departId)
	ReturnSpecifiedDepart(departId, &departIds)

	// 获取部门名字
	departNames := make([]string, 0)
	for _, id := range departIds {
		departName := dao.DepartIdToGetDepartName(id)
		departNames = append(departNames, departName)
	}

	logParent := model.AgentLogParent{}
	list, total := logParent.LogAll(query, departNames)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	data, err := DataPaging(query, list, total)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	JsonPagination(c, data, total, query)
}

// 查询所有子日志
func AgentLogs(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	ParentId := c.Query("id")
	parentId, err := strconv.Atoi(ParentId)
	if err != nil {
		response.Fail(c, nil, "父日志ID输入格式有误")
		return
	}

	logs := model.AgentLog{}
	list, tx := logs.AgentLog(query, parentId)

	total, err := CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	JsonPagination(c, list, total, query)
}

// 删除机器日志
func DeleteLog(c *gin.Context) {
	var logid model.AgentLogDel
	c.Bind(&logid)

	dao.LogDelete(logid.IDs)
	response.Success(c, nil, "日志删除成功!")
}
