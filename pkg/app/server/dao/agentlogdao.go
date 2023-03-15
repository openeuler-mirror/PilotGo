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
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

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

func (p *AgentLogParent) LogAll(q *PaginationQ) (list *[]AgentLogParent, tx *gorm.DB) {
	list = &[]AgentLogParent{}
	tx = global.PILOTGO_DB.Order("created_at desc").Find(&list)
	return
}

func (p *AgentLog) AgentLog(q *PaginationQ, parentId int) (list *[]AgentLog, tx *gorm.DB) {
	list = &[]AgentLog{}
	tx = global.PILOTGO_DB.Order("ID desc").Where("log_parent_id=?", parentId).Find(list)
	return
}

// 删除agent日志
func DeleteLog(PLogIds int) error {
	var logparent AgentLogParent
	var log AgentLog

	err := global.PILOTGO_DB.Where("log_parent_id=?", PLogIds).Unscoped().Delete(log).Error
	if err != nil {
		logger.Error(err.Error())
	}
	err = global.PILOTGO_DB.Where("id=?", PLogIds).Unscoped().Delete(logparent).Error
	return err
}

// 存储父日志
func ParentAgentLog(PLog AgentLogParent) (int, error) {
	err := global.PILOTGO_DB.Save(&PLog).Error
	return PLog.ID, err
}

// 存储子日志
func AgentLogMessage(Log AgentLog) error {
	return global.PILOTGO_DB.Save(&Log).Error
}

// 查询子日志
func Id2AgentLog(id int) ([]AgentLog, error) {
	var Log []AgentLog
	err := global.PILOTGO_DB.Where("log_parent_id = ?", id).Find(&Log).Error
	return Log, err
}

// 修改父日志的操作状态
func UpdateParentAgentLog(PLogId int, status string) error {
	var ParentLog AgentLogParent
	return global.PILOTGO_DB.Model(&ParentLog).Where("id=?", PLogId).Update("status", status).Error
}
