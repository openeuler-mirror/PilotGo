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
 * LastEditTime: 2022-02-28 14:50:08
 * Description: provide agent rpm manager functions.
 ******************************************************************************/
package handlers

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func AllRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_all, err := agent.AllRpm()
	if err != nil {
		response.Fail(c, nil, "获取已安装rpm包列表失败!")
		return
	}
	response.Success(c, gin.H{"rpm_all": rpm_all}, "Success")
}
func RpmSourceHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_source, err := agent.RpmSource(rpmname)
	if err != nil {
		response.Fail(c, nil, "获取源软件包名以及源失败!")
		return
	}
	response.Success(c, gin.H{"rpm_source": rpm_source}, "Success")
}
func RpmInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_info, Err, err := agent.RpmInfo(rpmname)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "获取源软件包信息失败!")
		return
	} else {
		response.Success(c, gin.H{"rpm_info": rpm_info}, "Success")
	}

}
func InstallRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "软件包安装/卸载"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = rpmname
	logParent.Action = model.RPMInstall

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	rpm_install, err := agent.InstallRpm(rpmname)
	rpm := rpm_install.(string)
	if len(rpm) != 0 || err != nil {
		response.Fail(c, gin.H{"error": rpm}, "Failed!")

		logParent.StatusCode = 400
		logParent.Message = rpm
		mysqlmanager.DB.Save(&logParent)
		return
	} else {
		response.Success(c, nil, "该rpm包安装成功!")
		logParent.StatusCode = 200
		logParent.Message = "安装成功"
		mysqlmanager.DB.Save(&logParent)
	}

}
func RemoveRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")
	username := c.Query("userName")

	var logParent model.AgentLog
	var machineNode model.MachineNode
	mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)
	logParent.Type = "软件包安装/卸载"
	logParent.IP = machineNode.IP
	logParent.UserName = username
	logParent.OperationObject = rpmname
	logParent.Action = model.RPMRemove

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		logParent.StatusCode = 400
		logParent.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	rpm_remove, err := agent.RemoveRpm(rpmname)
	if err != nil {
		response.Fail(c, nil, "rpm包卸载命令执行失败!")
		logParent.StatusCode = 400
		logParent.Message = "卸载命令执行失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	if rpm_remove != nil {
		response.Fail(c, nil, "软件包卸载失败!")
		logParent.StatusCode = 400
		logParent.Message = "软件包卸载失败"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"rpm_remove": "卸载成功!"}, "Success")
	logParent.StatusCode = 200
	logParent.Message = "卸载成功"
	mysqlmanager.DB.Save(&logParent)
}
