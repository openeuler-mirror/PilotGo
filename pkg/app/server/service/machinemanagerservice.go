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
 * Date: 2022-04-26 15:32:50
 * LastEditTime: 2022-04-27 09:37:48
 * Description: machinemanager service
 ******************************************************************************/
package service

import (
	"fmt"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
)

//返回全部的部门指针数组
func Returnptrchild(depart []model.DepartNode) (ptrchild []*model.MachineTreeNode, root model.MachineTreeNode) {
	departnode := make([]model.MachineTreeNode, 0)
	ptrchild = make([]*model.MachineTreeNode, 0)

	for _, value := range depart {
		if value.NodeLocate == 0 {
			root = model.MachineTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   0,
			}
		} else {
			departnode = append(departnode, model.MachineTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   value.PID,
			})
		}

	}
	ptrchild = append(ptrchild, &root)
	var a *model.MachineTreeNode
	for key := range departnode {
		a = &departnode[key]
		ptrchild = append(ptrchild, a)
	}
	return ptrchild, root
}

//生成部门树
func MakeTree(node *model.MachineTreeNode, ptrchild []*model.MachineTreeNode) {
	childs := findchild(node, ptrchild)
	for _, value := range childs {
		node.Children = append(node.Children, value)
		if IsChildExist(value, ptrchild) {
			MakeTree(value, ptrchild)
		}
	}
}

//返回节点的子节点
func findchild(node *model.MachineTreeNode, ptrchild []*model.MachineTreeNode) (ret []*model.MachineTreeNode) {
	for _, value := range ptrchild {
		if node.Id == value.Pid {
			ret = append(ret, value)
		}
	}
	return
}

//判断是否存在子节点
func IsChildExist(node *model.MachineTreeNode, ptrchild []*model.MachineTreeNode) bool {
	for _, child := range ptrchild {
		if node.Id == child.Pid {
			return true
		}
	}
	return false
}

func LoopTree(node *model.MachineTreeNode, ID int, res **model.MachineTreeNode) {
	if node.Children != nil {
		for _, value := range node.Children {
			if value.Id == ID {
				*res = value
			}

			LoopTree(value, ID, res)

		}

	}
}

func Deletedepartnode(DepartInfo []model.DepartNode, departid int) {
	needdelete := make([]int, 0)
	needdelete = append(needdelete, departid)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}

	for {
		if len(needdelete) == 0 {
			break
		}
		dao.Deletedepartdata(needdelete)
		str := fmt.Sprintf("%d", needdelete[0])
		needdelete = needdelete[1:]
		dao.Insertdepartlist(needdelete, str)
	}
}
