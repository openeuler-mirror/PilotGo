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
 * LastEditTime: 2022-04-22 14:18:14
 * Description: 公共函数
 ******************************************************************************/
package common

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

type PaginationQ struct {
	Ok             bool        `json:"ok"`
	Size           int         `form:"size" json:"size"`
	CurrentPageNum int         `form:"page" json:"page"`
	Data           interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"` // save pagination list
	TotalPage      int         `json:"total"`
}

// 返回所有子部门id
func ReturnSpecifiedDepart(id int, res *[]int) {
	temp, err := dao.SubDepartId(id)
	if err != nil {
		logger.Error(err.Error())
	}
	if len(temp) == 0 {
		return
	}
	for _, value := range temp {
		*res = append(*res, value)
		ReturnSpecifiedDepart(value, res)
	}
}

// 拼装json 分页数据
func JsonPagination(c *gin.Context, list interface{}, total int64, query *PaginationQ) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.CurrentPageNum,
		"size":  query.Size})
}
