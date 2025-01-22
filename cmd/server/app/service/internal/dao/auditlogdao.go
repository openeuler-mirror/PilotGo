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
)

type AuditLog struct {
	ID         int      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Action     string   `json:"action"`     //执行动作
	Status     string   `json:"status"`     //整个任务的状态
	Module     string   `json:"module"`     //所属模块
	User       string   `json:"user"`       //创建用户
	Batches    string   `json:"batches"`    //处理批次
	CreateTime string   `json:"createTime"` //发生时间
	LogChild   []SubLog `gorm:"foreignKey:LogId;references:ID" json:"SubLog"`
}

type SubLog struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	LogId        int    `gorm:"index" json:"logId"`
	ActionObject string `json:"action"`     //执行动作+对象
	UpdateTime   string `json:"updateTime"` //执行结果更新时间
	Status       string `json:"status"`     //执行状态
	Message      string `json:"message"`    //日志详情
}

// 存储日志
func (p *AuditLog) Record() (int, error) {
	err := mysqlmanager.MySQL().Save(p).Error
	return p.ID, err
}

func (sl *SubLog) Record() (int, error) {
	err := mysqlmanager.MySQL().Save(sl).Error
	return sl.ID, err
}

// 修改日志的操作状态
func UpdateLogStatus(logId int, status string) error {
	return mysqlmanager.MySQL().Model(&AuditLog{}).Where("id=?", logId).Update("status", status).Error
}
func UpdateSubLogMessage(subLogId int, status string, message string) error {
	return mysqlmanager.MySQL().Model(&SubLog{}).Where("id=?", subLogId).
		Updates(SubLog{Status: status, Message: message}).Error
}

// 分页查询
func GetAuditLogPaged(offset, size int) (int64, []AuditLog, error) {
	var count int64
	var auditlogs []AuditLog
	err := mysqlmanager.MySQL().Model(AuditLog{}).Order("id desc").Preload("LogChild").Offset(offset).Limit(size).Find(&auditlogs).Offset(-1).Limit(-1).Count(&count).Error
	return count, auditlogs, err
}
func GetSubLogStatus(logId int) bool {
	var count int64
	mysqlmanager.MySQL().Model(SubLog{}).Select("status").Where("log_id = ? AND status = ?", logId, "失败").Count(&count)
	return count > 0
}
