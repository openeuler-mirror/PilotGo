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
package model

import (
	"time"

	"gorm.io/gorm"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
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
func (c *CrontabList) CronList(uuid string, q *PaginationQ) (list *[]CrontabList, tx *gorm.DB) {
	list = &[]CrontabList{}
	tx = global.PILOTGO_DB.Order("created_at desc").Where("machine_uuid = ?", uuid).Find(&list)
	return
}
