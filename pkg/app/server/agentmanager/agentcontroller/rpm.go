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
 * LastEditTime: 2022-04-13 01:51:51
 * Description: provide agent rpm manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

type RPMS struct {
	UUIDs        []string `json:"uuid"`
	RPM          string   `json:"rpm"`
	UserName     string   `json:"userName"`
	UserDeptName string   `json:"userDept"`
}

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
	var rpm RPMS
	c.Bind(&rpm)

	logParent := model.AgentLogParent{
		UserName:   rpm.UserName,
		DepartName: rpm.UserDeptName,
		Type:       model.LogTypeRPM,
	}
	logParentId := dao.ParentAgentLog(logParent)

	StatusCodes := make([]string, 0)

	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {

			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMInstall,
				StatusCode:      http.StatusBadRequest,
				Message:         "获取uuid失败",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.InstallRpm(rpm.RPM)
		if err != nil || len(Err) != 0 {

			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMInstall,
				StatusCode:      http.StatusBadRequest,
				Message:         Err,
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {

			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMInstall,
				StatusCode:      http.StatusOK,
				Message:         "安装成功",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}
	status := service.BatchActionStatus(StatusCodes)
	dao.UpdateParentAgentLog(logParentId, status)
	response.Success(c, nil, "该rpm包安装完成!")
}
func RemoveRpmHandler(c *gin.Context) {
	var rpm RPMS
	c.Bind(&rpm)

	logParent := model.AgentLogParent{
		UserName:   rpm.UserName,
		DepartName: rpm.UserDeptName,
		Type:       model.LogTypeRPM,
	}
	logParentId := dao.ParentAgentLog(logParent)

	StatusCodes := make([]string, 0)
	for _, uuid := range rpm.UUIDs {

		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMRemove,
				StatusCode:      http.StatusBadRequest,
				Message:         "获取uuid失败",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.RemoveRpm(rpm.RPM)
		if len(Err) != 0 || err != nil {

			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMRemove,
				StatusCode:      http.StatusBadRequest,
				Message:         Err,
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {

			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: rpm.RPM,
				Action:          model.RPMRemove,
				StatusCode:      http.StatusOK,
				Message:         "卸载成功",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}

	status := service.BatchActionStatus(StatusCodes)
	dao.UpdateParentAgentLog(logParentId, status)
	response.Success(c, nil, "该rpm包卸载完成!")
}
