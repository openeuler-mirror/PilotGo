/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: wanghao
 * Date: 2022-02-18 13:03:16
 * LastEditTime: 2022-06-08 09:58:35
 * Description: provide machine manager functions.
 ******************************************************************************/
package controller

import (
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	machineservice "gitee.com/openeuler/PilotGo/app/server/service/machine"
	"gitee.com/openeuler/PilotGo/global"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type MachineModifyDepart struct {
	MachineID string `json:"machineid"`
	DepartID  int    `json:"departid"`
}

// 分页获取机器列表
func MachineInfoHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	depart := &machineservice.Depart{}
	if c.ShouldBind(&depart) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := machineservice.MachineInfo(depart, num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 资源池返回接口，查询没分配的机器，即departid=1,应该修改成维护状态的机器
func FreeMachineSource(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := machineservice.ReturnMachinePaged(global.UncateloguedDepartId, num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 返回所有机器指定字段，供插件使用
func MachineAllDataHandler(c *gin.Context) {
	datas, err := machineservice.MachineAllData()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, datas, "获取所有的机器数据")
}

// 删除机器
func DeleteMachineHandler(c *gin.Context) {
	var deleteuuid machineservice.DeleteUUID
	if c.Bind(&deleteuuid) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	machinelist := machineservice.DeleteMachine(deleteuuid.Deluuid)

	if len(machinelist) != 0 {
		response.Fail(c, gin.H{"machinelist": machinelist}, "机器删除失败")
	} else {
		response.Success(c, nil, "机器删除成功!")
	}
}

// 修改机器departid
func ModifyMachineDepartHandler(c *gin.Context) {
	var M MachineModifyDepart
	if err := c.Bind(&M); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := machineservice.ModifyMachineDepart(M.MachineID, M.DepartID)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "机器部门修改成功")
}

// 维护状态列表
func MaintStatusList(c *gin.Context) {
	var datas []string
	datas = append(datas, machineservice.NormalStatus)
	datas = append(datas, machineservice.MaintenanceStatus)
	response.Success(c, datas, "机器维护状态列表")
}
