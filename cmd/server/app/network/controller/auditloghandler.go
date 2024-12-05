/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 查询所有审计日志
func AuditLogAllHandler(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.PageSize * (query.Page - 1)
	total, data, err := auditlog.GetAuditLogPaged(num, query.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), query)
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
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.PageSize * (query.Page - 1)
	total, data, err := auditlog.GetParentLog(num, query.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), query)
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
