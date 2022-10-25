package controller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 查询插件清单
func GetPluginsHanlder(c *gin.Context) {
	logger.Info("query plugin")
	plugins := service.GetPlugins()

	response.NewSuccess(c, plugins, "插件查询")
}

// 添加插件
func AddPluginHanlder(c *gin.Context) {
	param := struct {
		Url string `json:"url"`
	}{}

	if err := c.Bind(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	logger.Info("add plugin from: %s", param.Url)
	if err := service.AddPlugin(param.Url); err != nil {
		response.Fail(c, nil, "add plugin failed:"+err.Error())
		return
	}

	response.Success(c, nil, "插件添加成功")
}

// 停用/启动插件
func TogglePluginHanlder(c *gin.Context) {
	logger.Info("disable plugin")
	// TODO:

	response.Success(c, nil, "插件信息更新成功")
}

// 卸载插件
func UnloadPluginHanlder(c *gin.Context) {
	logger.Info("disable plugin")
	// TODO:

	response.Success(c, nil, "插件信息更新成功")
}
