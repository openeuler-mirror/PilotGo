/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: wanghao
 * Date: 2022-01-18 02:33:45
 * LastEditTime: 2022-03-04 00:09:13
 * Description: 批次属性列表
 ******************************************************************************/
package model

import (
	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
)

type Batch struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(100)" json:"description"`
	Manager     string `gorm:"type:varchar(100)" json:"manager"`
	Machinelist string `json:"machinelist"`
	Depart      string `gorm:"type:varchar(100)"`
	DepartName  string `gorm:"type:varchar(100)"`
}

type CreateBatch struct {
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Manager     string   `json:"Manager"`
	DepartName  []string `json:"DepartName"`
	DepartID    []int    `json:"DepartID"`
	Machines    []int    `json:"Machines"`
	DepartIDs   []int    `json:"deptids"`
}

type BatchUpdate struct {
	BatchId     int    `json:"BatchID"`
	BatchName   string `json:"BatchName"`
	Description string `json:"Description"`
}

type BatchDel struct {
	BatchID []int `json:"BatchID"`
}

func (b *Batch) ReturnBatch(q *PaginationQ) (list *[]Batch, tx *gorm.DB) {
	list = &[]Batch{}
	tx = global.PILOTGO_DB.Order("created_at desc").Find(&list)
	return
}
