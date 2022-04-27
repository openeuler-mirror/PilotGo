/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhang&wang
 * Date: 2022-02-23 17:44:00
 * LastEditTime: 2022-04-22 14:18:14
 * Description: 公共函数
 ******************************************************************************/
package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// gorm分页查询方法
func CrudAll(p *model.PaginationQ, queryTx *gorm.DB, list interface{}) (uint, error) {
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
func SearchAll(p *model.PaginationQ, data []map[string]interface{}) (uint, []map[string]interface{}, error) {
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
func JsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQ) {
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
func Paging(len int, size int, page int, Res *[]interface{}) error {

	if len == 0 {
		*Res = nil
		return nil
	}
	num := size * (page - 1)
	if num > len {
		return fmt.Errorf("page overflow")
	}
	if page*size >= len {
		*Res = (*Res)[num:]
		return nil
	} else {
		if page*size < num {
			return fmt.Errorf("Read error")
		}
		if page*size == 0 {
			*Res = nil
			return nil
		} else {
			*Res = (*Res)[num : page*size-1]
			return nil
		}

	}

}
