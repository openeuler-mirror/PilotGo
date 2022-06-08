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
 * Date: 2022-06-02 10:25:52
 * LastEditTime: 2022-06-02 16:16:10
 * Description: depart info struct
 ******************************************************************************/
package model

type DepartNode struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"`
	PID          int    `gorm:"type:int(100);not null" json:"pid"`
	ParentDepart string `gorm:"type:varchar(100);not null" json:"parentdepart"`
	Depart       string `gorm:"type:varchar(100);not null" json:"depart"`
	NodeLocate   int    `gorm:"type:int(100);not null" json:"nodelocate"`
	//根节点为0,普通节点为1
}

type DepartTreeNode struct {
	Label    string            `json:"label"`
	Id       int               `json:"id"`
	Pid      int               `json:"pid"`
	Children []*DepartTreeNode `json:"children"`
}

type NewDepart struct {
	DepartID   int    `json:"DepartID"`
	DepartName string `json:"DepartName"`
}

type MachineModifyDepart struct {
	MachineID string `json:"machineid"`
	DepartID  int    `json:"departid"`
}

type DeleteDepart struct {
	DepartID int `json:"DepartID"`
}

type AddDepart struct {
	ParentID     int    `json:"PID"`
	ParentDepart string `json:"ParentDepart"`
	DepartName   string `json:"Depart"`
}

type Depart struct {
	ID int `form:"DepartId"`
}
