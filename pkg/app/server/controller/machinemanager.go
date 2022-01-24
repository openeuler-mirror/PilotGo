package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func AddDepart(c *gin.Context) {
	pid := c.PostForm("PID")
	parentDepart := c.PostForm("ParentDepart")
	depart := c.PostForm("Depart")
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
func AddMachine(c *gin.Context) {
	departID := c.PostForm("DepartID")
	machineuuid := c.PostForm("MachineUUID")
	if len(departID) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID不能为空")
		return
	}
	if len(machineuuid) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			442,
			nil,
			"机器uuid不能为空")
		return
	}
	tmp, err := strconv.Atoi(departID)
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
			"部门ID有误,数据库中不存在该部门ID")
		return
	}
	machinenode := model.MachineNode{
		DepartId:    tmp,
		MachineUUID: machineuuid,
	}
	mysqlmanager.DB.Create(&machinenode)
	response.Success(c, nil, "机器入库成功")
}
func DepartInfo(c *gin.Context) {
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
	for key, _ := range departnode {
		var a *model.MachineTreeNode
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
func Deletemachinedata(c *gin.Context) {
	uuid := c.PostForm("uuid")
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
	departid := c.PostForm("DepartID")
	needdelete := make([]int, 0)
	tmp, err := strconv.Atoi(departid)
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
	DepartInfo := dao.GetPid(departid)
	needdelete = append(needdelete, tmp)
	for _, value := range DepartInfo {
		needdelete = append(needdelete, value.ID)
	}

	for {
		if len(needdelete) == 0 {
			break
		}
		dao.Deletedepartdata(needdelete)
		dao.Insertdepartlist(needdelete)
	}
	response.Success(c, nil, "部门删除成功")
}

func MachineInfo(c *gin.Context) {
	departid := c.Query("DepartId")
	tmp, err := strconv.Atoi(departid)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"部门ID输入格式有误")
		return
	}
	machineinformation := dao.MachineStore(tmp)
	var uuid model.MachineInfo
	for _, value := range machineinformation {
		uuid.Uuid = append(uuid.Uuid, value.MachineUUID)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": uuid,
	})
}
