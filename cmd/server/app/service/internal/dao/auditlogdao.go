/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	LogUUID    string `gorm:"not null;unique" json:"log_uuid"`
	ParentUUID string `gorm:"type:varchar(60)" json:"parent_uuid"`
	Module     string `gorm:"type:varchar(30);not null" json:"module"`
	Status     string `gorm:"type:varchar(30);not null" json:"status"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	Action     string `gorm:"not null" json:"action"`
	Message    string `gorm:"type:text" json:"message"`
	Isempty    int
}

// 存储日志
func (p *AuditLog) Record() error {
	return mysqlmanager.MySQL().Save(p).Error
}

// 修改日志的操作状态
func (p *AuditLog) UpdateStatus(status string) error {
	return mysqlmanager.MySQL().Model(&p).Where("log_uuid=?", p.LogUUID).Update("status", status).Error
}

func (p *AuditLog) UpdateMessage(message string) error {
	return mysqlmanager.MySQL().Model(&p).Where("log_uuid=?", p.LogUUID).Update("message", message).Error
}

// 分页查询
func GetAuditLogPaged(offset, size int) (int64, []AuditLog, error) {
	var count int64
	var auditlogs []AuditLog
	err := mysqlmanager.MySQL().Model(AuditLog{}).Order("id desc").Offset(offset).Limit(size).Find(&auditlogs).Offset(-1).Limit(-1).Count(&count).Error
	return count, auditlogs, err
}

// 查询子日志
func GetAuditLogById(logUUId string) ([]AuditLog, error) {
	var list []AuditLog
	err := mysqlmanager.MySQL().Order("created_at desc").Where("parent_uuid=?", logUUId).Find(&list).Error
	return list, err
}

// 查询父日志为空的记录
func GetParentLog(offset, size int) (int64, []AuditLog, error) {
	var count int64
	var auditlogs []AuditLog
	rows := mysqlmanager.MySQL().Raw("SELECT count(*) FROM audit_log WHERE parent_uuid = ''").Scan(&count)
	if rows.Error != nil {
		return 0, nil, rows.Error
	}
	sql := "select id,log_uuid,parent_uuid,module,status,user_id,action,message,exists(select 1 from audit_log t2 where t2.parent_uuid=t1.log_uuid) as isempty from audit_log t1 where t1.parent_uuid='' order by id desc limit ? offset ?;"
	rows = mysqlmanager.MySQL().Raw(sql, size, offset).Scan(&auditlogs)
	if rows.Error != nil {
		return count, nil, rows.Error
	}
	return count, auditlogs, nil
}

// 根据模块名字查询日志
func GetAuditLogByModule(name string) ([]AuditLog, error) {
	var Log []AuditLog
	err := mysqlmanager.MySQL().Where("module = ?", name).Find(&Log).Error
	return Log, err
}
