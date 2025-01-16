/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package pluginapi

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/tag"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 插件返回tag数据
func GetTagHandler(c *gin.Context) {
	uuidTags := &struct {
		UUIDS []string `json:"uuids"`
	}{}
	if err := c.ShouldBindJSON(&uuidTags); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	//获取到机器标签
	data, err := tag.RequestTag(uuidTags.UUIDS)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	response.Success(c, data, "get tag成功")
}
