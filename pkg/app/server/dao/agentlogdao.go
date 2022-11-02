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
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

// 删除agent日志
func DeleteLog(PLogIds int) {
	var logparent model.AgentLogParent
	var log model.AgentLog

	global.PILOTGO_DB.Where("log_parent_id=?", PLogIds).Unscoped().Delete(log)
	global.PILOTGO_DB.Where("id=?", PLogIds).Unscoped().Delete(logparent)

}

// 存储父日志
func ParentAgentLog(PLog model.AgentLogParent) int {
	global.PILOTGO_DB.Save(&PLog)
	return PLog.ID
}

// 存储子日志
func AgentLog(Log model.AgentLog) {
	global.PILOTGO_DB.Save(&Log)
}

// 查询子日志
func Id2AgentLog(id int) []model.AgentLog {
	var Log []model.AgentLog
	global.PILOTGO_DB.Where("log_parent_id = ?", id).Find(&Log)
	return Log
}

// 修改父日志的操作状态
func UpdateParentAgentLog(PLogId int, status string) {
	var ParentLog model.AgentLogParent
	global.PILOTGO_DB.Model(&ParentLog).Where("id=?", PLogId).Update("status", status)

}
