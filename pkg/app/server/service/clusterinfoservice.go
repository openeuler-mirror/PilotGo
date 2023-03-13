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
package service

import (
	"errors"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

// 统计所有机器的状态
func AgentStatusCounts(machines []model.MachineNode) (normal, Offline, free int) {
	for _, agent := range machines {
		state := agent.State
		switch state {
		case global.Free:
			free++
		case global.OffLine:
			Offline++
		case global.Normal:
			normal++
		default:
			continue
		}
	}
	return
}

// 查找所有机器
func SelectAllMachine() ([]model.MachineNode, error) {
	machines, err := dao.AllMachine()
	if err != nil {
		return machines, err
	}
	if len(machines) == 0 {
		return nil, errors.New("未获取到机器")
	}
	return machines, nil
}

// 获取集群概览
func ClusterInfo() (dao.ClusterInfo, error) {
	data := dao.ClusterInfo{}
	machines, err := SelectAllMachine()
	if err != nil {
		return data, err
	}
	normal, Offline, free := AgentStatusCounts(machines)

	data.AgentTotal = len(machines)
	data.AgentStatus.Normal = normal
	data.AgentStatus.OffLine = Offline
	data.AgentStatus.Free = free
	return data, nil
}

// 获取各部门集群状态
func DepartClusterInfo() []dao.DepartMachineInfo {
	var departs []dao.DepartMachineInfo

	FirstDepartIds, err := dao.FirstDepartId()
	if err != nil {
		logger.Error(err.Error())
	}
	for _, depart_Id := range FirstDepartIds {
		Departids := make([]int, 0)
		Departids = append(Departids, depart_Id)
		ReturnSpecifiedDepart(depart_Id, &Departids) //某一级部门及其下属部门id

		lists, err := dao.SomeDepartMachine(Departids) //某一级部门及其下属部门所有机器
		if err != nil {
			logger.Error(err.Error())
		}
		departName, err := dao.DepartIdToGetDepartName(depart_Id)
		if err != nil {
			logger.Error(err.Error())
		}
		normal, Offline, free := AgentStatusCounts(lists)

		departInfo := dao.DepartMachineInfo{}
		departInfo.DepartName = departName
		departInfo.AgentStatus.Normal = normal
		departInfo.AgentStatus.OffLine = Offline
		departInfo.AgentStatus.Free = free

		departs = append(departs, departInfo)
	}
	return departs
}
