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
 * LastEditTime: 2022-03-14 10:35:54
 * Description: casbin权限控制
 ******************************************************************************/
package middleware

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

// 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Query("roleId")
		// 获取请求的PATH
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := role
		//判断策略中是否存在
		if ok := global.PILOTGO_E.Enforce(sub, obj, act); ok {
			logger.Info("恭喜您,权限验证通过")
			c.Next()
		} else {
			logger.Fatal("很遗憾,权限验证没有通过")
			c.Abort()
		}
	}
}

func AllPolicy() (interface{}, int) {
	casbin := make([]map[string]interface{}, 0)
	list := global.PILOTGO_E.GetPolicy()
	for _, vlist := range list {
		policy := make(map[string]interface{})
		policy["role"] = vlist[0]
		policy["url"] = vlist[1]
		policy["method"] = vlist[2]
		casbin = append(casbin, policy)
	}
	total := len(casbin)
	return casbin, total
}

func PolicyRemove(rule model.CasbinRule) bool {
	if ok := global.PILOTGO_E.RemovePolicy(rule.RoleType, rule.Url, rule.Method); !ok {
		return false
	} else {
		return true
	}
}

func PolicyAdd(rule model.CasbinRule) bool {
	if ok := global.PILOTGO_E.AddPolicy(rule.RoleType, rule.Url, rule.Method); !ok {
		return false
	} else {
		return true
	}
}
