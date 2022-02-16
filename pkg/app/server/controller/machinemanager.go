package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func LoopTree(node *model.MachineTreeNode, ID int) *model.MachineTreeNode {
	for _, value := range node.Children {
		if value.Id == ID {
			return value
		}
		LoopTree(value, ID)
	}
	return &model.MachineTreeNode{}
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
	for key, _ := range departnode {
		var a *model.MachineTreeNode
		a = &departnode[key]
		ptrchild = append(ptrchild, a)
	}
	node := &root
	makeTree(node, ptrchild)
	if node.Id != tmp {
		node = LoopTree(node, tmp)
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
