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
 * LastEditTime: 2023-06-28 15:58:39
 * Description: os scheduled task
 ******************************************************************************/
package dao

import (
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
)

type CrontabList struct {
	ID          int `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MachineUUID string `json:"uuid"`
	TaskName    string `json:"taskname"`
	Description string `json:"description"`
	CronSpec    string `json:"spec"`
	Command     string `json:"cmd"`
	Status      *bool  `json:"status"`
}

type CrontabUpdate struct {
	ID          int    `json:"id"`
	MachineUUID string `json:"uuid"`
	TaskName    string `json:"taskname"`
	Description string `json:"description"`
	CronSpec    string `json:"spec"`
	Command     string `json:"cmd"`
	Status      bool   `json:"status"`
}

type DelCrons struct {
	IDs         []int  `json:"ids"`
	MachineUUID string `json:"uuid"`
}

// 根据uuid获取所有机器
func CronListPaged(uuid string, offset, size int) (int64, []CrontabList, error) {
	var crontabLists []CrontabList
	var count int64
	err := mysqlmanager.MySQL().Model(CrontabList{}).Order("id desc").Where("machine_uuid = ?", uuid).Offset(offset).Limit(size).Find(&crontabLists).Offset(-1).Limit(-1).Count(&count).Error
	return count, crontabLists, err
}

// 任务名称是否存在
func IsTaskNameExist(name string) (bool, error) {
	var cron CrontabList
	err := mysqlmanager.MySQL().Where("task_name=?", name).Find(&cron).Error
	return cron.ID != 0, err
}

// 判断任务状态
func IsTaskStatus(id int, status bool) (bool, error) {
	var cron CrontabList
	err := mysqlmanager.MySQL().Where("id = ?", id).Find(&cron).Error
	return cron.Status == &status, err
}

// 新建定时任务
func NewCron(c CrontabList) (int, error) {
	err := mysqlmanager.MySQL().Save(&c).Error
	return c.ID, err
}

// 删除任务
func DeleteTask(id int) error {
	var cron CrontabList
	return mysqlmanager.MySQL().Where("id=?", id).Unscoped().Delete(cron).Error
}

// 更新任务
func UpdateTask(id int, c CrontabList) error {
	var cron CrontabList
	return mysqlmanager.MySQL().Model(&cron).Where("id=?", id).Updates(&c).Error
}

// 任务状态更新
func CronTaskStatus(id int, status bool) error {
	var cron CrontabList
	flag := !status
	UpdateCron := CrontabList{
		Status: &flag,
	}
	return mysqlmanager.MySQL().Model(&cron).Where("id=?", id).Updates(&UpdateCron).Error
}

// 根据任务id获取spec和command
func Id2CronInfo(id int) (spec, command string, err error) {
	var cron CrontabList
	err = mysqlmanager.MySQL().Where("id =?", id).Find(&cron).Error
	return cron.CronSpec, cron.Command, err
}
