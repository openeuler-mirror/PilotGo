//go:build production
// +build production

/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */

package resource

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

//go:embed all:assets index.html pilotgo.ico
var StaticFiles embed.FS

func StaticRouter(router *gin.Engine) {
	sf, err := fs.Sub(StaticFiles, "assets")
	if err != nil {
		logger.Error("failed to load frontend assets files: %s", err.Error())
		return
	}

	router.StaticFS("/assets", http.FS(sf))
	router.GET("/", func(c *gin.Context) {
		c.FileFromFS("/", http.FS(StaticFiles))
	})
	router.GET("/pilotgo.ico", func(c *gin.Context) {
		c.FileFromFS("/pilotgo.ico", http.FS(StaticFiles))
	})

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Debug("process noroute: %s", c.Request.URL.String())
		if !strings.HasPrefix(c.Request.RequestURI, "/api/") && !strings.HasPrefix(c.Request.RequestURI, "/plugin/") {
			c.FileFromFS("/", http.FS(StaticFiles))
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
