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
