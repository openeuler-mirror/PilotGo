package controller

import (
	"net/http"
	"strconv"
	"strings"

	"gitee.com/PilotGo/PilotGo/app/server/config"
	"gitee.com/PilotGo/PilotGo/app/server/service/auditlog"
	"gitee.com/PilotGo/PilotGo/app/server/service/batch"
	"gitee.com/PilotGo/PilotGo/app/server/service/common"
	"gitee.com/PilotGo/PilotGo/app/server/service/jwt"
	"gitee.com/PilotGo/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func CreateBatchHandler(c *gin.Context) {
	params := &struct {
		Name        string   `json:"Name"`
		Description string   `json:"Descrip"`
		Manager     string   `json:"Manager"`
		DepartName  []string `json:"DepartName"`
		DepartID    []int    `json:"DepartID"`
		Machines    []int    `json:"Machines"`
	}{}
	if err := c.Bind(params); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypeBatch, "创建批次", "", user)
	auditlog.Add(log)

	batchinfo := &batch.CreateBatchParam{
		Name:        params.Name,
		Description: params.Description,
		Manager:     params.Manager,
		DepartName:  params.DepartName,
		DepartID:    params.DepartID,
		Machines:    params.Machines,
	}

	if err := batch.CreateBatch(batchinfo); err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, params.Name, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, nil, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, params.Name, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
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

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypeBatch, "删除批次", "", user)
	auditlog.Add(log)

	batchesname := []string{}
	for _, batchid := range batchdel.BatchID {
		batchesname = append(batchesname, strconv.Itoa(batchid))
	}

	if err := batch.DeleteBatch(batchdel.BatchID); err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, strings.Join(batchesname, "/"), http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, strings.Join(batchesname, "/"), http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
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

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log := auditlog.NewByUser(auditlog.LogTypeBatch, "编辑批次", "", user)
	auditlog.Add(log)

	err = batch.UpdateBatch(batchinfo.BatchId, batchinfo.BatchName, batchinfo.Description)
	if err != nil {
		log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, err.Error(), log.Module, batchinfo.BatchName, http.StatusBadRequest)
		auditlog.Add(log_s)
		auditlog.UpdateStatus(log, auditlog.ActionFalse)
		response.Fail(c, gin.H{"status": false}, "update batch failed: "+err.Error())
		return
	}

	log_s := auditlog.New_sub(log.LogUUID, strings.Split(config.Config().HttpServer.Addr, ":")[0], log.Action, "", log.Module, batchinfo.BatchName, http.StatusOK)
	auditlog.Add(log_s)
	auditlog.UpdateStatus(log, auditlog.ActionOK)
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
