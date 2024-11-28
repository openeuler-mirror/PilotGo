/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package middleware

import (
	"net/http"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

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
