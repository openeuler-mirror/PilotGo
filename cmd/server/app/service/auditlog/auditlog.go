/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package auditlog

import (
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

	ScriptAdd    = "自定义脚本/脚本创建"
	ScriptExec   = `自定义脚本/脚本运行`
	ScriptDelete = "自定义脚本/脚本删除"
	ScriptEdit   = "自定义脚本/脚本编辑"
)

type AuditLog = dao.AuditLog
type SubLog = dao.SubLog

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
