/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaginationQ struct {
	Ok             bool        `json:"ok"`
	Size           int         `form:"size" json:"size"`
	CurrentPageNum int         `form:"page" json:"page"`
	Data           interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"` // save pagination list
	TotalPage      int         `json:"total"`
}

// 拼装json 分页数据
func JsonPagination(c *gin.Context, list interface{}, total int64, query *PaginationQ) {
	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.CurrentPageNum,
		"size":  query.Size})
}
