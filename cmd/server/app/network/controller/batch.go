/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"fmt"
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/net"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 添加批次
func CreateBatchHandler(c *gin.Context) {
	params := &batch.CreateBatchParam{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, net.GetValidMsg(err, params))
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
		Module:     auditlog.ModuleBatch,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "创建批次",
	}
	auditlog.Add(log)
	params.Manager = u.Email

	if err := batch.CreateBatch(params); err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, err.Error())
		return
	}

	global.SendRemindMsg(
		global.ServerSendMsg,
		fmt.Sprintf("用户 %s 创建批次 %s", u.Username, params.Name),
	)

	response.Success(c, nil, "批次入库成功")
}

// 分页查询所有批次
func BatchInfoHandler(c *gin.Context) {
	p := &response.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := p.PageSize * (p.Page - 1)
	total, data, err := batch.GetBatchPaged(num, p.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), p)
}

// 删除批次
func DeleteBatchHandler(c *gin.Context) {
	batchdel := struct {
		BatchID []int `json:"BatchID"`
	}{}
	if err := c.Bind(&batchdel); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if len(batchdel.BatchID) == 0 {
		response.Fail(c, nil, "请输入删除批次ID")
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
		Module:     auditlog.ModuleBatch,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "删除批次",
	}
	auditlog.Add(log)

	if err := batch.DeleteBatch(batchdel.BatchID); err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	global.SendRemindMsg(
		global.ServerSendMsg,
		fmt.Sprintf("用户 %s 删除批次 %v", u.Username, batchdel.BatchID),
	)

	response.Success(c, nil, "批次删除成功")
}

// 更改批次
func UpdateBatchHandler(c *gin.Context) {
	batchinfo := struct {
		BatchId     int    `json:"BatchID"`
		BatchName   string `json:"BatchName"`
		Description string `json:"Description"`
	}{}
	if err := c.Bind(&batchinfo); err != nil {
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
		Module:     auditlog.ModuleBatch,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "编辑批次",
	}
	auditlog.Add(log)

	err = batch.UpdateBatch(batchinfo.BatchId, batchinfo.BatchName, batchinfo.Description)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, gin.H{"status": false}, "update batch failed: "+err.Error())
		return
	}

	global.SendRemindMsg(
		global.ServerSendMsg,
		fmt.Sprintf("用户 %s 更新批次 %s", u.Username, batchinfo.BatchName),
	)

	response.Success(c, nil, "批次修改成功")
}

// 查询某一个批次
func BatchMachineInfoHandler(c *gin.Context) {
	p := &response.PaginationQ{}
	err := c.ShouldBindQuery(p)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	Batchid := c.Query("ID")
	batchid, err := strconv.Atoi(Batchid)
	if err != nil {
		response.Fail(c, nil, "批次ID输入格式有误")
		return
	}

	num := p.PageSize * (p.Page - 1)
	total, data, err := batch.GetBatchMachines(num, p.PageSize, batchid)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), p)
}

// 一次性获取素有批次，供下拉列表选择
func SelectBatchHandler(c *gin.Context) {
	batch, err := batch.SelectBatch()
	if err != nil {
		response.Fail(c, nil, "获取批次信息错误"+err.Error())
		return
	}

	if len(batch) == 0 {
		response.Fail(c, nil, "未获取到批次信息")
		return
	}
	response.Success(c, gin.H{"data": batch}, "批次信息获取成功")
}
