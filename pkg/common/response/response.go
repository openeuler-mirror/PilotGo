package response

/**
 * @Author: zhang han
 * @Date: 2021/11/1 9:52
 * @Description: 封装response的返回参数
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg})
}

func Success(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 200, data, msg)
}

func Fail(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusBadRequest, 400, data, msg)
}
