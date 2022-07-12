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
 * LastEditTime: 2022-07-12 14:10:23
 * Description: static router
 ******************************************************************************/
package resource

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed dist/static/css dist/static/fonts dist/static/img dist/static/js dist/static/pilotgo.ico dist/index.html
var Static embed.FS

func StaticRouter(r *gin.Engine) {
	r.StaticFS("/static", http.FS(Static))
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})
}
