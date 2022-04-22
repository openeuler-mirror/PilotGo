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
 * Date: 2022-02-23 17:46:13
 * LastEditTime: 2022-03-25 03:04:30
 * Description: provide agent log manager functions.
 ******************************************************************************/
package model

import (
	"fmt"
	"time"

	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
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

const (
	RPMInstall     = "软件包安装"
	RPMRemove      = "软件包卸载"
	SysctlChange   = "修改内核参数"
	ServiceRestart = "重启服务"
	ServiceStop    = "关闭服务"
	ServiceStart   = "开启服务"
)

func (p *AgentLogParent) LogAll(q *utils.PaginationQ, Dnames []string) (list []AgentLogParent, total uint, err error) {
	list = []AgentLogParent{}
	lists := []AgentLogParent{}
	for _, name := range Dnames {
		mysqlmanager.DB.Order("created_at desc").Where("depart_name = ?", name).Find(&list)
		lists = append(lists, list...)
	}
	list, total, err = SliceAll(q, lists)
	return
}

func (p *AgentLog) AgentLog(q *utils.PaginationQ, parentId int) (list *[]AgentLog, total uint, err error) {
	list = &[]AgentLog{}
	tx := mysqlmanager.DB.Order("ID desc").Where("log_parent_id=?", parentId).Find(list)
	total, err = utils.CrudAll(q, tx, list)
	return
}
func SliceAll(p *utils.PaginationQ, data []AgentLogParent) ([]AgentLogParent, uint, error) {
	if p.Size < 1 {
		p.Size = 10
	}
	if p.CurrentPageNum < 1 {
		p.CurrentPageNum = 1
	}
	total := len(data)
	if total == 0 {
		p.TotalPage = 1
	}
	num := p.Size * (p.CurrentPageNum - 1)
	if num > uint(total) {
		return nil, uint(total), fmt.Errorf("页码超出")
	}
	if p.Size*p.CurrentPageNum > uint(total) {
		return data[num:], uint(total), nil
	} else {
		if p.Size*p.CurrentPageNum < num {
			return nil, uint(total), fmt.Errorf("读取错误")
		}
		if p.Size*p.CurrentPageNum == 0 {
			return data, uint(total), nil
		} else {
			return data[num : p.CurrentPageNum*p.Size], uint(total), nil
		}
	}
}
