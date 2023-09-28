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

	"gitee.com/PilotGo/PilotGo/app/server/agentmanager"
)

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
