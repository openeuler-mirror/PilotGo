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
package cron

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
)

type CrontabUpdate = dao.CrontabUpdate
type CrontabList = dao.CrontabList
type DelCrons = dao.DelCrons

// 开启任务
func CronStart(uuid string, id int, spec string, command string) (interface{}, error) {
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		return nil, fmt.Errorf("server端获取uuid失败")
	}
	cron_start, Err, err := agent.CronStart(id, spec, command)
	if len(Err) != 0 || err != nil {
		return nil, fmt.Errorf("任务执行失败:%s", Err)
	}
	return cron_start, nil
}

// 暂停任务
func StopAndDel(uuid string, id int) (interface{}, error) {
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		return nil, fmt.Errorf("server端获取uuid失败")
	}
	cron_stop, err := agent.CronStopAndDel(id)
	if err != nil {
		return nil, fmt.Errorf("任务暂停失败:%s", err)
	}
	return cron_stop, nil
}

// 任务名称是否存在
func IsTaskNameExist(name string) (bool, error) {
	return dao.IsTaskNameExist(name)
}

// 新建定时任务
func NewCron(c CrontabList) (int, error) {
	return dao.NewCron(c)
}

// 根据任务id获取spec和command
func Id2CronInfo(id int) (spec, command string, err error) {
	return dao.Id2CronInfo(id)
}

// 删除任务
func DeleteTask(id int) error {
	return dao.DeleteTask(id)
}

// 更新任务
func UpdateTask(id int, c CrontabList) error {
	return dao.UpdateTask(id, c)
}

// 判断任务状态
func IsTaskStatus(id int, status bool) (bool, error) {
	return dao.IsTaskStatus(id, status)
}

// 任务状态更新
func CronTaskStatus(id int, status bool) error {
	return dao.CronTaskStatus(id, status)
}

// 根据uuid获取所有机器
func CronListPaged(uuid string, offset, size int) (int64, []CrontabList, error) {
	return dao.CronListPaged(uuid, offset, size)
}
