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
