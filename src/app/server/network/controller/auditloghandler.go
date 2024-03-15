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
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 查询所有审计日志
func AuditLogAllHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := auditlog.GetAuditLogPaged(num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 根据模块名字查询日志
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

// 查询所有父日志为空的日志
func LogAllHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := auditlog.GetParentLog(num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 查询其子日志
func GetAuditLogByIdHandler(c *gin.Context) {
	parent_uuid := c.Query("uuid")
	data, err := auditlog.GetAuditLogById(parent_uuid)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, data, "审计子日志查询成功!")
}
