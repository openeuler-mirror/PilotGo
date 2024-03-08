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
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auditlog"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 查询所有审计日志
func AuditLogAllHandler(c *gin.Context) {
	loglist, _, err := auditlog.Get()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, loglist, "审计日志查询成功!")
}

func ModuleLogHandler(c *gin.Context) {
	var moduleName string
	if c.Bind(&moduleName) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	loglist, err := auditlog.GetByModule(moduleName)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, loglist, "模块审计日志查询成功!")
}

// 查询所有父日志
func LogAllHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// logParent := dao.AgentLogParent{}
	// list, tx := logParent.LogAll()
	list, tx, err := auditlog.Get()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	common.JsonPagination(c, list, total, query)
}

// 查询所有子日志
func AgentLogsHandler(c *gin.Context) {
	ParentId := c.Query("id")
	// parentId, err := strconv.Atoi(ParentId)
	// if err != nil {
	// 	response.Fail(c, nil, "父日志ID输入格式有误")
	// 	return
	// }

	// agentlog, err := service.AgentLogs(parentId)
	agentlog, _, err := auditlog.GetByParentId(ParentId)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, "获取子日志失败: "+err.Error())
		return
	}

	response.Success(c, agentlog, "子日志查询成功!")
}
