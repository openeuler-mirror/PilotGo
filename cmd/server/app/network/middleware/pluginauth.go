/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Aug 07 16:18:53 2025 +0800
 */
package middleware

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/sdk/plugin/jwt"
	"github.com/gin-gonic/gin"
)

// 检查plugin接口调用权限
func PluginAuthServiceCheck(c *gin.Context) {
	_, err := jwt.ParsePluginServiceClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "plugin token check error:" + err.Error()})
		c.Abort()
		return
	}

	c.Next()
}
