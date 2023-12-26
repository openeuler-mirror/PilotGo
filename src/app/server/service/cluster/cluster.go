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
 * Date: 2021-04-29 09:08:08
 * LastEditTime: 2022-04-29 09:25:41
 * Description: 集群概览业务逻辑
 ******************************************************************************/
package cluster

import (
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/app/server/service/machine"
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
	Normal  int `json:"normal"`
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

	//Normal显示的事在线机器数量，前端字段需要做修改
	data.AgentStatus.Normal = online
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
		maint, err := dao.CountMaintStatus(machine.MaintenanceStatus, depart.ID)
		if err != nil {
			logger.Error(err.Error())
		}
		departInfo.AgentStatus.Normal = online
		departInfo.AgentStatus.OffLine = offline
		departInfo.AgentStatus.Free = maint

		departs = append(departs, departInfo)
	}
	return departs
}
