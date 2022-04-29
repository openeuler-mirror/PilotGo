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
 * Date: 2021-04-29 09:08:08
 * LastEditTime: 2022-04-29 09:25:41
 * Description: 集群概览数据结构体
 ******************************************************************************/
package model

type ClusterInfo struct {
	AgentTotal  int `json:"total"`
	AgentStatus AgentStatus
}

type DepartMachineInfo struct {
	DepartName  string `json:"depart"`
	AgentStatus AgentStatus
}

type AgentStatus struct {
	Normal  int `json:"normal"`
	OffLine int `json:"offline"`
	Free    int `json:"free"`
}
