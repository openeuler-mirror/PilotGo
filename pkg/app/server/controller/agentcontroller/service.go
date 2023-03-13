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
 * Date: 2022-02-16 09:28:46
 * LastEditTime: 2022-04-11 17:07:55
 * Description: provide agent service manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

type AgentService struct {
	UUID         string `json:"uuid"`
	Service      string `json:"service"`
	UserName     string `json:"userName"`
	UserDeptName string `json:"userDept"`
}

func ServiceListHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_list, err := agent.ServiceList()
	if err != nil {
		response.Fail(c, nil, "获取服务列表失败!")
		return
	}
	response.Success(c, gin.H{"service_list": service_list}, "Success")
}
func ServiceStatusHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	service := c.Query("service")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	service_status, err := agent.ServiceStatus(service)
	if err != nil {
		response.Fail(c, nil, "获取服务状态失败!")
		return
	}
	response.Success(c, gin.H{"service_status": service_status}, "Success")
}

func ServiceStartHandler(c *gin.Context) {
	var agentservice AgentService
	c.Bind(&agentservice)

	logParent := dao.AgentLogParent{
		UserName:   agentservice.UserName,
		DepartName: agentservice.UserDeptName,
		Type:       global.LogTypeService,
	}
	logParentId, err := dao.ParentAgentLog(logParent)
	if err != nil {
		logger.Error(err.Error())
	}
	agent := agentmanager.GetAgent(agentservice.UUID)
	UUID_iP, err := dao.UUID2MacIP(agentservice.UUID)
	if err != nil {
		logger.Error(err.Error())
	}
	if agent == nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceStart,
			StatusCode:      http.StatusBadRequest,
			Message:         "获取uuid失败",
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, nil, "获取uuid失败")
		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	service_start, Err, err := agent.ServiceStart(agentservice.Service)
	if len(Err) != 0 || err != nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceStart,
			StatusCode:      http.StatusBadRequest,
			Message:         Err,
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, gin.H{"error": Err}, "Failed!")

		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	log := dao.AgentLog{
		LogParentID:     logParentId,
		IP:              UUID_iP,
		OperationObject: agentservice.Service,
		Action:          global.ServiceStart,
		StatusCode:      http.StatusOK,
		Message:         "启动服务成功",
	}
	if dao.AgentLogMessage(log) != nil {
		logger.Error(err.Error())
	}
	if dao.UpdateParentAgentLog(logParentId, global.ActionOK) != nil {
		logger.Error(err.Error())
	}

	response.Success(c, gin.H{"service_start": service_start}, "Success")
}
func ServiceStopHandler(c *gin.Context) {
	var agentservice AgentService
	c.Bind(&agentservice)

	logParent := dao.AgentLogParent{
		UserName:   agentservice.UserName,
		DepartName: agentservice.UserDeptName,
		Type:       global.LogTypeService,
	}
	logParentId, err := dao.ParentAgentLog(logParent)
	if err != nil {
		logger.Error(err.Error())
	}
	agent := agentmanager.GetAgent(agentservice.UUID)
	UUID_iP, err := dao.UUID2MacIP(agentservice.UUID)
	if err != nil {
		logger.Error(err.Error())
	}
	if agent == nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceStop,
			StatusCode:      http.StatusBadRequest,
			Message:         "获取uuid失败",
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, nil, "获取uuid失败")

		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	service_stop, Err, err := agent.ServiceStop(agentservice.Service)
	if len(Err) != 0 || err != nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceStop,
			StatusCode:      http.StatusBadRequest,
			Message:         Err,
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, gin.H{"error": Err}, "Failed!")

		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	log := dao.AgentLog{
		LogParentID:     logParentId,
		IP:              UUID_iP,
		OperationObject: agentservice.Service,
		Action:          global.ServiceStop,
		StatusCode:      http.StatusOK,
		Message:         "关闭服务成功",
	}
	if dao.AgentLogMessage(log) != nil {
		logger.Error(err.Error())
	}
	if dao.UpdateParentAgentLog(logParentId, global.ActionOK) != nil {
		logger.Error(err.Error())
	}

	response.Success(c, gin.H{"service_stop": service_stop}, "Success")
}
func ServiceRestartHandler(c *gin.Context) {
	var agentservice AgentService
	c.Bind(&agentservice)

	logParent := dao.AgentLogParent{
		UserName:   agentservice.UserName,
		DepartName: agentservice.UserDeptName,
		Type:       global.LogTypeService,
	}
	logParentId, err := dao.ParentAgentLog(logParent)
	if err != nil {
		logger.Error(err.Error())
	}
	agent := agentmanager.GetAgent(agentservice.UUID)
	UUID_iP, err := dao.UUID2MacIP(agentservice.UUID)
	if err != nil {
		logger.Error(err.Error())
	}
	if agent == nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceRestart,
			StatusCode:      http.StatusBadRequest,
			Message:         "获取uuid失败",
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, nil, "获取uuid失败")

		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	service_restart, Err, err := agent.ServiceRestart(agentservice.Service)
	if len(Err) != 0 || err != nil {

		log := dao.AgentLog{
			LogParentID:     logParentId,
			IP:              UUID_iP,
			OperationObject: agentservice.Service,
			Action:          global.ServiceRestart,
			StatusCode:      http.StatusBadRequest,
			Message:         Err,
		}
		if dao.AgentLogMessage(log) != nil {
			logger.Error(err.Error())
		}
		response.Fail(c, gin.H{"error": Err}, "Failed!")

		if dao.UpdateParentAgentLog(logParentId, global.ActionFalse) != nil {
			logger.Error(err.Error())
		}
		return
	}

	log := dao.AgentLog{
		LogParentID:     logParentId,
		IP:              UUID_iP,
		OperationObject: agentservice.Service,
		Action:          global.ServiceRestart,
		StatusCode:      http.StatusOK,
		Message:         "重启服务成功",
	}
	if dao.AgentLogMessage(log) != nil {
		logger.Error(err.Error())
	}
	if dao.UpdateParentAgentLog(logParentId, global.ActionOK) != nil {
		logger.Error(err.Error())
	}

	response.Success(c, gin.H{"service_restart": service_restart}, "Success")
}
