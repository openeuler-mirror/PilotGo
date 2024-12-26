/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
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
	if err := c.Bind(&rpm); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	if !utils.CheckString(rpm.RPM) {
		response.Fail(c, nil, "软件包名有除_+-.以外的特殊字符")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleMachine,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "rpm软件包安装",
	}
	auditlog.Add(log)

	StatusCodes := make([]string, 0)

	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)
		log_s := &auditlog.AuditLog{
			LogUUID:    uuidservice.New().String(),
			ParentUUID: log.LogUUID,
			Module:     auditlog.ModuleMachine,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "rpm软件包安装",
			Message:    "agentuuid:" + uuid,
		}
		auditlog.Add(log_s)
		if agent == nil {
			message := "获取uuid失败" + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		info, err := agent.AgentOverview()
		if err != nil {
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			logger.Error("%v",err.Error())
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			logger.Error("Install rpm is not supported on NestOS For Container")
			message := "Install rpm is not supported on NestOS For Container" + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.InstallRpm(rpm.RPM)
		if err != nil || len(Err) != 0 {
			message := Err + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			message := rpm.RPM + strconv.Itoa(http.StatusOK)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}
	status := auditlog.BatchActionStatus(StatusCodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error("%v",err.Error())
	}

	switch strings.Split(status, ",")[2] {
	case "0.00":
		response.Fail(c, nil, "软件包安装失败")
		return
	case "1.00":
		response.Success(c, nil, "软件包安装成功")
		return
	default:
		response.Success(c, nil, "软件包安装部分成功")
	}
}

func RemoveRpmHandler(c *gin.Context) {
	var rpm RPMS
	if err := c.Bind(&rpm); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	if !utils.CheckString(rpm.RPM) {
		response.Fail(c, nil, "软件包名有除_+-.以外的特殊字符")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleMachine,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "rpm软件包卸载",
	}
	auditlog.Add(log)

	StatusCodes := make([]string, 0)
	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)
		log_s := &auditlog.AuditLog{
			LogUUID:    uuidservice.New().String(),
			ParentUUID: log.LogUUID,
			Module:     auditlog.ModuleMachine,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "rpm软件包卸载",
			Message:    "agentuuid:" + uuid,
		}
		auditlog.Add(log_s)

		if agent == nil {
			message := "获取uuid失败" + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		info, err := agent.AgentOverview()
		if err != nil {
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			logger.Error("%v",err.Error())
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			logger.Error("Remove rpm is not supported on NestOS For Container")
			message := "Remove rpm is not supported on NestOS For Container" + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.RemoveRpm(rpm.RPM)
		if len(Err) != 0 || err != nil {
			message := Err + rpm.RPM + strconv.Itoa(http.StatusBadRequest)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			message := rpm.RPM + strconv.Itoa(http.StatusOK)
			auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+message)
			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}

	status := auditlog.BatchActionStatus(StatusCodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error("%v",err.Error())
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
