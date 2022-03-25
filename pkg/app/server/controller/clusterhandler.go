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
 * LastEditTime: 2022-03-24 06:29:44
 * Description: 集群概览
 ******************************************************************************/
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func ClusterInfo(c *gin.Context) {
	var agents []model.MachineNode
	var nolmal, Offline, free int
	data := make(map[string]interface{})
	mysqlmanager.DB.Find(&agents)
	total := len(agents)
	if total == 0 {
		response.Response(c, http.StatusOK, 400, gin.H{"data": data}, "未获取到机器")
	}
	data["total"] = total
	for _, agent := range agents {
		state := agent.State
		switch state {
		case model.Free:
			free++
		case model.OffLine:
			Offline++
		case model.Normal:
			nolmal++
		default:
			continue
		}
	}
	data["normal"] = nolmal
	data["offline"] = Offline
	data["free"] = free
	response.Response(c, http.StatusOK, 200, gin.H{"data": data}, "集群概览获取成功")
}

func DepartClusterInfo(c *gin.Context) {
	var agents []model.MachineNode
	mysqlmanager.DB.Find(&agents)
	departId := make([]int, 0)
	for _, agent := range agents {
		depart := model.DepartNode{}
		mysqlmanager.DB.Where("id = ?", agent.DepartId).Find(&depart)
		if depart.PID == 1 {
			departId = append(departId, depart.ID)
		}
	}
	departs := make([]AgentStatus, 0)
	departId = RemoveRepeated(departId)
	for _, depart_Id := range departId {
		Dids := make([]int, 0)
		ReturnSpecifiedDepart(depart_Id, &Dids)
		Dids = append(Dids, depart_Id)
		list := []model.MachineNode{}
		lists := []model.MachineNode{}
		for _, id := range Dids {
			mysqlmanager.DB.Where("depart_id = ?", id).Find(&list)
			lists = append(lists, list...)
		}
		dep := model.DepartNode{}
		mysqlmanager.DB.Where("id = ?", depart_Id).Find(&dep)
		status := AgentStatus{}
		status.Depart = dep.Depart
		for _, list := range lists {
			state := list.State
			switch state {
			case model.Free:
				status.Free++
			case model.OffLine:
				status.OffLine++
			case model.Normal:
				status.Normal++
			default:
				continue
			}
		}
		departs = append(departs, status)
	}
	response.Response(c, http.StatusOK, 200, gin.H{"data": departs}, "获取各部门集群状态成功")
}

// 去重
func RemoveRepeated(s []int) []int {
	result := make([]int, 0)
	m := make(map[int]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

type AgentStatus struct {
	Depart  string `json:"depart"`
	Normal  int    `json:"normal"`
	OffLine int    `json:"offline"`
	Free    int    `json:"free"`
}
