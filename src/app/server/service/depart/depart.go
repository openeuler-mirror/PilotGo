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
 * Description: depart info service
 ******************************************************************************/
package depart

import (
	"errors"

	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type DepartNode = dao.DepartNode

type DeleteDeparts struct {
	DepartID int `json:"DepartID"`
}

type AddDepartNode struct {
	ParentID   int    `json:"PID"`
	DepartName string `json:"Depart" binding:"required" msg:"部门名称不能为空"`
}

type DepartTreeNode struct {
	Label    string            `json:"label"`
	Id       int               `json:"id"`
	Pid      int               `json:"pid"`
	Children []*DepartTreeNode `json:"children"`
}

// 将数据库DepartNode表数据转换成DepartTreeNode类型指针，并找到根节点
func Returnptrchild(depart []dao.DepartNode) (ptrchild []*DepartTreeNode, deptRoot DepartTreeNode) {
	ptrchild = make([]*DepartTreeNode, 0)
	//将DepartNode表数据转换成*DepartTreeNode类型
	for _, value := range depart {
		if value.NodeLocate == 0 {
			deptRoot = DepartTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   0,
			}
		} else {
			ptrchild = append(ptrchild, &DepartTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   value.PID,
			})
		}
	}
	//在第一个位置添加根节点
	ptrchild = append([]*DepartTreeNode{&deptRoot}, ptrchild...)
	return ptrchild, deptRoot
}

// 生成部门树
func MakeTree(node *DepartTreeNode, ptrchild []*DepartTreeNode) {
	//查找本节点的所有子节点
	childs := findchild(node, ptrchild)
	for _, value := range childs {
		//添加到本节点的Children字段中
		node.Children = append(node.Children, value)
		//检查value节点是否存在子节点
		if IsChildExist(value, ptrchild) {
			MakeTree(value, ptrchild)
		}
	}
}

// 返回节点的子节点
func findchild(node *DepartTreeNode, ptrchild []*DepartTreeNode) (ret []*DepartTreeNode) {
	for _, value := range ptrchild {
		if node.Id == value.Pid {
			ret = append(ret, value)
		}
	}
	return
}

/*
// 返回节点的所有子节点的ID
func AllChild(node *DepartTreeNode) (dts []int) {
	if node.Children != nil {
		for _, value := range node.Children {
			dts = append(dts, value.Id)
			dts = append(dts, AllChild(value)...)
		}
	}
	return dts
}*/

// 判断是否存在子节点
func IsChildExist(node *DepartTreeNode, ptrchild []*DepartTreeNode) bool {
	for _, child := range ptrchild {
		if node.Id == child.Pid {
			return true
		}
	}
	return false
}

func LoopTree(node *DepartTreeNode, ID int, res **DepartTreeNode) {
	if node.Children != nil {
		for _, value := range node.Children {
			if value.Id == ID {
				*res = value
			}
			LoopTree(value, ID, res)
		}
	}
}

// 根据部门id在部门树中查询部门节点
func Dept(tmp int, root *DepartTreeNode) (*DepartTreeNode, error) {
	if root == nil {
		return nil, errors.New("部门根节点为空")
	}
	var d *DepartTreeNode
	if root.Id != tmp {
		LoopTree(root, tmp, &d)
		root = d
	}
	return root, nil
}

// 生成部门树
func DepartInfo() (*DepartTreeNode, error) {
	depart, err := dao.GetAllDepart()
	if err != nil {
		return nil, err
	}
	if len(depart) == 0 {
		return nil, errors.New("当前无部门节点")
	}
	//返回所有根节点喝孩子节点
	ptrchild, departRoot := Returnptrchild(depart)
	//构造树
	MakeTree(&departRoot, ptrchild)
	return &departRoot, nil
}

// 添加部门
func AddDepart(newDepart *AddDepartNode) error {
	pid := newDepart.ParentID
	depart := newDepart.DepartName
	departNode := dao.DepartNode{
		PID:    pid,
		Depart: depart,
	}
	if pid == 0 {
		temp, err := dao.IsRootExist()
		if err != nil {
			return err
		}
		if temp {
			return errors.New("已存在根节点,即组织名称")
		} else {
			departNode.NodeLocate = global.Departroot
			if departNode.Add() != nil {
				return errors.New("部门节点添加失败")
			}
		}
	} else {
		parent, err := dao.GetDepartById(pid)
		if err != nil {
			return err
		}
		if parent.ID == 0 {
			return errors.New("部门PID有误,数据库中不存在该部门PID")
		}

		temp, err := dao.IsDepartNodeExist(pid, depart)
		if err != nil {
			return err
		}
		if temp {
			return errors.New("该部门节点已存在")
		}

		departNode.NodeLocate = global.DepartUnroot
		if departNode.Add() != nil {
			return errors.New("部门节点添加失败")
		}
	}
	return nil
}

// 删除部门
func DeleteDepart(DelDept *DeleteDeparts) error {
	//查看该部门是否存在
	temp, err := dao.GetDepartById(DelDept.DepartID)
	if err != nil {
		return err
	}
	if temp.ID == 0 {
		return errors.New("不存在该部门")
	}
	//先查询所有子部门，然后查询所有子部门机器和用户，若不为空则不可删除，若为空直接删除
	//查询所有子节点
	var departs []int
	ReturnSpecifiedDepart(DelDept.DepartID, &departs)

	//将要删除的节点添加到数组中
	departs = append([]int{DelDept.DepartID}, departs...)

	//查询与此部门相关的机器和用户
	var users []dao.User
	nodes, err := dao.MachineList(departs)
	if err != nil {
		return err
	}
	for _, depart := range departs {
		us, err := dao.GetUserBypid(depart)
		if err != nil {
			return err
		}
		users = append(users, us...)
	}
	if len(nodes) > 0 {
		return errors.New("该部门有机器不允许删除")
	}
	if len(users) > 0 {
		return errors.New("该部门有用户不允许删除")
	}
	if len(nodes) == 0 && len(users) == 0 {
		return dao.Deletedepartdata(departs)
	}
	return nil
}

// 更改部门名字
func UpdateDepart(DepartID int, DepartName string) error {
	return dao.UpdateDepart(DepartID, DepartName)
}

// 获取部门下所有机器列表
func MachineList(DepId int) ([]dao.Res, error) {
	var departId []int
	ReturnSpecifiedDepart(DepId, &departId)
	departId = append(departId, DepId)
	machinelist1, err := dao.MachineList(departId)
	return machinelist1, err
}

// 返回所有子部门id
func ReturnSpecifiedDepart(id int, res *[]int) {
	temp, err := dao.SubDepartId(id)
	if err != nil {
		logger.Error(err.Error())
	}
	if len(temp) == 0 {
		return
	}
	for _, value := range temp {
		*res = append(*res, value)
		ReturnSpecifiedDepart(value, res)
	}
}

// 创建公司组织
func CreateOrganization() error {
	return dao.CreateOrganization()
}
