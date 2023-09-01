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
 * LastEditTime: 2023-09-01 09:34:45
 * Description: casbin权限控制
 ******************************************************************************/
package auth

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
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
		ok, err := G_Enfocer.Enforce(sub, obj, act)
		if err == nil {
			if ok {
				logger.Info("恭喜您,权限验证通过")
				c.Next()
			} else {
				logger.Error("很遗憾,权限验证没有通过")
				c.Abort()
			}
		}
		logger.Error("auth check error:%s", err)
		c.Abort()
	}
}

func AllPolicy1() (interface{}, int) { //暂时没有调用，与casbin文件中的方法重复，所以暂时修改为AllPolicy1
	casbin := make([]map[string]interface{}, 0)
	list := G_Enfocer.GetPolicy()
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
