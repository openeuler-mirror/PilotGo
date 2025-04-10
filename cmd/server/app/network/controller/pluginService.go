/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Apr 10 16:18:53 2025 +0800
 */
package controller

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 查询插件服务清单
func GetPluginServices(c *gin.Context) {
	pluginServices := global.GW.GetAllServices()

	response.Success(c, pluginServices, "插件查询成功")
}

// 分页查询插件服务清单
func GetPluginServicesPaged(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	pluginServices := global.GW.GetAllServices()
	data, err := response.DataPaging(query, pluginServices, len(pluginServices))
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.DataPagination(c, data, len(pluginServices), query)
}

// 停用/启动插件服务
func TogglePluginService(c *gin.Context) {
	param := struct {
		ServiceName string `json:"serviceName"`
		Enable      bool   `json:"enable"`
	}{}
	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	if err := global.GW.SetServiceStatus(param.ServiceName, param.Enable); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, nil, fmt.Sprintf("插件信息更新成功,当前%v服务状态为: %v", param.ServiceName, global.GW.GetServiceStatus(param.ServiceName)))
}
