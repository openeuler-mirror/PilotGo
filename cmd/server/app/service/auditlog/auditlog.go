/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package auditlog

import (
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
)

// 日志执行操作状态
const (
	StatusSuccess = "成功"
	StatusFail    = "失败"
)

// 日志记录归属模块
const (
	BatchCreate = "系统/创建批次"
	BatchDelete = "系统/删除批次"
	BatchEdit   = "系统/编辑批次"

	DepartAdd    = "系统/部门添加"
	DepartDelete = "系统/部门删除"
	DepartEdit   = "系统/部门更名"

	RPMInstall = "系统/软件包安装"
	RPMRemove  = "系统/软件包卸载"

	RoleChange = "角色管理/权限变更"

	ScriptAdd = "自定义脚本/创建脚本"
)

type AuditLog = dao.AuditLog
type SubLog = dao.SubLog

// 单机操作成功状态:是否成功，机器数量，成功率
const (
	ActionOK    = "1,1,1.00"
	ActionFalse = "0,1,0.00"
)

// 计算批量机器操作的状态：成功数，总数目，比率
func BatchActionStatus(StatusCodes []string) (status string) {
	var StatusOKCounts int
	for _, success := range StatusCodes {
		if success == strconv.Itoa(http.StatusOK) {
			StatusOKCounts++
		}
	}
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(StatusOKCounts)/float64(len(StatusCodes))), 64)
	rate := strconv.FormatFloat(num, 'f', 2, 64)
	status = strconv.Itoa(StatusOKCounts) + "," + strconv.Itoa(len(StatusCodes)) + "," + rate
	return
}

func Add(log *dao.AuditLog) (int, error) {
	return log.Record()
}
func AddSubLog(log *dao.SubLog) (int, error) {
	return log.Record()
}

// 修改日志的操作状态
func UpdateLog(logId int, status string) error {
	return dao.UpdateLogStatus(logId, status)
}

// 添加message信息
func UpdateSubLog(subLogId int, status, message string) error {
	return dao.UpdateSubLogMessage(subLogId, status, message)
}

// 分页查询
func GetAuditLogPaged(offset, size int) (int64, []AuditLog, error) {
	return dao.GetAuditLogPaged(offset, size)
}
