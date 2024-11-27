/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package pluginapi

import (
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func BatchListHandler(c *gin.Context) {
	batch, err := batch.SelectBatch()
	if err != nil {
		response.Fail(c, nil, "获取批次信息错误"+err.Error())
		return
	}

	if len(batch) == 0 {
		response.Fail(c, nil, "未获取到批次信息")
		return
	}

	resp := []*common.BatchList{}
	for _, item := range batch {
		res := &common.BatchList{
			ID:          int(item.ID),
			Name:        item.Name,
			Description: item.Description,
			Manager:     item.Manager,
		}
		resp = append(resp, res)
	}
	response.Success(c, resp, "批次信息获取成功")
}

func MachineListOfBatch(c *gin.Context) {
	BatchId := c.Query("batchId")
	batchId, err := strconv.Atoi(BatchId)
	if err != nil {
		response.Fail(c, nil, "批次ID输入格式有误")
		return
	}
	machine_uuids := batch.GetMachineUUIDS(batchId)
	response.Success(c, machine_uuids, "获取到批次的uuid列表")
}
