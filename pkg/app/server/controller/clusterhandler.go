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
 * Date: 2022-03-24 00:46:05
 * LastEditTime: 2022-04-29 11:29:44
 * Description: 集群概览
 ******************************************************************************/
package controller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func ClusterInfo(c *gin.Context) {
	machines := dao.AllMachine()
	if len(machines) == 0 {
		response.Fail(c, nil, "未获取到机器")
	}

	normal, Offline, free := service.AgentStatusCounts(machines)

	data := model.ClusterInfo{}
	data.AgentTotal = len(machines)
	data.AgentStatus.Normal = normal
	data.AgentStatus.OffLine = Offline
	data.AgentStatus.Free = free

	response.Success(c, gin.H{"data": data}, "集群概览获取成功")
}

func DepartClusterInfo(c *gin.Context) {
	var departs []model.DepartMachineInfo

	FirstDepartIds := dao.FirstDepartId()
	for _, depart_Id := range FirstDepartIds {
		Departids := make([]int, 0)
		Departids = append(Departids, depart_Id)
		ReturnSpecifiedDepart(depart_Id, &Departids) //某一级部门及其下属部门id

		lists := dao.SomeDepartMachine(Departids) //某一级部门及其下属部门所有机器

		departName := dao.DepartIdToGetDepartName(depart_Id)
		normal, Offline, free := service.AgentStatusCounts(lists)

		departInfo := model.DepartMachineInfo{}
		departInfo.DepartName = departName
		departInfo.AgentStatus.Normal = normal
		departInfo.AgentStatus.OffLine = Offline
		departInfo.AgentStatus.Free = free

		departs = append(departs, departInfo)
	}
	response.Success(c, gin.H{"data": departs}, "获取各部门集群状态成功")
}
