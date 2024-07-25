package controller

import (
	scriptservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/script"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 存储脚本文件
func AddScriptHandler(c *gin.Context) {
	var script scriptservice.Script
	err := scriptservice.AddScript(&script)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "脚本文件添加失败")
		return
	}
	response.Success(c, nil, "脚本文件添加成功")
}
