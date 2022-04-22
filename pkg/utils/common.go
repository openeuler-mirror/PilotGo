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
 * Date: 2022-01-24 15:08:08
 * LastEditTime: 2022-03-17 14:55:44
 * Description: 公共函数
 ******************************************************************************/
package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

//分页查询结构体
type PaginationQ struct {
	Ok             bool        `json:"ok"`
	Size           uint        `form:"size" json:"size"`
	CurrentPageNum uint        `form:"page" json:"page"`
	Data           interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"` // save pagination list
	TotalPage      uint        `json:"total"`
}

// gorm分页查询方法
func CrudAll(p *PaginationQ, queryTx *gorm.DB, list interface{}) (uint, error) {
	if p.Size < 1 {
		p.Size = 10
	}
	if p.CurrentPageNum < 1 {
		p.CurrentPageNum = 1
	}

	var total uint
	err := queryTx.Count(&total).Error
	if err != nil {
		return 0, err
	}
	offset := p.Size * (p.CurrentPageNum - 1)
	err = queryTx.Limit(p.Size).Offset(offset).Find(list).Error
	if err != nil {
		return 0, err
	}
	return total, err
}

// map分页查询方法
func SearchAll(p *PaginationQ, data []map[string]interface{}) (uint, []map[string]interface{}, error) {
	if p.Size < 1 {
		p.Size = 10
	}
	if p.CurrentPageNum < 1 {
		p.CurrentPageNum = 1
	}
	total := len(data)
	if total == 0 {
		p.TotalPage = 1
	}
	num := p.Size * (p.CurrentPageNum - 1)
	if num > uint(total) {
		return uint(total), nil, fmt.Errorf("页码超出")
	}
	if p.Size*p.CurrentPageNum > uint(total) {
		return uint(total), data[num:], nil
	} else {
		if p.Size*p.CurrentPageNum < num {
			return uint(total), nil, fmt.Errorf("读取错误")
		}
		if p.Size*p.CurrentPageNum == 0 {
			return uint(total), data, nil
		} else {
			return uint(total), data[num : p.CurrentPageNum*p.Size-1], nil
		}
	}
}

// []map倒序函数
func Reverse(arr *[]map[string]interface{}) {
	var temp map[string]interface{}
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

// 拼装json 分页数据
func JsonPagination(c *gin.Context, list interface{}, total uint, query *PaginationQ) {
	c.AbortWithStatusJSON(200, gin.H{
		"code":  200,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.CurrentPageNum,
		"size":  query.Size})
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return true
	}
	return false
}
