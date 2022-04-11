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
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

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

type AgentService struct {
	UUID     string `json:"uuid"`
	Service  string `json:"service"`
	UserName string `json:"userName"`
}

func ServiceStartHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var AS AgentService
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &AS)
	fmt.Println(bodys)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var logParent model.AgentLogParent
	var user model.User
	var log model.AgentLog
	var machineNode model.MachineNode

	logParent.Type = "运行服务"
	logParent.UserName = AS.UserName
	mysqlmanager.DB.Where("email = ?", AS.UserName).Find(&user)
	logParent.DepartName = user.DepartName
	mysqlmanager.DB.Save(&logParent)

	mysqlmanager.DB.Where("machine_uuid=?", AS.UUID).Find(&machineNode)

	log.IP = machineNode.IP
	log.OperationObject = AS.Service
	log.Action = model.ServiceStart
	log.LogParentID = logParent.ID

	agent := agentmanager.GetAgent(AS.UUID)
	if agent == nil {
		response.Success(c, gin.H{"code": 400}, "获取uuid失败")

		log.StatusCode = 400
		log.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_start, Err, err := agent.ServiceStart(AS.Service)
	if len(Err) != 0 || err != nil {
		response.Success(c, gin.H{"code": 400, "error": Err}, "Failed!")

		log.StatusCode = 400
		log.Message = Err
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"code": 200, "service_start": service_start}, "Success")
	log.StatusCode = 200
	log.Message = "启动服务成功"
	mysqlmanager.DB.Save(&log)
	logParent.Status = "1,1,1.00"
	mysqlmanager.DB.Save(&logParent)
}
func ServiceStopHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var AS AgentService
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &AS)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	var logParent model.AgentLogParent
	var user model.User
	var log model.AgentLog
	var machineNode model.MachineNode

	logParent.Type = "运行服务"
	logParent.UserName = AS.UserName
	mysqlmanager.DB.Where("email = ?", AS.UserName).Find(&user)
	logParent.DepartName = user.DepartName
	mysqlmanager.DB.Save(&logParent)

	mysqlmanager.DB.Where("machine_uuid=?", AS.UUID).Find(&machineNode)

	log.IP = machineNode.IP
	log.OperationObject = AS.Service
	log.Action = model.ServiceStop
	log.LogParentID = logParent.ID

	agent := agentmanager.GetAgent(AS.UUID)
	if agent == nil {
		response.Success(c, gin.H{"code": 400}, "获取uuid失败")

		log.StatusCode = 400
		log.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_stop, Err, err := agent.ServiceStop(AS.Service)
	if len(Err) != 0 || err != nil {
		response.Success(c, gin.H{"code": 400, "error": Err}, "Failed!")

		log.StatusCode = 400
		log.Message = Err
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"code": 200, "service_stop": service_stop}, "Success")
	log.StatusCode = 200
	log.Message = "关闭服务成功"
	mysqlmanager.DB.Save(&log)
	logParent.Status = "1,1,1.00"
	mysqlmanager.DB.Save(&logParent)
}
func ServiceRestartHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var AS AgentService
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &AS)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	var logParent model.AgentLogParent
	var user model.User
	var log model.AgentLog
	var machineNode model.MachineNode

	logParent.Type = "运行服务"
	logParent.UserName = AS.UserName
	mysqlmanager.DB.Where("email = ?", AS.UserName).Find(&user)
	logParent.DepartName = user.DepartName
	mysqlmanager.DB.Save(&logParent)

	mysqlmanager.DB.Where("machine_uuid=?", AS.UUID).Find(&machineNode)

	log.IP = machineNode.IP
	log.OperationObject = AS.Service
	log.Action = model.ServiceRestart
	log.LogParentID = logParent.ID

	agent := agentmanager.GetAgent(AS.UUID)
	if agent == nil {
		response.Success(c, gin.H{"code": 400}, "获取uuid失败")

		log.StatusCode = 400
		log.Message = "获取uuid失败"
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}

	service_restart, Err, err := agent.ServiceRestart(AS.Service)
	if len(Err) != 0 || err != nil {
		response.Success(c, gin.H{"code": 400, "error": Err}, "重启服务失败!")
		log.StatusCode = 400
		log.Message = Err
		mysqlmanager.DB.Save(&log)
		logParent.Status = "0,1,0.00"
		mysqlmanager.DB.Save(&logParent)
		return
	}
	response.Success(c, gin.H{"code": 200, "service_restart": service_restart}, "Success")
	log.StatusCode = 200
	log.Message = "重启服务成功"
	mysqlmanager.DB.Save(&log)
	logParent.Status = "1,1,1.00"
	mysqlmanager.DB.Save(&logParent)
}
