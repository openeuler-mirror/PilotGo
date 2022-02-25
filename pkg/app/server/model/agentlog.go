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
 * LastEditTime: 2022-02-24 16:54:59
 * Description: provide agent log manager functions.
 ******************************************************************************/
package model

import "github.com/jinzhu/gorm"

type AgentHistory struct {
	gorm.Model
	Type       string `json:"type"`
	IP         string `json:"ip"`
	UserName   string `json:"username"`
	Status     string `json:"status"` // 成功、失败
	Rpm        string `json:"rpm"`
	Service    string `json:"service"`
	SysctlArgs string `json:"args"`
	SourceDisk string `json:"source"`
	MountPath  string `json:"mountpath"`
	FileType   string `json:"filetype"`
	AgentUser  string `json:"agentuser"`
	Password   string `json:"agentpass"`
	File       string `json:"file"`
	Permission string `json:"per"`
	Operation  string `json:"operation"`
	Message    string `json:"msg"`
}
