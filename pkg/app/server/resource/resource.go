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
 * Date: 2022-07-12 13:03:16
 * LastEditTime: 2023-07-04 14:31:05
 * Description: static router
 ******************************************************************************/
package resource

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

// //go:embed css fonts img js pilotgo.ico index.html
// var Static embed.FS

func StaticRouter(router *gin.Engine) *gin.Engine {
	// router.StaticFS("/static", http.FS(Static))
	// t, err := template.ParseFS(Static, "index.html")
	// if err != nil {
	// 	logger.Error("parse template failed !!!")
	// }
	// router.SetHTMLTemplate(t)
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })

	// TODO: for test
	router.Static("/static", "./frontend/dist/static")
	router.StaticFile("/", "./frontend/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Info("process noroute: %s", c.Request.URL.RawPath)
		// c.HTML(http.StatusOK, "index.html", nil)

		// TODO: for test
		if !strings.HasPrefix(c.Request.RequestURI, "/api/") && !strings.HasPrefix(c.Request.RequestURI, "/plugin/") {
			c.File("./frontend/dist/index.html")
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
	return router
}
