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
 * Date: 2021-04-28 13:08:08
 * LastEditTime: 2022-04-28 14:25:41
 * Description: agent操作日志相关数据获取
 ******************************************************************************/
package dao

import (
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
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
	err := mysqlmanager.MySQL().Model(AuditLog{}).Order("id desc").Where("parent_uuid=?", "").Offset(offset).Limit(size).Find(&auditlogs).Offset(-1).Limit(-1).Count(&count).Error
	return count, auditlogs, err
}

// 根据模块名字查询日志
func GetAuditLogByModule(name string) ([]AuditLog, error) {
	var Log []AuditLog
	err := mysqlmanager.MySQL().Where("module = ?", name).Find(&Log).Error
	return Log, err
}
