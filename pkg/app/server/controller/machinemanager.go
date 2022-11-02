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
	"net/http"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func MachineInfo(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	depart := &model.Depart{}
	if c.ShouldBind(depart) != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	var TheDeptAndSubDeptIds []int
	service.ReturnSpecifiedDepart(depart.ID, &TheDeptAndSubDeptIds)
	TheDeptAndSubDeptIds = append(TheDeptAndSubDeptIds, depart.ID)
	machinelist := dao.MachineList(TheDeptAndSubDeptIds)

	lens := len(machinelist)
	data, err := service.DataPaging(query, machinelist, lens)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	service.JsonPagination(c, data, int64(lens), query)
}

// 资源池返回接口
func FreeMachineSource(c *gin.Context) {
	machine := model.MachineNode{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	list, tx, res := machine.ReturnMachine(query, global.UncateloguedDepartId)
	total, err := service.CrudAll(query, tx, &res)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 返回数据开始拼装分页的json
	service.JsonPagination(c, list, total, query)
}

func MachineAllData(c *gin.Context) {
	AllData := dao.MachineAllData()
	datas := make([]map[string]string, 0)
	for _, data := range AllData {
		datas = append(datas, map[string]string{"uuid": data.UUID, "ip_dept": data.IP + "-" + data.Departname, "ip": data.IP})
	}
	response.JSON(c, http.StatusOK, http.StatusOK, datas, "获取所有的机器数据")
}

// 删除机器
func DeleteMachine(c *gin.Context) {
	var deleteuuid model.DeleteUUID
	c.Bind(&deleteuuid)
	machinelist := make(map[string]string)
	for _, machinedeluuid := range deleteuuid.Deluuid {
		if err := dao.DeleteMachine(machinedeluuid); err != nil {
			machinelist[machinedeluuid] = err.Error()
		}
	}
	if len(machinelist) != 0 {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"machinelist": machinelist}, "机器删除失败")
	} else {
		response.Response(c, http.StatusOK, http.StatusOK, nil, "机器删除成功!")
	}
}
