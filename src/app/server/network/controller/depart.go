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

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/depart"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/utils/message/net"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NewDepart struct {
	DepartID   int    `json:"DepartID" binding:"required" msg:"部门id不能为空"`
	DepartName string `json:"DepartName" binding:"required" msg:"部门名称不能为空"`
}

// 获取部门下所有机器列表，若想获取根节点部门，传入departid=1
func MachineListHandler(c *gin.Context) {
	DepartId := c.Query("DepartId")
	DepId, err := strconv.Atoi(DepartId)
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	machinelist, err := depart.MachineList(DepId)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(machinelist) == 0 {
		response.Success(c, []interface{}{}, "部门下所属机器获取成功")
	}
	response.Success(c, machinelist, "部门下所属机器获取成功")
}

// 根据部门id在部门树中查询部门节点
func DepartHandler(c *gin.Context) {
	departID := c.Query("DepartID")
	tmp, err := strconv.Atoi(departID)
	if err != nil {
		response.Fail(c, nil, "部门ID有误")
		return
	}
	//构造部门树
	departRoot, err := depart.DepartInfo()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	//查询节点
	node, err := depart.Dept(tmp, departRoot)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, node, "获取当前部门及子部门信息")
}

// 生成部门树
func DepartInfoHandler(c *gin.Context) {
	departRoot, err := depart.DepartInfo()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, departRoot, "获取全部的部门信息")
}

// 添加部门节点
func AddDepartHandler(c *gin.Context) {
	newDepart := depart.AddDepartNode{}
	if err := c.Bind(&newDepart); err != nil {
		response.Fail(c, nil, net.GetValidMsg(err, &newDepart))
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleDepart,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "添加部门信息",
	}
	auditlog.Add(log)

	err = depart.AddDepart(&newDepart)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门信息入库成功")
}

// 删除部门节点
func DeleteDepartDataHandler(c *gin.Context) {
	var DelDept depart.DeleteDeparts
	if err := c.Bind(&DelDept); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleDepart,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "删除部门信息",
	}
	auditlog.Add(log)

	err = depart.DeleteDepart(&DelDept)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门删除成功")
}

// 更新部门节点名字
func UpdateDepartHandler(c *gin.Context) {
	var new NewDepart
	if err := c.Bind(&new); err != nil {
		response.Fail(c, nil, net.GetValidMsg(err, &new))
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModuleDepart,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "修改部门信息",
	}
	auditlog.Add(log)

	err = depart.UpdateDepart(new.DepartID, new.DepartName)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "部门更新成功")
}
