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
 * LastEditTime: 2022-03-25 01:51:51
 * Description: provide agent rpm manager functions.
 ******************************************************************************/
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

type RPMS struct {
	UUIDs    []string `json:"uuid"`
	RPM      string   `json:"rpm"`
	UserName string   `json:"userName"`
}

func InstallRpmHandler(c *gin.Context) {
	var rpm RPMS
	var user model.User
	var logParent model.AgentLogParent
	var log model.AgentLog
	var machineNode model.MachineNode

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	bodys := string(body)

	err = json.Unmarshal([]byte(bodys), &rpm)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	logParent.Type = "软件包安装/卸载"
	logParent.UserName = rpm.UserName
	mysqlmanager.DB.Where("email = ?", rpm.UserName).Find(&user)
	logParent.DepartName = user.DepartName
	mysqlmanager.DB.Save(&logParent)
	StatusCodes := make([]string, 0)

	for _, uuid := range rpm.UUIDs {
		mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)

		log.IP = machineNode.IP
		log.OperationObject = rpm.RPM
		log.Action = model.RPMInstall
		log.LogParentID = logParent.ID

		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			response.Success(c, gin.H{"code": 400}, "获取uuid失败")

			log.StatusCode = 400
			log.Message = "获取uuid失败"
			mysqlmanager.DB.Save(&log)
			StatusCodes = append(StatusCodes, "400")
			continue
		}
		rpm_install, Err, err := agent.InstallRpm(rpm.RPM)
		if err != nil || len(Err) != 0 {
			response.Success(c, gin.H{"code": 400, "error": Err}, "Failed!")

			log.StatusCode = 400
			log.Message = Err
			mysqlmanager.DB.Save(&log)
			StatusCodes = append(StatusCodes, "400")
			continue
		} else {
			response.Success(c, gin.H{"code": 200, "install": rpm_install}, "该rpm包安装成功!")
			log.StatusCode = 200
			StatusCodes = append(StatusCodes, "200")
			log.Message = "安装成功"
			mysqlmanager.DB.Save(&log)
		}
	}
	var s int
	for _, success := range StatusCodes {
		if success == "200" {
			s++
		}
	}
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(s)/float64(len(StatusCodes))), 64)
	rate := strconv.FormatFloat(num, 'f', 2, 64)
	logParent.Status = strconv.Itoa(s) + "," + strconv.Itoa(len(StatusCodes)) + "," + rate
	mysqlmanager.DB.Save(&logParent)
}
func RemoveRpmHandler(c *gin.Context) {
	var rpm RPMS
	var user model.User
	var logParent model.AgentLogParent
	var log model.AgentLog
	var machineNode model.MachineNode

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &rpm)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	logParent.Type = "软件包安装/卸载"
	logParent.UserName = rpm.UserName
	mysqlmanager.DB.Where("email = ?", rpm.UserName).Find(&user)
	logParent.DepartName = user.DepartName
	mysqlmanager.DB.Save(&logParent)
	StatusCodes := make([]string, 0)
	for _, uuid := range rpm.UUIDs {
		mysqlmanager.DB.Where("machine_uuid=?", uuid).Find(&machineNode)

		log.IP = machineNode.IP
		log.OperationObject = rpm.RPM
		log.Action = model.RPMRemove
		log.LogParentID = logParent.ID

		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			response.Success(c, gin.H{"code": 400}, "获取uuid失败")
			log.StatusCode = 400
			log.Message = "获取uuid失败"
			mysqlmanager.DB.Save(&log)
			StatusCodes = append(StatusCodes, "400")
			continue
		}

		rpm_remove, Err, err := agent.RemoveRpm(rpm.RPM)
		if len(Err) != 0 || err != nil {
			response.Success(c, gin.H{"code": 400, "error": Err}, "Failed!")

			log.StatusCode = 400
			log.Message = Err
			mysqlmanager.DB.Save(&log)
			StatusCodes = append(StatusCodes, "400")
			continue
		} else {
			response.Success(c, gin.H{"code": 200, "remove": rpm_remove}, "卸载成功")
			log.StatusCode = 200
			StatusCodes = append(StatusCodes, "200")
			log.Message = "卸载成功"
			mysqlmanager.DB.Save(&log)
		}
	}
	var s int
	for _, success := range StatusCodes {
		if success == "200" {
			s++
		}
	}
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(s)/float64(len(StatusCodes))), 64)
	rate := strconv.FormatFloat(num, 'f', 2, 64)
	logParent.Status = strconv.Itoa(s) + "," + strconv.Itoa(len(StatusCodes)) + "," + rate
	mysqlmanager.DB.Save(&logParent)
}
