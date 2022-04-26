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
package controller

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// gorm分页查询方法
func CrudAll(p *model.PaginationQ, queryTx *gorm.DB, list interface{}) (int, error) {
	if p.Size < 1 {
		p.Size = 10
	}
	if p.CurrentPageNum < 1 {
		p.CurrentPageNum = 1
	}

	var total int
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

// 结构体分页查询方法
func DataPaging(p *model.PaginationQ, list interface{}, total int) (interface{}, error) {
	data := make([]interface{}, 0)
	if reflect.TypeOf(list).Kind() == reflect.Slice {
		s := reflect.ValueOf(list)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			data = append(data, ele.Interface())
		}
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.CurrentPageNum < 1 {
		p.CurrentPageNum = 1
	}
	if total == 0 {
		p.TotalPage = 1
	}
	num := p.Size * (p.CurrentPageNum - 1)
	if num > uint(total) {
		return nil, fmt.Errorf("页码超出")
	}
	if p.Size*p.CurrentPageNum > uint(total) {
		return data[num:], nil
	} else {
		if p.Size*p.CurrentPageNum < num {
			return nil, fmt.Errorf("读取错误")
		}
		if p.Size*p.CurrentPageNum == 0 {
			return data, nil
		} else {
			return data[num : p.CurrentPageNum*p.Size], nil
		}
	}
}

// 拼装json 分页数据
func JsonPagination(c *gin.Context, list interface{}, total int, query *model.PaginationQ) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.CurrentPageNum,
		"size":  query.Size})
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return true
	}
	return false
}
