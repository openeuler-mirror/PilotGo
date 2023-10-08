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
 * Date: 2022-07-11 15:25:53
 * LastEditTime: 2022-07-11 10:35:54
 * Description: panic recover
 ******************************************************************************/
package middleware

import (
	"net/http"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func LoggerDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Debug("status_code:%d latency_time:%s client_ip:%s req_method:%s req_uri:%s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 开始时间
			startTime := time.Now()

			// 处理请求
			c.Next()

			// 结束时间
			endTime := time.Now()

			// 执行时间
			latencyTime := endTime.Sub(startTime)

			// 请求方式
			reqMethod := c.Request.Method

			// 请求路由
			reqUri := c.Request.RequestURI

			// 状态码
			statusCode := c.Writer.Status()

			// 请求IP
			clientIP := c.ClientIP()

			// 日志格式
			logger.Error("status_code:%d latency_time:%s client_ip:%s req_method:%s req_uri:%s panic:%v",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
				r,
			)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}
