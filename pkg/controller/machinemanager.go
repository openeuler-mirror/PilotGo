package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/model"
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
func Postmachinedata(c *gin.Context) {}

// func Getmachinedata(c *gin.Context) {
// 	company := c.PostForm("company")
// 	primdepart := c.PostForm(("primaryDepart"))
// 	secondepart := c.PostForm("secondaryDepart")
// 	tertpart := c.PostForm("tertiaryDepart")
// 	id := c.PostForm("UserID")
// 	uid := c.PostForm("machineUID")
// 	if (len(company)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"公司不能为空")
// 		return
// 	}
// 	if (len(primdepart)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"一级部门不能为空")
// 		return
// 	}
// 	if (len(secondepart)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"二级部门不能为空")
// 		return
// 	}
// 	if (len(tertpart)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"三级部门不能为空")
// 		return
// 	}
// 	if (len(id)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"id不能为空")
// 		return
// 	}
// 	if (len(uid)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"uid不能为空")
// 		return
// 	}
// 	machine := fmt.Sprintf("%s/%s/%s/%s/%s/%s", company, primdepart, secondepart, tertpart, id, uid)
// 	tert := fmt.Sprintf("%s/%s/%s/%s", company, primdepart, secondepart, tertpart)
// 	second := fmt.Sprintf("%s/%s/%s", company, primdepart, secondepart)
// 	prim := fmt.Sprintf("%s/%s", company, primdepart)

// 	// machineman := model.MachineNode{
// 	// 	Companyname: company,
// 	// 	Primdepart:  primdepart,
// 	// 	Secondepart: secondepart,
// 	// 	Tertpart:    tertpart,
// 	// 	ID:          id,
// 	// 	UID:         uid,
// 	// }
// 	machineman := model.MachineManage{
// 		Data: company,
// 	}
// 	if dao.IsMachineinfoExist(company) {
// 		logger.Info("该机器的所属公司已入库")
// 	} else {
// 		mysqlmanager.DB.Create(&machineman)
// 	}

// 	machineman = model.MachineManage{
// 		Data: prim,
// 	}
// 	if dao.IsMachineinfoExist(prim) {
// 		logger.Info("该机器的一级部门已入库")
// 	} else {
// 		mysqlmanager.DB.Create(&machineman)
// 	}

// 	machineman = model.MachineManage{
// 		Data: second,
// 	}
// 	if dao.IsMachineinfoExist(second) {
// 		logger.Info("该机器的二级部门已入库")
// 	} else {
// 		mysqlmanager.DB.Create(&machineman)
// 	}

// 	machineman = model.MachineManage{
// 		Data: tert,
// 	}
// 	if dao.IsMachineinfoExist(tert) {
// 		logger.Info("该机器的三级部门已入库")
// 	} else {
// 		mysqlmanager.DB.Create(&machineman)
// 	}

// 	machineman = model.MachineManage{
// 		Data: machine,
// 	}
// 	if dao.IsMachineinfoExist(machine) {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"该机器已经存在")
// 		return
// 	} else {
// 		mysqlmanager.DB.Create(&machineman)
// 	}
// 	response.Success(c, nil, "机器入库成功")
// }
// func DeleteMachineInfo(c *gin.Context) {
// 	company := c.PostForm("company")
// 	primdepart := c.PostForm(("primaryDepart"))
// 	secondepart := c.PostForm("secondaryDepart")
// 	tertpart := c.PostForm("tertiaryDepart")
// 	id := c.PostForm("UserID")
// 	uid := c.PostForm("machineUID")
// 	if (len(company)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"公司不能为空")
// 		return
// 	} else if (len(primdepart)) == 0 {
// 		if (len(secondepart)) != 0 || (len(tertpart)) != 0 || (len(id)) != 0 || (len(uid)) != 0 {
// 			response.Response(c, http.StatusUnprocessableEntity,
// 				422,
// 				nil,
// 				"一级部门不能为空")
// 			return
// 		} else {
// 			if dao.IsMachineinfoExist(company) {
// 				var machineman model.MachineManage
// 				mysqlmanager.DB.Where("data like ?", company+"%%").Delete(machineman)
// 			} else {
// 				response.Response(c, http.StatusUnprocessableEntity,
// 					422,
// 					nil,
// 					"数据库不存在该公司!")
// 				return
// 			}
// 		}
// 	} else if (len(secondepart)) == 0 {
// 		if (len(tertpart)) != 0 || (len(id)) != 0 || (len(uid)) != 0 {
// 			response.Response(c, http.StatusUnprocessableEntity,
// 				422,
// 				nil,
// 				"二级部门不能为空")
// 			return
// 		} else {
// 			if dao.IsMachineinfoExist(company + "/" + primdepart) {
// 				var machineman model.MachineManage
// 				mysqlmanager.DB.Where("data like ?", company+"/"+primdepart+"%%").Delete(machineman)
// 			} else {
// 				response.Response(c, http.StatusUnprocessableEntity,
// 					422,
// 					nil,
// 					"数据库不存在该公司一级部门!")
// 				return
// 			}
// 		}
// 	} else if (len(tertpart)) == 0 {
// 		if (len(id)) != 0 || (len(uid)) != 0 {
// 			response.Response(c, http.StatusUnprocessableEntity,
// 				422,
// 				nil,
// 				"三级部门不能为空")
// 			return
// 		} else {
// 			if dao.IsMachineinfoExist(company + "/" + primdepart + "/" + secondepart) {
// 				var machineman model.MachineManage
// 				mysqlmanager.DB.Where("data like ?", company+"/"+primdepart+"/"+secondepart+"%%").Delete(machineman)
// 			} else {
// 				response.Response(c, http.StatusUnprocessableEntity,
// 					422,
// 					nil,
// 					"数据库不存在该公司二级部门!")
// 				return
// 			}
// 		}
// 	} else if (len(id)) == 0 {
// 		if (len(uid)) != 0 {
// 			response.Response(c, http.StatusUnprocessableEntity,
// 				422,
// 				nil,
// 				"id不能为空")
// 			return
// 		} else {
// 			if dao.IsMachineinfoExist(company + "/" + primdepart + "/" + secondepart + "/" + tertpart) {
// 				var machineman model.MachineManage
// 				mysqlmanager.DB.Where("data like ?", company+"/"+primdepart+"/"+secondepart+"/"+tertpart+"%%").Delete(machineman)
// 			} else {
// 				response.Response(c, http.StatusUnprocessableEntity,
// 					422,
// 					nil,
// 					"数据库不存在该公司三级部门!")
// 				return
// 			}
// 		}
// 	} else if (len(uid)) == 0 {
// 		response.Response(c, http.StatusUnprocessableEntity,
// 			422,
// 			nil,
// 			"uid不能为空")
// 		return
// 	} else {
// 		machine := fmt.Sprintf("%s/%s/%s/%s/%s/%s", company, primdepart, secondepart, tertpart, id, uid)
// 		if dao.IsMachineinfoExist(machine) {
// 			var machineman model.MachineManage
// 			mysqlmanager.DB.Where("data=?", machine).Delete(machineman)
// 		} else {
// 			response.Response(c, http.StatusUnprocessableEntity,
// 				422,
// 				nil,
// 				"数据库不存在该机器!")
// 			return
// 		}
// 	}
// 	response.Success(c, nil, "机器删除成功")

// }
