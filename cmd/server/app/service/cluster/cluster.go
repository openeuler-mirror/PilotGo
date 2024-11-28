/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package cluster

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type ClusterInfoParam struct {
	AgentTotal  int `json:"total"`
	AgentStatus AgentStatus
}

type DepartMachineInfo struct {
	DepartName  string `json:"depart"`
	AgentStatus AgentStatus
}

type AgentStatus struct {
	Online  int `json:"normal"` //normal显示的是在线机器数量，前端字段需要做修改
	OffLine int `json:"offline"`
	Free    int `json:"free"`
}

// 获取集群概览
func ClusterInfo() (ClusterInfoParam, error) {
	data := ClusterInfoParam{}
	count, err := dao.CountMachineNode(nil)
	if err != nil {
		return data, err
	}
	data.AgentTotal = count
	//从数据库进行状态统计
	online, err := dao.CountRunStatus(machine.OnlineStatus, nil)
	if err != nil {
		return data, err
	}
	offline, err := dao.CountRunStatus(machine.OfflineStatus, nil)
	if err != nil {
		return data, err
	}
	maint, err := dao.CountMaintStatus(machine.MaintenanceStatus, nil)
	if err != nil {
		return data, err
	}

	data.AgentStatus.Online = online
	data.AgentStatus.OffLine = offline
	data.AgentStatus.Free = maint
	return data, nil
}

// 获取各部门集群状态
func DepartClusterInfo() []DepartMachineInfo {
	var departs []DepartMachineInfo
	//获取所有部门
	departnode, err := dao.GetAllDepart()
	if err != nil {
		logger.Error(err.Error())
	}
	//获取每个部门的机器状态数量
	for _, depart := range departnode {
		departInfo := DepartMachineInfo{}
		departInfo.DepartName = depart.Depart
		online, err := dao.CountRunStatus(machine.OnlineStatus, depart.ID)
		if err != nil {
			logger.Error(err.Error())
		}
		offline, err := dao.CountRunStatus(machine.OfflineStatus, depart.ID)
		if err != nil {
			logger.Error(err.Error())
		}

		//TODO:各部门中不存在维护状态的机器，应该统计正常使用的机器的数量，只有根节点下才会存在维护状态的机器
		maint, err := dao.CountMaintStatus(machine.MaintenanceStatus, depart.ID)
		if err != nil {
			logger.Error(err.Error())
		}
		departInfo.AgentStatus.Online = online
		departInfo.AgentStatus.OffLine = offline
		departInfo.AgentStatus.Free = maint

		departs = append(departs, departInfo)
	}
	return departs
}
