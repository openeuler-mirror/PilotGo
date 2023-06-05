package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auditlog"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/batch"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	userservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/user"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func CreateBatchHandler(c *gin.Context) {
	var batchinfo batch.CreateBatchParam
	if err := c.Bind(&batchinfo); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeBatch, "创建批次", "", user)
	auditlog.Add(log)

	if err := batch.CreateBatch(&batchinfo); err != nil {
		logger.Debug(err.Error())
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "批次入库成功")
}

func BatchInfoHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	list, total, err := batch.GetBatches(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 返回数据开始拼装分页的json
	common.JsonPagination(c, list, total, query)
}

func DeleteBatchHandler(c *gin.Context) {
	batchdel := struct {
		BatchID []int `json:"BatchID"`
	}{}
	if err := c.Bind(&batchdel); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if len(batchdel.BatchID) == 0 {
		response.Fail(c, nil, "请输入删除批次ID")
		return
	}

	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeBatch, "删除批次", "", user)
	auditlog.Add(log)

	if err := batch.DeleteBatch(batchdel.BatchID); err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "批次删除成功")
}

func UpdateBatchHandler(c *gin.Context) {
	batchinfo := struct {
		BatchId     int    `json:"BatchID"`
		BatchName   string `json:"BatchName"`
		Description string `json:"Description"`
	}{}
	if err := c.Bind(&batchinfo); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	//TODO:
	var user userservice.User
	log := auditlog.New(auditlog.LogTypeBatch, "修改批次", "", user)
	auditlog.Add(log)

	err := batch.UpdateBatch(batchinfo.BatchId, batchinfo.BatchName, batchinfo.Description)
	if err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, gin.H{"status": false}, "update batch failed: "+err.Error())
		return
	}

	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "批次修改成功")
}

func BatchMachineInfoHandler(c *gin.Context) {
	query := &common.PaginationQ{}
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

	machinesInfo, err := batch.GetBatchMachines(batchid)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, "get batch machines failed: "+err.Error())
		return
	}

	// 分页
	data, err := common.DataPaging(query, machinesInfo, len(machinesInfo))
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, int64(len(machinesInfo)), query)
}

func SelectBatchHandler(c *gin.Context) {
	batch, err := batch.SelectBatch()
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
