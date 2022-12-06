package controller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 存储脚本文件
func AddScriptHandler(c *gin.Context) {
	var script model.Script
	err := service.AddScript(&script)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "脚本文件添加失败")
		return
	}
	response.Success(c, nil, "脚本文件添加成功")
}
