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
 * Date: 2022-05-26 10:25:52
 * LastEditTime: 2023-07-11 19:25:23
 * Description: agent config file handler
 ******************************************************************************/

package agentcontroller

import (
	"net/http"
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/app/server/service"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ReadConfigFile(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	filepath := c.Query("file")
	result, Err, err := agent.ReadConfigFile(filepath)
	if err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, gin.H{"file": result}, "Success")
}

func GetAgentRepo(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	repos, Err, err := agent.GetRepoSource()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, repos, "获取到repo源")
}

func ConfigFileBroadcastToAgents(c *gin.Context) {
	fd := &dao.Frontdata{}
	if err := c.Bind(fd); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	batchIds := fd.FileBroadcast_BatchId
	UUIDs := dao.BatchIds2UUIDs(batchIds)

	path := fd.FileBroadcast_Path
	filename := fd.FileBroadcast_FileName
	text := fd.FileBroadcast_Text

	if len(path) == 0 {
		response.Fail(c, nil, "路径为空，请检查配置文件路径")
		return
	}
	if len(filename) == 0 {
		response.Fail(c, nil, "文件名为空，请检查配置文件名字")
		return
	}
	if len(text) == 0 {
		response.Fail(c, nil, "文件内容为空，请重新检查文件内容")
		return
	}

	log := auditlog.New(auditlog.LogTypeBroadcast, "配置文件下发", "", fd)
	auditlog.Add(log)

	statuscodes := []string{}

	for _, uuid := range UUIDs {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			log_s := auditlog.New_sub(log.LogUUID, agentmanager.GetAgent(uuid).IP, log.Action, "获取uuid失败", log.Module, path+"/"+filename, http.StatusBadRequest)
			auditlog.Add(log_s)
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.UpdateConfigFile(path, filename, text)
		if len(Err) != 0 || err != nil {
			log_s := auditlog.New_sub(log.LogUUID, agentmanager.GetAgent(uuid).IP, log.Action, Err, log.Module, path+"/"+filename, http.StatusBadRequest)
			auditlog.Add(log_s)
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			log_s := auditlog.New_sub(log.LogUUID, agentmanager.GetAgent(uuid).IP, log.Action, "配置文件下发成功", log.Module, path+"/"+filename, http.StatusOK)
			auditlog.Add(log_s)
			statuscodes = append(statuscodes, strconv.Itoa(http.StatusOK))
		}
	}

	status := service.BatchActionStatus(statuscodes)
	if err := auditlog.UpdateStatus(log, status); err != nil {
		logger.Error("failed to update father log status: %s", err.Error())
	}

	switch strings.Split(status, ",")[2] {
	case "0.00":
		response.Fail(c, nil, "配置文件下发失败")
		return
	case "1.00":
		response.Success(c, nil, "配置文件下发完成")
	default:
		response.Success(c, nil, "配置文件下发部分完成")
	}
}
