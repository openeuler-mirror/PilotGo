package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func CreateBatchHandler(c *gin.Context) {
	var batchinfo model.CreateBatch
	if err := c.Bind(&batchinfo); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	if err := service.CreateBatch(&batchinfo); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, nil, "批次入库成功")
}

func BatchInfoHandler(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	list, total, err := service.GetBatches(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 返回数据开始拼装分页的json
	service.JsonPagination(c, list, total, query)
}

func DeleteBatchHandler(c *gin.Context) {
	var batchdel model.BatchDel
	if err := c.Bind(&batchdel); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if len(batchdel.BatchID) == 0 {
		response.Fail(c, nil, "请输入删除批次ID")
		return
	}

	if err := service.DeleteBatch(batchdel.BatchID); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	response.Success(c, nil, "批次删除成功")
}

func UpdateBatchHandler(c *gin.Context) {
	var batchinfo model.BatchUpdate
	if err := c.Bind(&batchinfo); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := service.UpdateBatch(batchinfo.BatchId, batchinfo.BatchName, batchinfo.Description)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, "update batch failed: "+err.Error())
		return
	}

	response.Success(c, nil, "批次修改成功")
}

func BatchMachineInfoHandler(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	Batchid := c.Query("ID")
	batchid, err := strconv.Atoi(Batchid)
	if err != nil {
		response.Fail(c, nil, "批次ID输入格式有误")
		return
	}

	machinesInfo, err := service.GetBatchMachines(batchid)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, "get batch machines failed: "+err.Error())
		return
	}

	// 分页
	data, err := service.DataPaging(query, machinesInfo, len(machinesInfo))
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	service.JsonPagination(c, data, int64(len(machinesInfo)), query)
}

func SelectBatchHandler(c *gin.Context) {
	batch, err := service.SelectBatch()
	if err != nil {
		response.Fail(c, nil, "获取批次信息错误"+err.Error())
		return
	}

	if len(batch) == 0 {
		response.Fail(c, nil, "未获取到批次信息")
		return
	}
	response.Success(c, gin.H{"data": batch}, "批次信息获取成功")
}
