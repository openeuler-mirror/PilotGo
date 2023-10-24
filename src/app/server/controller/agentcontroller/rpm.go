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
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/service"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/jwt"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
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
	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.New(auditlog.LogTypeRPM, "rpm软件包安装", user.ID)
	if err := auditlog.Add(log); err != nil {
		logger.Error(err.Error())
	}

	StatusCodes := make([]string, 0)

	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			message := "获取uuid失败" + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包安装", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		info, err := agent.AgentOverview()
		if err != nil {
			logger.Error(err.Error())
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			logger.Error("Install rpm is not supported on NestOS For Container")
			message := "Install rpm is not supported on NestOS For Container" + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包安装", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.InstallRpm(rpm.RPM)
		if err != nil || len(Err) != 0 {
			message := Err + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包安装", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			message := rpm.RPM + string(http.StatusOK)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包安装", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}
	status := service.BatchActionStatus(StatusCodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error(err.Error())
	}

	switch strings.Split(status, ",")[2] {
	case "0.00":
		response.Fail(c, nil, "软件包安装失败")
		return
	case "1.00":
		response.Success(c, nil, "软件包安装成功")
	default:
		response.Success(c, nil, "软件包安装部分成功")
	}
}

func RemoveRpmHandler(c *gin.Context) {
	var rpm RPMS
	c.Bind(&rpm)

	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.New(auditlog.LogTypeRPM, "rpm软件包卸载", user.ID)
	if err := auditlog.Add(log); err != nil {
		logger.Error(err.Error())
	}

	StatusCodes := make([]string, 0)
	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			message := "获取uuid失败" + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包卸载", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		info, err := agent.AgentOverview()
		if err != nil {
			logger.Error(err.Error())
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			logger.Error("Remove rpm is not supported on NestOS For Container")
			message := "Remove rpm is not supported on NestOS For Container" + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包卸载", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.RemoveRpm(rpm.RPM)
		if len(Err) != 0 || err != nil {
			message := Err + rpm.RPM + string(http.StatusBadRequest)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包卸载", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			message := rpm.RPM + string(http.StatusOK)
			log_s := auditlog.New_sub(auditlog.LogTypeRPM, "rpm软件包卸载", uuid, message, auditlog.StatusSuccess, user.ID)
			if err := auditlog.Add(log_s); err != nil {
				logger.Error(err.Error())
			}
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}

	status := service.BatchActionStatus(StatusCodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error(err.Error())
	}

	switch strings.Split(status, ",")[2] {
	case "0.00":
		response.Fail(c, nil, "软件包卸载失败")
		return
	case "1.00":
		response.Success(c, nil, "软件包卸载成功")
	default:
		response.Success(c, nil, "软件包卸载部分成功")
	}
}
