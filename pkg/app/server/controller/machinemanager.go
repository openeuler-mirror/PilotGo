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
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func MachineInfoHandler(c *gin.Context) {
	query := &service.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	depart := &service.Depart{}
	if c.ShouldBind(&depart) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	data, lens, err := service.MachineInfo(depart, query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	service.JsonPagination(c, data, int64(lens), query)
}

// 资源池返回接口
func FreeMachineSource(c *gin.Context) {
	machine := service.MachineNode{}
	query := &service.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	list, tx, res := machine.ReturnMachine(global.UncateloguedDepartId)
	total, err := service.CrudAll(query, tx, &res)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 返回数据开始拼装分页的json
	service.JsonPagination(c, list, total, query)
}

func MachineAllDataHandler(c *gin.Context) {
	datas, err := service.MachineAllData()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, datas, "获取所有的机器数据")
}

// 删除机器
func DeleteMachineHandler(c *gin.Context) {
	var deleteuuid service.DeleteUUID
	if c.Bind(&deleteuuid) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	machinelist := service.DeleteMachine(deleteuuid.Deluuid)

	if len(machinelist) != 0 {
		response.Fail(c, gin.H{"machinelist": machinelist}, "机器删除失败")
	} else {
		response.Success(c, nil, "机器删除成功!")
	}
}
