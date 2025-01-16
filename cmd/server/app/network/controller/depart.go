/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"strconv"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/depart"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/net"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
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

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "部门添加",
		Module:     auditlog.DepartAdd,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "部门添加:" + newDepart.DepartName,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	err = depart.AddDepart(&newDepart)
	if err != nil {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
		auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "部门添加失败："+err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")

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

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "部门删除",
		Module:     auditlog.DepartDelete,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "部门删除：" + DelDept.DepartName,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	err = depart.DeleteDepart(&DelDept)
	if err != nil {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
		auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "部门删除失败："+err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")
	response.Success(c, nil, "部门删除成功")
}

// 更改部门节点名字
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

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "修改部门名称",
		Module:     auditlog.DepartEdit,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "修改部门名称" + new.DepartName,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	err = depart.UpdateDepart(new.DepartID, new.DepartName)
	if err != nil {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
		auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "部门名称修改失败："+err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")

	response.Success(c, nil, "部门更新成功")
}
