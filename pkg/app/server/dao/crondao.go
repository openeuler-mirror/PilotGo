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
 * Date: 2022-05-23 10:25:52
 * LastEditTime: 2022-05-23 15:16:10
 * Description: os scheduled task
 ******************************************************************************/
package dao

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

// 任务名称是否存在
func IsTaskNameExist(name string) bool {
	var cron model.CrontabList
	mysqlmanager.DB.Where("task_name=?", name).Find(&cron)
	return cron.ID != 0
}

// 判断任务状态
func IsTaskStatus(id int, status bool) bool {
	var cron model.CrontabList
	mysqlmanager.DB.Where("id = ?", id).Find(&cron)
	return cron.Status == &status
}

// 新建定时任务
func NewCron(c model.CrontabList) (id int) {
	mysqlmanager.DB.Save(&c)
	return c.ID
}

// 删除任务
func DeleteTask(id int) {
	var cron model.CrontabList
	mysqlmanager.DB.Where("id=?", id).Unscoped().Delete(cron)
}

// 更新任务
func UpdateTask(id int, c model.CrontabList) {
	var cron model.CrontabList
	mysqlmanager.DB.Model(&cron).Where("id=?", id).Update(&c)
}

// 任务状态更新
func CronTaskStatus(id int, status bool) {
	var cron model.CrontabList
	flag := !status
	UpdateCron := model.CrontabList{
		Status: &flag,
	}
	mysqlmanager.DB.Model(&cron).Where("id=?", id).Update(&UpdateCron)
}

// 根据任务id获取spec和command
func Id2CronInfo(id int) (spec, command string) {
	var cron model.CrontabList
	mysqlmanager.DB.Where("id =?", id).Find(&cron)
	return cron.CronSpec, cron.Command
}
