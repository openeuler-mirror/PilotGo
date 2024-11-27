//go:build !production
// +build !production

/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

package resource

import (
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func StaticRouter(router *gin.Engine) {
	router.Static("/assets", "./frontend/dist/assets")
	router.Static("/pilotgo.ico", "./frontend/dist/pilotgo.ico")
	router.StaticFile("/", "./frontend/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Debug("process noroute: %s", c.Request.URL.String())
		if !strings.HasPrefix(c.Request.RequestURI, "/api/") && !strings.HasPrefix(c.Request.RequestURI, "/plugin/") {
			c.File("./frontend/dist/index.html")
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
