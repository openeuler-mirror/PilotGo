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
 * LastEditTime: 2023-03-16 15:08:01
 * Description: 封装response的返回参数
 ******************************************************************************/
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func result(c *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg})
}

func Success(c *gin.Context, data interface{}, msg string) {
	result(c, http.StatusOK, http.StatusOK, data, msg)
}

func Fail(c *gin.Context, data interface{}, msg string) {
	result(c, http.StatusOK, http.StatusBadRequest, data, msg)
}
