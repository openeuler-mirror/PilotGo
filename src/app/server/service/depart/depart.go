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
	"fmt"

	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
)

type DepartNode = dao.DepartNode

type DeleteDepart struct {
	DepartID int `json:"DepartID"`
}

type AddDepart struct {
	ParentID   int    `json:"PID"`
	DepartName string `json:"Depart" binding:"required" msg:"部门名称不能为空"`
}

type MachineModifyDepart struct {
	MachineID string `json:"machineid"`
	DepartID  int    `json:"departid"`
}

type NewDepart struct {
	DepartID   int    `json:"DepartID" binding:"required" msg:"部门id不能为空"`
	DepartName string `json:"DepartName" binding:"required" msg:"部门名称不能为空"`
}

type DepartTreeNode struct {
	Label    string            `json:"label"`
	Id       int               `json:"id"`
	Pid      int               `json:"pid"`
	Children []*DepartTreeNode `json:"children"`
}

// 返回全部的部门指针数组
func Returnptrchild(depart []dao.DepartNode) (ptrchild []*DepartTreeNode, deptRoot DepartTreeNode) {
	departnode := make([]DepartTreeNode, 0)
	ptrchild = make([]*DepartTreeNode, 0)

	for _, value := range depart {
		if value.NodeLocate == 0 {
			deptRoot = DepartTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   0,
			}
		} else {
			departnode = append(departnode, DepartTreeNode{
				Label: value.Depart,
				Id:    value.ID,
				Pid:   value.PID,
			})
		}

	}
	ptrchild = append(ptrchild, &deptRoot)
	var a *DepartTreeNode
	for key := range departnode {
		a = &departnode[key]
		ptrchild = append(ptrchild, a)
	}
	return ptrchild, deptRoot
}

// 生成部门树
func MakeTree(node *DepartTreeNode, ptrchild []*DepartTreeNode) {
	childs := findchild(node, ptrchild)
	for _, value := range childs {
		node.Children = append(node.Children, value)
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

func DeleteDepartNode(DepartInfo []dao.DepartNode, departid int) {
	needdelete := make([]int, 0)
	needdelete = append(needdelete, departid)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}

	for {
		if len(needdelete) == 0 {
			break
		}
		err := dao.Deletedepartdata(needdelete)
		if err != nil {
			logger.Error(err.Error())
		}
		str := fmt.Sprintf("%d", needdelete[0])
		needdelete = needdelete[1:]
		dao.Insertdepartlist(needdelete, str)
	}
}

// 获取部门下所有机器列表
func MachineList(DepId int) ([]dao.Res, error) {
	var departId []int
	common.ReturnSpecifiedDepart(DepId, &departId)
	departId = append(departId, DepId)
	machinelist1, err := dao.MachineList(departId)
	if err != nil {
		return machinelist1, err
	}
	return machinelist1, nil
}

func Dept(tmp int) (*DepartTreeNode, error) {
	depart, err := dao.DepartStore()
	if err != nil {
		return nil, err
	}
	ptrchild, departRoot := Returnptrchild(depart)
	MakeTree(&departRoot, ptrchild)
	node := &departRoot
	var d *DepartTreeNode
	if node.Id != tmp {
		LoopTree(node, tmp, &d)
		node = d
	}
	if node == nil {
		return nil, errors.New("部门ID有误")
	}
	return node, nil
}

func DepartInfo() (*DepartTreeNode, error) {
	depart, err := dao.DepartStore()
	if err != nil {
		return nil, err
	}
	if len(depart) == 0 {
		return nil, errors.New("当前无部门节点")
	}
	ptrchild, departRoot := Returnptrchild(depart)
	MakeTree(&departRoot, ptrchild)
	return &departRoot, nil
}

func AddDepartMethod(newDepart *AddDepart) error {
	pid := newDepart.ParentID
	depart := newDepart.DepartName
	parent, err := dao.GetDepartById(pid)
	if err != nil {
		return err
	}
	if parent.ID == 0 {
		return errors.New("部门PID有误,数据库中不存在该部门PID")
	}

	departNode := dao.DepartNode{
		PID:    pid,
		Depart: depart,
	}
	temp, err := dao.IsDepartNodeExist(pid, depart)
	if err != nil {
		return err
	}
	if temp {
		return errors.New("该部门节点已存在")
	}

	if parent.ID == 1 && parent.PID == 0 {
		temp, err := dao.IsRootExist()
		if err != nil {
			return err
		}
		if temp {
			return errors.New("已存在根节点,即组织名称")
		} else {
			departNode.NodeLocate = global.Departroot
			if dao.AddDepartMessage(mysqlmanager.MySQL(), &departNode) != nil {
				return errors.New("部门节点添加失败")
			}
		}
	} else {
		departNode.NodeLocate = global.DepartUnroot
		if dao.AddDepartMessage(mysqlmanager.MySQL(), &departNode) != nil {
			return errors.New("部门节点添加失败")
		}
	}
	return nil
}

func DeleteDepartData(DelDept *DeleteDepart) error {
	temp, err := dao.IsDepartIDExist(DelDept.DepartID)
	if err != nil {
		return err
	}
	if !temp {
		return errors.New("不存在该部门")
	}
	macli, err := dao.MachineStore(DelDept.DepartID)
	if err != nil {
		return err
	}
	if len(macli) > 0 {
		return errors.New("该部门有机器不允许删除")
	}
	departs, err := dao.SubDepartId(DelDept.DepartID)
	if err != nil {
		return err
	}
	for _, depart := range departs {
		machine, err := dao.MachineStore(depart)
		if err != nil {
			return err
		}
		for _, m := range machine {
			err := dao.UpdateMachineDepartState(m.MachineUUID, global.UncateloguedDepartId)
			if err != nil {
				return err
			}
		}
	}

	DepartInfo, err := dao.Pid2Depart(DelDept.DepartID)
	if err != nil {
		return err
	}
	DeleteDepartNode(DepartInfo, DelDept.DepartID)
	err = dao.DelUserByDeptId(DelDept.DepartID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDepart(DepartID int, DepartName string) error {
	err := dao.UpdateDepart(DepartID, DepartName)
	if err != nil {
		return err
	}
	err = dao.UpdateParentDepart(DepartID, DepartName)
	return err
}

func ModifyMachineDepart(MachineID string, DepartID int) error {
	//查询部门节点是否存在
	is, err := dao.IsDepartIDExist(DepartID)
	if err != nil {
		return err
	}
	if !is {
		return errors.New("此部门不存在")
	}
	ResIds := utils.String2Int(MachineID)
	for _, id := range ResIds {
		machine, err := dao.MachineInfo(id)
		if err != nil {
			return err
		}
		err = dao.UpdateMachineDepartState(machine.MachineUUID, DepartID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建公司组织
func CreateOrganization() error {
	return dao.CreateOrganization()
}
