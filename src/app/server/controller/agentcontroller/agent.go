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
 * Date: 2022-01-24 15:08:08
 * LastEditTime: 2023-09-08 16:45:13
 * Description: Get the basic information of the machine
 ******************************************************************************/
package agentcontroller

import (
	"net/http"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func AgentOverviewHandler(c *gin.Context) {
	logger.Debug("process get agent overview request")
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		c.JSON(http.StatusOK, `{"status":-1}`)
	}

	info, err := agent.AgentOverview()
	if err != nil {
		response.Fail(c, nil, err.Error())
	}

	ip, state, dept, err := dao.MachineBasic(uuid)
	if err != nil {
		response.Fail(c, gin.H{"IP": ip, "state": state, "depart": dept}, err.Error())
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
		State           int         `json:"state"`
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
		Department:      dept,
		State:           state,
		Platform:        info.SysInfo.Platform,
		PlatformVersion: info.SysInfo.PlatformVersion,
		KernelArch:      info.SysInfo.KernelArch,
		KernelVersion:   info.SysInfo.KernelVersion,
		CPUNum:          info.CpuInfo.CpuNum,
		ModelName:       info.CpuInfo.ModelName,
		MemoryTotal:     info.MemoryInfo.MemTotal,
		DiskUsage:       dus,
		Immutable:       info.IsImmutable,
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
	ip, state, dept, err := dao.MachineBasic(uuid)
	if err != nil {
		response.Fail(c, gin.H{"IP": ip, "state": state, "depart": dept}, err.Error())
	}
	response.Success(c, gin.H{"IP": ip, "state": state, "depart": dept}, "Success")
}
