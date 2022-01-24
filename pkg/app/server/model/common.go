package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
)

/**
 * @Author: zhang han
 * @Date: 2021/11/12 17:01
 * @Description:定义分页结构体及公共的结构体
 */

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

// 分页查询方法
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
		response.Fail(c, gin.H{"ok": false}, "获取失败!")
		return true
	}
	return false
}
