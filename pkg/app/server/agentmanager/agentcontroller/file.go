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
 * LastEditTime: 2022-06-02 10:16:10
 * Description: agent config file handler
 ******************************************************************************/

package agentcontroller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	uos "openeluer.org/PilotGo/PilotGo/pkg/utils/os"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func ReadFile(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	filepath := c.Query("file")
	result, Err, err := agent.ReadFile(filepath)
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
	response.JSON(c, http.StatusOK, http.StatusOK, repos, "获取到repo源")
}

func GetAgentNetworkConnect(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net, Err, err := agent.GetNetWorkConnInfo()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.JSON(c, http.StatusOK, http.StatusOK, net, "获取到网络连接信息")
}

func ConfigNetworkConnect(c *gin.Context) {
	var network uos.NetworkConfig
	c.Bind(&network)

	ip_assignment := network.BootProto
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ip分配方式不能为空")
		return
	}
	ipv4_addr := network.IPAddr
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4地址不能为空")
		return
	}
	ipv4_netmask := network.NetMask
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4子网掩码不能为空")
		return
	}
	if ok := strings.Contains(ipv4_netmask, "."); !ok {
		prefix, _ := strconv.Atoi(ipv4_netmask)
		ipv4_netmask = uos.LenToSubNetMask(prefix)
	}
	ipv4_gateway := network.GateWay
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4网关不能为空")
		return
	}
	ipv4_dns1 := network.DNS1
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4 DNS1 不能为空")
		return
	}

	agent := agentmanager.GetAgent(network.MachineUUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	nic_name, Err, err := agent.GetNICName()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}

	oldnet, Err, err := agent.GetNetWorkConnectInfo()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	oldnets := controller.InterfaceToSlice(oldnet)

	switch ip_assignment {
	case "static":
		text := service.NetworkStatic(oldnets, ipv4_addr, ipv4_netmask, ipv4_gateway, ipv4_dns1, network.DNS2)
		_, Err, err := agent.UpdateFile(global.NetWorkPath, nic_name.(string), text)
		if len(Err) != 0 || err != nil {
			response.JSON(c, http.StatusOK, http.StatusOK, nil, Err)
			return
		}
		_, Err, err = agent.RestartNetWork(nic_name.(string))
		if len(Err) != 0 || err != nil {
			response.JSON(c, http.StatusOK, http.StatusOK, nil, Err)
			return
		}
		response.JSON(c, http.StatusOK, http.StatusOK, nil, "网络配置更新成功")

	case "dhcp":
		text := service.NetworkDHCP(oldnets)
		_, Err, err := agent.UpdateFile(global.NetWorkPath, nic_name.(string), text)
		if len(Err) != 0 || err != nil {
			response.JSON(c, http.StatusOK, http.StatusOK, nil, Err)
			return
		}
		_, Err, err = agent.RestartNetWork(nic_name.(string))
		if len(Err) != 0 || err != nil {
			response.JSON(c, http.StatusOK, http.StatusOK, nil, Err)
			return
		}
		response.JSON(c, http.StatusOK, http.StatusOK, nil, "网络配置更新成功")

	default:
		response.Fail(c, nil, "请重新检查ip分配方式")
	}
}

func FileBroadcastToAgents(c *gin.Context) {
	var fb model.FileBroadcast
	c.Bind(&fb)

	batchIds := fb.BatchId
	UUIDs := dao.BatchIds2UUIDs(batchIds)

	path := fb.Path
	filename := fb.FileName
	text := fb.Text

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
	logParent := model.AgentLogParent{
		UserName:   fb.User,
		DepartName: fb.UserDept,
		Type:       global.LogTypeBroadcast,
	}
	logParentId := dao.ParentAgentLog(logParent)

	StatusCodes := make([]string, 0)

	for _, uuid := range UUIDs {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: filename,
				Action:          global.BroadcastFile,
				StatusCode:      http.StatusBadRequest,
				Message:         "获取uuid失败",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		_, Err, err := agent.UpdateFile(path, filename, text)
		if len(Err) != 0 || err != nil {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: filename,
				Action:          global.BroadcastFile,
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
				OperationObject: filename,
				Action:          global.BroadcastFile,
				StatusCode:      http.StatusOK,
				Message:         "配置文件下发成功",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}
	status := service.BatchActionStatus(StatusCodes)
	dao.UpdateParentAgentLog(logParentId, status)
	response.Success(c, nil, "配置文件下发完成!")
}
