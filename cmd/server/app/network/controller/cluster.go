/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/cluster"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ClusterInfoHandler(c *gin.Context) {
	data, err := cluster.ClusterInfo()
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, gin.H{"data": data}, "集群概览获取成功")
}

func DepartClusterInfoHandler(c *gin.Context) {
	departs := cluster.DepartClusterInfo()
	if len(departs) == 0 {
		response.Success(c, gin.H{"data": []interface{}{}}, "获取各部门集群状态成功")
	} else {
		response.Success(c, gin.H{"data": departs}, "获取各部门集群状态成功")
	}
}
