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
 * LastEditTime: 2022-04-06 16:23:06
 * Description: provide machine manager functions.
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func AddDepart(c *gin.Context) {
	pid := c.Query("PID")
	parentDepart := c.Query("ParentDepart")
	depart := c.Query("Depart")
	tmp, err := strconv.Atoi(pid)
	if len(pid) != 0 && err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"pid识别失败")
		return
	}
	if len(pid) != 0 && !dao.IsDepartIDExist(tmp) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门PID有误,数据库中不存在该部门PID")
		return
	}
	if len(pid) == 0 && len(parentDepart) != 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
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
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"该部门节点已存在")
		return
	}
	if len(parentDepart) != 0 && !dao.IsParentDepartExist(parentDepart) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"该部门上级部门不存在")
		return
	}
	if len(depart) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门节点不能为空")
		return
	} else if len(parentDepart) == 0 {
		if dao.IsRootExist() {
			response.Response(c, http.StatusUnprocessableEntity,
				422,
				nil,
				"已存在根节点,即组织名称")
			return
		} else {
			departNode.NodeLocate = 0
			mysqlmanager.DB.Create(&departNode)
		}
	} else {
		departNode.NodeLocate = 1
		mysqlmanager.DB.Create(&departNode)
	}
	response.Success(c, nil, "部门信息入库成功")
}

func DepartInfo(c *gin.Context) {
	depart := dao.DepartStore()
	if len(depart) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": model.MachineTreeNode{},
		})
		return
	}
	var root model.MachineTreeNode
	departnode := make([]model.MachineTreeNode, 0)
	ptrchild := make([]*model.MachineTreeNode, 0)

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
	node := &root
	makeTree(node, ptrchild)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": node,
	})
}
func makeTree(node *model.MachineTreeNode, ptrchild []*model.MachineTreeNode) {
	childs := findchild(node, ptrchild)
	for _, value := range childs {
		node.Children = append(node.Children, value)
		if IsChildExist(value, ptrchild) {
			makeTree(value, ptrchild)
		}
	}
}
func findchild(node *model.MachineTreeNode, ptrchild []*model.MachineTreeNode) (ret []*model.MachineTreeNode) {
	for _, value := range ptrchild {
		if node.Id == value.Pid {
			ret = append(ret, value)
		}
	}
	return
}
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
func Deletemachinedata(c *gin.Context) {
	uuid := c.Query("uuid")
	logger.Info("%s", uuid)
	var Machine model.MachineNode
	logger.Info("%+v", Machine)
	if !dao.IsUUIDExist(uuid) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"不存在该机器")
		return
	} else {
		dao.Deleteuuid(uuid)
		response.Success(c, nil, "机器删除成功")
	}
}

func Deletedepartdata(c *gin.Context) {
	a := c.Query("DepartID")
	logger.Info("%s", a)
	tmp, err := strconv.Atoi(a)
	logger.Info("%d", tmp)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID有误")
		return
	}
	if !dao.IsDepartIDExist(tmp) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"不存在该机器")
		return
	}
	needdelete := make([]int, 0)
	DepartInfo := dao.GetPid(a)
	needdelete = append(needdelete, tmp)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}

	for {
		if len(needdelete) == 0 {
			break
		}
		logger.Info("%d", needdelete[0])
		dao.Deletedepartdata(needdelete)
		str := fmt.Sprintf("%d", needdelete[0])
		needdelete = needdelete[1:]
		dao.Insertdepartlist(needdelete, str)

	}
	var user model.User
	mysqlmanager.DB.Where("depart_second=?", a).Unscoped().Delete(user)
	response.Success(c, nil, "部门删除成功")
}

func MachineInfo(c *gin.Context) {
	departid := c.Query("DepartId")
	machine := model.MachineNode{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	tmp, err := strconv.Atoi(departid)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID输入格式有误")
		return
	}
	list, total, err := machine.ReturnMachine(query, tmp)
	if model.HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	model.JsonPagination(c, list, total, query)

}

func MachineAllData(c *gin.Context) {
	AllData := model.MachineAllData()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": AllData,
	})
}
func Dep(c *gin.Context) {
	departID := c.Query("DepartID")
	tmp, err := strconv.Atoi(departID)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID有误")
		return
	}
	depart := dao.DepartStore()
	var root model.MachineTreeNode
	departnode := make([]model.MachineTreeNode, 0)
	ptrchild := make([]*model.MachineTreeNode, 0)

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
	node := &root
	makeTree(node, ptrchild)
	var d *model.MachineTreeNode
	if node.Id != tmp {
		LoopTree(node, tmp, &d)
		node = d
	}
	if node == nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"部门ID有误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": node,
	})
}

func UpdateDepart(c *gin.Context) {
	DepartID := c.Query("DepartID")
	DepartName := c.Query("DepartName")
	tmp, err := strconv.Atoi(DepartID)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID有误")
		return
	}
	dao.UpdateDepart(tmp, DepartName)
	dao.UpdateParentDepart(tmp, DepartName)
	response.Success(c, nil, "部门更新成功")
}

type modify struct {
	MachineID int `json:"machineid"`
	DepartID  int `json:"departid"`
}

func ModifyMachineDepart(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	logger.Info(string(j))
	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	var M modify
	err = json.Unmarshal(j, &M)
	logger.Info("%+v", M)

	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	dao.ModifyMachineDepart(M.MachineID, M.DepartID)
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

func AddAgents() {
	var agent_list model.MachineNode
	agents := agentmanager.GetAgentList()
	for _, agent := range agents {
		uuid := agent["agent_uuid"]
		agent_uuid := agentmanager.GetAgent(uuid)
		if agent_uuid == nil {
			logger.Fatal("获取uuid失败!")
			return
		}
		agent_OS, err := agent_uuid.GetAgentOSInfo()
		if err != nil {
			logger.Fatal("初始化系统信息失败!")
			return
		}
		agentOS := strings.Split(agent_OS.(string), ";")
		agent_list.DepartId = 1
		agent_list.MachineUUID = uuid
		if dao.IsUUIDExist(uuid) {
			logger.Warn("机器%s已经存在!", agentOS[0])
			continue
		}
		agent_list.IP = agentOS[0]
		if dao.IsIPExist(agentOS[0]) {
			mysqlmanager.DB.Where("ip=?", agentOS[0]).Unscoped().Delete(agent_list)
			// mysqlmanager.DB.Model(&agent_list).Where("ip=?", agentOS[0]).Update("state", model.OffLine)
		}
		agent_list.Systeminfo = agentOS[1] + " " + agentOS[2]
		agent_list.CPU = agentOS[3]
		agent_list.State = model.Free
		mysqlmanager.DB.Save(&agent_list)
	}
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
func ReturnSpecifiedDepart(id int, res *[]int) {
	if len(ReturnID(id)) == 0 {
		return
	}
	for _, value := range ReturnID(id) {
		*res = append(*res, value)
		ReturnSpecifiedDepart(value, res)
	}
}
