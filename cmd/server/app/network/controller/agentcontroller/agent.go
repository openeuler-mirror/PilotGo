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

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func AgentOverviewHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		c.JSON(http.StatusOK, `{"status":-1}`)
	}

	info, err := agent.AgentOverview()
	if err != nil {
		response.Fail(c, nil, err.Error())
	}

	node, err := machine.MachineInfoByUUID(uuid)
	if err != nil {
		response.Fail(c, gin.H{"IP": node.IP, "state": node.Maintstatus, "depart": node.Departname}, err.Error())
	}

	type DiskUsage struct {
		Device      string `json:"device"`
		Path        string `json:"path"`
		Total       string `json:"total"`
		UsedPercent string `json:"used_percent"`
	}
	dus := []DiskUsage{}
	for _, du := range info.DiskUsage {
		dus = append(dus, DiskUsage{
			Device:      du.Device,
			Path:        du.Path,
			Total:       du.Total,
			UsedPercent: du.UsedPercent,
		})
	}

	result := struct {
		IP              string      `json:"ip"`
		Department      string      `json:"department"`
		State           string      `json:"state"`
		Platform        string      `json:"platform"`
		PlatformVersion string      `json:"platform_version"`
		KernelArch      string      `json:"kernel_arch"`
		KernelVersion   string      `json:"kernel_version"`
		CPUNum          int         `json:"cpu_num"`
		ModelName       string      `json:"model_name"`
		MemoryTotal     int64       `json:"memory_total"`
		DiskUsage       []DiskUsage `json:"disk_usage"`
		Immutable       bool        `json:"immutable"`
	}{
		IP:              info.IP,
		Department:      node.Departname,
		Platform:        info.SysInfo.Platform,
		PlatformVersion: info.SysInfo.PlatformVersion,
		KernelArch:      info.SysInfo.KernelArch,
		KernelVersion:   info.SysInfo.KernelVersion,
		CPUNum:          info.CpuInfo.CpuNum,
		ModelName:       info.CpuInfo.ModelName,
		MemoryTotal:     info.MemoryInfo.MemTotal,
		DiskUsage:       dus,
		Immutable:       info.IsImmutable,
		State:           node.Runstatus,
	}

	response.Success(c, result, "Success")
}

func AgentListHandler(c *gin.Context) {
	logger.Debug("process get agent list request")

	agent_list := agentmanager.GetAgentList()

	c.JSON(http.StatusOK, agent_list)
}

func OsBasic(c *gin.Context) {
	uuid := c.Query("uuid")
	node, err := machine.MachineInfoByUUID(uuid)
	if err != nil {
		response.Fail(c, gin.H{"IP": node.IP, "state": node.Maintstatus, "depart": node.Departname}, err.Error())
	}
	response.Success(c, gin.H{"IP": node.IP, "state": node.Maintstatus, "depart": node.Departname}, "Success")
}
