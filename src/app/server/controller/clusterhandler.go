/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-03-24 00:46:05
 * LastEditTime: 2022-04-29 11:29:44
 * Description: 集群概览
 ******************************************************************************/
package controller

import (
	"gitee.com/PilotGo/PilotGo/app/server/service/cluster"
	"gitee.com/PilotGo/PilotGo/sdk/response"
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
