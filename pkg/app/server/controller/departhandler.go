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
 * LastEditTime: 2022-06-08 16:16:10
 * Description: depart info handler
 ******************************************************************************/
package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func Dept(c *gin.Context) {
	departID := c.Query("DepartID")
	tmp, err := strconv.Atoi(departID)
	if err != nil {
		response.Fail(c, nil, "部门ID有误")
		return
	}

	depart := dao.DepartStore()
	ptrchild, departRoot := service.Returnptrchild(depart)
	service.MakeTree(&departRoot, ptrchild)

	node := &departRoot
	var d *model.DepartTreeNode
	if node.Id != tmp {
		service.LoopTree(node, tmp, &d)
		node = d
	}
	if node == nil {
		response.Fail(c, nil, "部门ID有误")
		return
	}
	response.JSON(c, http.StatusOK, http.StatusOK, node, "获取当前部门及子部门信息")
}

func DepartInfo(c *gin.Context) {
	depart := dao.DepartStore()
	if len(depart) == 0 {
		response.Fail(c, nil, "当前无部门节点")
		return
	}
	ptrchild, departRoot := service.Returnptrchild(depart)
	service.MakeTree(&departRoot, ptrchild)
	response.JSON(c, http.StatusOK, http.StatusOK, departRoot, "获取全部的部门信息")
}

func AddDepart(c *gin.Context) {
	newDepart := model.AddDepart{}
	c.Bind(&newDepart)
	pid := newDepart.ParentID
	parentDepart := newDepart.ParentDepart
	depart := newDepart.DepartName

	if !dao.IsDepartIDExist(pid) {
		response.Fail(c, nil, "部门PID有误,数据库中不存在该部门PID")
		return
	}

	departNode := model.DepartNode{
		PID:          pid,
		ParentDepart: parentDepart,
		Depart:       depart,
	}
	if dao.IsDepartNodeExist(parentDepart, depart) {
		response.Fail(c, nil, "该部门节点已存在")
		return
	}
	if len(parentDepart) != 0 && !dao.IsParentDepartExist(parentDepart) {
		response.Fail(c, nil, "该部门上级部门不存在")
		return
	}
	if len(depart) == 0 {
		response.Fail(c, nil, "部门节点不能为空")
		return
	} else if len(parentDepart) == 0 {
		if dao.IsRootExist() {
			response.Fail(c, nil, "已存在根节点,即组织名称")
			return
		} else {
			departNode.NodeLocate = global.Departroot
			if dao.AddDepart(global.PILOTGO_DB, &departNode) != nil {
				response.Fail(c, nil, "部门节点添加失败")
				return
			}
		}
	} else {
		departNode.NodeLocate = global.DepartUnroot
		if dao.AddDepart(global.PILOTGO_DB, &departNode) != nil {
			response.Fail(c, nil, "部门节点添加失败")
			return
		}
	}
	response.Success(c, nil, "部门信息入库成功")
}

func DeleteDepartData(c *gin.Context) {
	var DelDept model.DeleteDepart
	c.Bind(&DelDept)

	if !dao.IsDepartIDExist(DelDept.DepartID) {
		response.Fail(c, nil, "不存在该机器")
		return
	}

	for _, mac := range dao.MachineStore(DelDept.DepartID) {
		dao.ModifyMachineDepart2(mac.ID, global.UncateloguedDepartId)
	}
	for _, depart := range dao.SubDepartId(DelDept.DepartID) {
		machine := dao.MachineStore(depart)
		for _, m := range machine {
			dao.ModifyMachineDepart2(m.ID, global.UncateloguedDepartId)
		}
	}

	DepartInfo := dao.Pid2Depart(DelDept.DepartID)
	service.DeleteDepartNode(DepartInfo, DelDept.DepartID)

	dao.DelUser(DelDept.DepartID)

	response.Success(c, nil, "部门删除成功")
}

func UpdateDepart(c *gin.Context) {
	var new model.NewDepart
	c.Bind(&new)

	dao.UpdateDepart(new.DepartID, new.DepartName)
	dao.UpdateParentDepart(new.DepartID, new.DepartName)
	response.Success(c, nil, "部门更新成功")
}

func ModifyMachineDepart(c *gin.Context) {
	var M model.MachineModifyDepart
	c.Bind(&M)

	Ids := strings.Split(M.MachineID, ",")
	ResIds := make([]int, len(Ids))
	for index, val := range Ids {
		ResIds[index], _ = strconv.Atoi(val)
	}

	for _, id := range ResIds {
		dao.ModifyMachineDepart(id, M.DepartID)
	}
	response.Success(c, nil, "机器部门修改成功")
}
