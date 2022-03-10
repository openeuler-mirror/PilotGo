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
 * Date: 2022-03-07 15:25:53
 * LastEditTime: 2022-03-10 11:18:22
 * Description: casbin权限控制
 ******************************************************************************/
package common

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/dto"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

var (
	E *casbin.Enforcer
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("x-user")
		if !ok {
			response.Response(c, http.StatusOK,
				200,
				nil,
				"登录已失效")
		}
		var User model.User
		email := dto.ToUserDto(user.(model.User)).Email
		mysqlmanager.DB.Where("email=?", email).Find(&User)
		// 获取请求的PATH
		obj := c.Request.URL.Path
		// 获取请求方法
		act := User.RoleType
		// 获取用户的角色
		sub := email
		//判断策略中是否存在
		if ok := E.Enforce(sub, obj, act); ok {
			logger.Info("恭喜您,权限验证通过")
			c.Next()
		} else {
			logger.Fatal("很遗憾,权限验证没有通过")
			c.Abort()
		}
	}
}
