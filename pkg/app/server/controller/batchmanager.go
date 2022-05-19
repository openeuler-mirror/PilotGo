package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

// 创建批次页面机器列表树
func MachineList(c *gin.Context) {
	DepartId := c.Query("DepartId")
	DepId, err := strconv.Atoi(DepartId)
	if err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	var departId []int
	ReturnSpecifiedDepart(DepId, &departId)
	departId = append(departId, DepId)

	machinelist := dao.MachineList(departId)
	response.JSON(c, http.StatusOK, http.StatusOK, machinelist, "部门下所属机器获取成功")
}

func CreateBatch(c *gin.Context) {
	var batchinfo model.CreateBatch
	c.Bind(&batchinfo)

	if len(batchinfo.Name) == 0 {
		response.Fail(c, nil, "请输入批次名称")
		return
	}
	if dao.IsExistName(batchinfo.Name) {
		response.Fail(c, nil, "已存在该名称批次")
		return
	}
	if len(batchinfo.Manager) == 0 {
		response.Fail(c, nil, "创建人未输入")
		return
	}

	if len(batchinfo.Machines) == 0 && len(batchinfo.DepartIDs) == 0 {
		response.Fail(c, nil, "请先选择机器IP或部门")
		return
	}

	// 机器id列表
	var machinelist string
	Departids := make([]int, 0)
	if len(batchinfo.Machines) == 0 {
		// 点选部门创建批次
		var machineids []int
		for _, ids := range batchinfo.DepartIDs {
			Departids = append(Departids, ids)
			ReturnSpecifiedDepart(ids, &Departids)
		}

		machines := dao.MachineList(Departids)
		for _, mamachine := range machines {
			machineids = append(machineids, mamachine.ID)
		}
		for _, id := range machineids {
			machinelist = machinelist + "," + strconv.Itoa(id)
			machinelist = strings.Trim(machinelist, ",")
		}
	} else {
		// 点选ip创建批次
		for _, id := range batchinfo.Machines {
			machinelist = machinelist + "," + strconv.Itoa(id)
			machinelist = strings.Trim(machinelist, ",")
		}
	}

	// 机器所属部门ids
	var departIdlist string
	if len(batchinfo.DepartID) == 0 {
		for _, id := range Departids {
			departIdlist = departIdlist + "," + strconv.Itoa(id)
			departIdlist = strings.Trim(departIdlist, ",")
		}
	} else {
		for _, id := range batchinfo.DepartID {
			departIdlist = departIdlist + "," + strconv.Itoa(id)
			departIdlist = strings.Trim(departIdlist, ",")
		}
	}

	// 机器所属部门
	var departNamelist string
	if len(batchinfo.DepartID) == 0 {
		list := dao.DepartIdsToGetDepartNames(Departids)
		departNamelist = strings.Join(list, ",")
	} else {
		List := dao.DepartIdsToGetDepartNames(batchinfo.DepartID)
		departNamelist = strings.Join(List, ",")
	}

	Batch := model.Batch{
		Name:        batchinfo.Name,
		Description: batchinfo.Description,
		Manager:     batchinfo.Manager,
		Depart:      departIdlist,
		DepartName:  departNamelist,
		Machinelist: machinelist,
	}
	dao.CreateBatch(Batch)
	response.Success(c, nil, "批次入库成功")
}

func BatchInfo(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	batch := model.Batch{}
	list, tx := batch.ReturnBatch(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, err := CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	JsonPagination(c, list, total, query)
}

func DeleteBatch(c *gin.Context) {
	var batchdel model.BatchDel
	c.Bind(&batchdel)

	if len(batchdel.BatchID) == 0 {
		response.Response(c, http.StatusOK, http.StatusUnprocessableEntity, nil, "请输入删除批次ID")
		return
	}
	for _, value := range batchdel.BatchID {
		dao.DeleteBatch(value)
	}
	response.Success(c, nil, "批次删除成功")
}

func UpdateBatch(c *gin.Context) {
	var batchinfo model.BatchUpdate
	c.Bind(&batchinfo)

	if !dao.IsExistID(batchinfo.BatchId) {
		response.Response(c, http.StatusOK, http.StatusUnprocessableEntity, nil, "不存在该批次")
		return
	}
	dao.UpdateBatch(batchinfo.BatchId, batchinfo.BatchName, batchinfo.Description)
	response.Success(c, nil, "批次修改成功")
}

func BatchMachineInfo(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	Batchid := c.Query("ID")
	batchid, err := strconv.Atoi(Batchid)
	if err != nil {
		response.Fail(c, nil, "批次ID输入格式有误")
		return
	}

	machinelist := dao.GetMachineID(batchid)
	machineIdlist := String2Int(machinelist) // 获取批次里所有机器的id

	// 获取机器的所有信息
	MachinesInfo := make([]model.MachineNode, 0)
	for _, macId := range machineIdlist {
		MacInfo := dao.MachineData(macId)
		MachinesInfo = append(MachinesInfo, MacInfo)
	}
	// 分页
	data, err := DataPaging(query, MachinesInfo, len(MachinesInfo))
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	JsonPagination(c, data, len(MachinesInfo), query)
}
