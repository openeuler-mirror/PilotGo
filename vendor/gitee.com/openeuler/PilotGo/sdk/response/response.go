/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}, msg string) {
	result(c, http.StatusOK, http.StatusOK, data, msg)
}
func Fail(c *gin.Context, data interface{}, msg string) {
	result(c, http.StatusOK, http.StatusBadRequest, data, msg)
}

// 拼装json 分页数据
func DataPagination(c *gin.Context, list interface{}, total int, query *PaginationQ) {
	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.Page,
		"size":  query.PageSize})
}

func result(c *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg})
}

type PaginationQ struct {
	Ok        bool        `json:"ok"`
	PageSize  int         `form:"size" json:"size"`
	Page      int         `form:"page" json:"page"`
	Data      interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"`
	TotalSize int         `json:"total"`
}
