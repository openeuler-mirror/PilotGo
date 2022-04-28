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
 * Date: 2022-02-18 13:03:16
 * LastEditTime: 2022-04-27 09:58:35
 * Description: provide machine manager functions.
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

const (
	UncateloguedDepartId = 1 // 新注册机器添加到部门根节点
	Departroot           = 0
	DepartUnroot         = 1
)

func AddDepart(c *gin.Context) {
	pid := c.Query("PID")
	parentDepart := c.Query("ParentDepart")
	depart := c.Query("Depart")
	tmp, err := strconv.Atoi(pid)
	if len(pid) != 0 && err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"pid识别失败")
		return
	}
	if len(pid) != 0 && !dao.IsDepartIDExist(tmp) {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"部门PID有误,数据库中不存在该部门PID")
		return
	}
	if len(pid) == 0 && len(parentDepart) != 0 {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"请输入PID")
		return
	}
	departNode := model.DepartNode{
		PID:          tmp,
		ParentDepart: parentDepart,
		Depart:       depart,
	}
	if dao.IsDepartNodeExist(parentDepart, depart) {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"该部门节点已存在")
		return
	}
	if len(parentDepart) != 0 && !dao.IsParentDepartExist(parentDepart) {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"该部门上级部门不存在")
		return
	}
	if len(depart) == 0 {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"部门节点不能为空")
		return
	} else if len(parentDepart) == 0 {
		if dao.IsRootExist() {
			response.Response(c, http.StatusOK,
				http.StatusUnprocessableEntity,
				nil,
				"已存在根节点,即组织名称")
			return
		} else {
			departNode.NodeLocate = Departroot
			if dao.AddDepart(mysqlmanager.DB, &departNode) != nil {
				logger.Error("部门节点添加失败")
				return
			}
		}
	} else {
		departNode.NodeLocate = DepartUnroot
		if dao.AddDepart(mysqlmanager.DB, &departNode) != nil {
			logger.Error("部门节点添加失败")
			return
		}
	}
	response.Success(c, nil, "部门信息入库成功")
}

func DepartInfo(c *gin.Context) {
	depart := dao.DepartStore()
	if len(depart) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": model.MachineTreeNode{},
		})
		return
	}
	ptrchild, departRoot := service.Returnptrchild(depart)
	service.MakeTree(&departRoot, ptrchild)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": departRoot})
}

func Deletedepartdata(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
	var a model.DeleteDepart
	err = json.Unmarshal(j, &a)
	if err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
	tmp := strconv.Itoa(a.DepartID)

	for _, n := range dao.MachineStore(a.DepartID) {
		dao.ModifyMachineDepart2(n.ID, 1)
	}
	for _, depart := range ReturnID(a.DepartID) {
		machine := dao.MachineStore(depart)
		for _, m := range machine {
			dao.ModifyMachineDepart2(m.ID, 1)
		}
	}
	if !dao.IsDepartIDExist(a.DepartID) {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"不存在该机器")
		return
	}

	DepartInfo := dao.GetPid(tmp)
	service.Deletedepartnode(DepartInfo, a.DepartID)
	var user model.User
	mysqlmanager.DB.Where("depart_second=?", a).Unscoped().Delete(user)
	response.Success(c, nil, "部门删除成功")
}

func MachineInfo(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	depart := &model.Depart{}
	if c.ShouldBind(depart) != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"parameter error")
		return
	}

	var a []int
	ReturnSpecifiedDepart(depart.ID, &a)
	a = append(a, depart.ID)
	machinelist := make([]model.Res, 0)
	for _, value := range a {
		list := &[]model.Res{}
		mysqlmanager.DB.Table("machine_node").Where("depart_id=?", value).Select("machine_node.id as id,machine_node.depart_id as departid," +
			"depart_node.depart as departname,machine_node.ip as ip,machine_node.machine_uuid as uuid, " +
			"machine_node.cpu as cpu,machine_node.state as state, machine_node.systeminfo as systeminfo").Joins("left join depart_node on machine_node.depart_id = depart_node.id").Scan(&list)
		for _, value1 := range *list {
			if value1.Departid == value {
				machinelist = append(machinelist, value1)
			}
		}
	}
	lens := len(machinelist)

	data, err := DataPaging(query, machinelist, lens)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	JsonPagination(c, data, lens, query)
}

//资源池返回接口
func FreeMachineSource(c *gin.Context) {
	departid := 1
	machine := model.MachineNode{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if HandleError(c, err) {
		return
	}

	list, tx, res := machine.ReturnMachine(query, departid)
	total, err := CrudAll(query, tx, &res)
	if HandleError(c, err) {
		return
	}

	// 返回数据开始拼装分页的json
	JsonPagination(c, list, total, query)
}
func MachineAllData(c *gin.Context) {
	AllData := model.MachineAllData()
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": AllData,
	})
}
func Dep(c *gin.Context) {
	departID := c.Query("DepartID")
	tmp, err := strconv.Atoi(departID)
	if err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"部门ID有误")
		return
	}
	depart := dao.DepartStore()
	var root model.MachineTreeNode
	departnode := make([]model.MachineTreeNode, 0)
	ptrchild := make([]*model.MachineTreeNode, 0)

	for _, value := range depart {
		if value.NodeLocate == Departroot {
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
	node := &root
	service.MakeTree(node, ptrchild)
	var d *model.MachineTreeNode
	if node.Id != tmp {
		service.LoopTree(node, tmp, &d)
		node = d
	}
	if node == nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			"部门ID有误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": node,
	})
}

type NewDepart struct {
	DepartID   int    `json:"DepartID"`
	DepartName string `json:"DepartName"`
}

func UpdateDepart(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
	var new NewDepart
	err = json.Unmarshal(j, &new)
	if err != nil {
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
	dao.UpdateDepart(new.DepartID, new.DepartName)
	dao.UpdateParentDepart(new.DepartID, new.DepartName)
	response.Success(c, nil, "部门更新成功")
}

type modify struct {
	MachineID string `json:"machineid"`
	DepartID  int    `json:"departid"`
}

func ModifyMachineDepart(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	logger.Info(string(j))
	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
	var M modify
	err = json.Unmarshal(j, &M)

	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			http.StatusUnprocessableEntity,
			nil,
			err.Error())
		return
	}
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
func AddIP(c *gin.Context) {
	IP := c.Query("ip")
	uuid := c.Query("uuid")
	var MachineInfo model.MachineNode
	Machine := model.MachineNode{
		IP: IP,
	}
	mysqlmanager.DB.Model(&MachineInfo).Where("machine_uuid=?", uuid).Update(&Machine)
	response.Success(c, nil, "ip更新成功")
}

func ReturnID(id int) []int {
	var depart []model.DepartNode
	mysqlmanager.DB.Where("p_id=?", id).Find(&depart)

	res := make([]int, 0)
	for _, value := range depart {
		res = append(res, value.ID)
	}
	return res
}

//返回所有子部门函数
func ReturnSpecifiedDepart(id int, res *[]int) {
	if len(ReturnID(id)) == 0 {
		return
	}
	for _, value := range ReturnID(id) {
		*res = append(*res, value)
		ReturnSpecifiedDepart(value, res)
	}
}
