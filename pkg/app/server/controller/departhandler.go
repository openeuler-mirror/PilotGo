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
 * Date: 2022-06-02 10:25:52
 * LastEditTime: 2022-06-08 16:16:10
 * Description: depart info handler
 ******************************************************************************/
package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 获取部门下所有机器列表
func MachineListHandler(c *gin.Context) {
	DepartId := c.Query("DepartId")
	DepId, err := strconv.Atoi(DepartId)
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	machinelist, err := service.MachineList(DepId)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(machinelist) == 0 {
		response.Success(c, []interface{}{}, "部门下所属机器获取成功")
	}
	response.Success(c, machinelist, "部门下所属机器获取成功")
}

func DepartHandler(c *gin.Context) {
	departID := c.Query("DepartID")
	tmp, err := strconv.Atoi(departID)
	if err != nil {
		response.Fail(c, nil, "部门ID有误")
		return
	}
	node, err := service.Dept(tmp)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, node, "获取当前部门及子部门信息")
}

func DepartInfoHandler(c *gin.Context) {
	departRoot, err := service.DepartInfo()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, departRoot, "获取全部的部门信息")
}

func AddDepartHandler(c *gin.Context) {
	newDepart := dao.AddDepart{}
	if err := c.Bind(&newDepart); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.AddDepart(&newDepart)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门信息入库成功")
}

func DeleteDepartDataHandler(c *gin.Context) {
	var DelDept dao.DeleteDepart
	if err := c.Bind(&DelDept); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.DeleteDepartData(&DelDept)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门删除成功")
}

func UpdateDepartHandler(c *gin.Context) {
	var new dao.NewDepart
	if err := c.Bind(&new); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.UpdateDepart(new.DepartID, new.DepartName)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门更新成功")
}

func ModifyMachineDepartHandler(c *gin.Context) {
	var M dao.MachineModifyDepart
	if err := c.Bind(&M); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.ModifyMachineDepart(M.MachineID, M.DepartID)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "机器部门修改成功")
}
