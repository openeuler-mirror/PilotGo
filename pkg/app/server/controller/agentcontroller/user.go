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
 * Date: 2022-02-17 02:43:29
 * LastEditTime: 2022-02-24 14:32:25
 * Description: provide agent user manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/sdk/response"
)

func CurrentUserInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_info, err := agent.CurrentUser()
	if err != nil {
		response.Fail(c, nil, "获取当前登录用户信息失败!")
		return
	}
	response.Success(c, gin.H{"user_info": user_info}, "获取当前登录用户信息成功!")
}
func AllUserInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_all, err := agent.AllUser()
	if err != nil {
		response.Fail(c, nil, "获取机器所有用户数据失败!")
		return
	}
	response.Success(c, gin.H{"user_all": user_all}, "获取机器所有用户数据成功!")
}
func AddLinuxUserHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	username := c.Query("username")
	password := c.Query("password")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_add, err := agent.AddLinuxUser(username, password)
	if err != nil {
		response.Fail(c, nil, "新增用户失败!")
		return
	}
	response.Success(c, gin.H{"user_add": user_add}, "新增用户成功!")
}
func DelUserHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	username := c.Query("username")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_del, Err, err := agent.DelUser(username)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "删除用户失败!")
		return
	} else {
		response.Success(c, gin.H{"user_del": user_del}, "删除用户成功!")
	}

}
func ChangeFileOwnerHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	user := c.Query("username")
	file := c.Query("file")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_ower, err := agent.ChangeFileOwner(user, file)
	if err != nil {
		response.Fail(c, nil, "改变文件或目录所有者失败!")
		return
	}
	response.Success(c, gin.H{"user_ower": user_ower}, "改变文件或目录所有者成功!")
}
func ChangePermissionHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	permission := c.Query("per")
	file := c.Query("file")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	user_per, err := agent.ChangePermission(permission, file)
	if err != nil {
		response.Fail(c, nil, "改变文件权限失败!")
		return
	}
	response.Success(c, gin.H{"user_per": user_per}, "改变文件权限成功!")
}
