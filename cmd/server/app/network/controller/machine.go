/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	machineservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type MachineModifyDepart struct {
	MachineID string `json:"machineid"`
	DepartID  int    `json:"departid"`
}

// 获取机器列表
func MachineInfoNoPageHandler(c *gin.Context) {
	data, err := machineservice.MachineAll()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	resp := []*common.MachineNode{}
	for _, item := range data {
		d := &common.MachineNode{
			UUID:        item.UUID,
			IP:          item.IP,
			Department:  item.Departname,
			CPUArch:     item.CPU,
			OS:          item.Systeminfo,
			RunStatus:   item.Runstatus,
			MaintStatus: item.Maintstatus,
		}

		resp = append(resp, d)
	}

	response.Success(c, resp, "获取所有的机器数据")
}

// 分页获取机器列表
func MachineInfoHandler(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	search := c.Query("search")

	depart := &machineservice.Depart{}
	if c.ShouldBind(&depart) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	num := query.PageSize * (query.Page - 1)
	total, data, err := machineservice.MachineInfo(depart, num, query.PageSize, search)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), query)
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
