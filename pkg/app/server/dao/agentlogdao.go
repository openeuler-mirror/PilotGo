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
	"time"

	"gorm.io/gorm"
	"openeuler.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

type AuditLog struct {
	ID              uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	LogUUID         string `gorm:"not null;unique" json:"log_uuid"`
	ParentLogUUID   string `json:"parent_log_uuid"`
	AgentUUID       string `json:"agent_uuid"`
	Module          string `gorm:"type:varchar(30);not null" json:"module"`
	Status          string `gorm:"type:varchar(30);not null" json:"status"`
	UserName        string `gorm:"not null" json:"userName"`
	DepartName      string `gorm:"not null" json:"departName"`
	Email           string `gorm:"not null" json:"email"`
	Action          string `gorm:"not null" json:"action"`
	Message         string `json:"message"`
	IP              string `json:"ip"`
	StatusCode      int    `json:"code"`
	OperationObject string `json:"object"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	// OperatorID    uint   `gorm:"not null" json:"operator_id"`
}

// 存储日志
func (p *AuditLog) Record() error {
	return mysqlmanager.MySQL().Save(p).Error
}

// 修改日志的操作状态
func (p *AuditLog) UpdateStatus(status string) error {
	// TODO:
	return mysqlmanager.MySQL().Model(&p).Where("log_uuid=?", p.LogUUID).Update("status", status).Error
}

// 查询所有父日志
func GetAuditLog() (list *[]AuditLog, tx *gorm.DB, err error) {
	list = &[]AuditLog{}
	tx = mysqlmanager.MySQL().Order("created_at desc").Where("parent_log_uuid=?", "").Find(&list)
	err = tx.Error
	return
}

// 根据父UUid查询日志
func GetAuditLogByParentId(parentUUId string) (list *[]AuditLog, tx *gorm.DB, err error) {
	list = &[]AuditLog{}
	tx = mysqlmanager.MySQL().Order("created_at desc").Where("parent_log_uuid=?", parentUUId).Find(&list)
	err = tx.Error
	return
}

// 查询子日志
func GetAuditLogById(logUUId string) (AuditLog, error) {
	var Log AuditLog
	err := mysqlmanager.MySQL().Where("log_uuid = ?", logUUId).Find(&Log).Error
	return Log, err
}

// 根据模块名字查询日志
func GetAuditLogByModule(name string) ([]AuditLog, error) {
	var Log []AuditLog
	err := mysqlmanager.MySQL().Where("modulename = ?", name).Find(&Log).Error
	return Log, err
}

type AgentLogParent struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UserName   string    `json:"userName"`
	DepartName string    `json:"departName"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
}
type AgentLog struct {
	ID              int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	LogParentID     int    `gorm:"index" json:"logparent_id"`
	IP              string `json:"ip"`
	StatusCode      int    `json:"code"`
	OperationObject string `json:"object"`
	Action          string `json:"action"`
	Message         string `json:"message"`
}
type AgentLogDel struct {
	IDs []int `json:"ids"`
}

func (p *AgentLogParent) LogAll() (list *[]AgentLogParent, tx *gorm.DB) {
	list = &[]AgentLogParent{}
	tx = mysqlmanager.MySQL().Order("created_at desc").Find(&list)
	return
}

func (p *AgentLog) AgentLog(parentId int) (list *[]AgentLog, tx *gorm.DB) {
	list = &[]AgentLog{}
	tx = mysqlmanager.MySQL().Order("ID desc").Where("log_parent_id=?", parentId).Find(list)
	return
}

// 删除agent日志
func DeleteLog(PLogIds int) error {
	var logparent AgentLogParent
	var log AgentLog

	err := mysqlmanager.MySQL().Where("log_parent_id=?", PLogIds).Unscoped().Delete(log).Error
	if err != nil {
		logger.Error(err.Error())
	}
	return mysqlmanager.MySQL().Where("id=?", PLogIds).Unscoped().Delete(logparent).Error
}

// 存储父日志
func ParentAgentLog(PLog AgentLogParent) (int, error) {
	err := mysqlmanager.MySQL().Save(&PLog).Error
	return PLog.ID, err
}

// 存储子日志
func AgentLogMessage(Log AgentLog) error {
	return mysqlmanager.MySQL().Save(&Log).Error
}

// 查询子日志
func Id2AgentLog(id int) ([]AgentLog, error) {
	var Log []AgentLog
	err := mysqlmanager.MySQL().Where("log_parent_id = ?", id).Find(&Log).Error
	return Log, err
}

// 修改父日志的操作状态
func UpdateParentAgentLog(PLogId int, status string) error {
	var ParentLog AgentLogParent
	return mysqlmanager.MySQL().Model(&ParentLog).Where("id=?", PLogId).Update("status", status).Error
}
