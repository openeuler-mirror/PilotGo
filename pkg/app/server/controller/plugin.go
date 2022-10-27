package controller

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

// 查询插件清单
func GetPluginsHanlder(c *gin.Context) {
	plugins := service.GetPlugins()

	logger.Info("find %d plugins", len(plugins))
	response.NewSuccess(c, plugins, "插件查询成功")
}

// 添加插件
func AddPluginHanlder(c *gin.Context) {
	param := struct {
		Url string `json:"url"`
	}{}

	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	if err := service.AddPlugin(param.Url); err != nil {
		response.Fail(c, nil, "add plugin failed:"+err.Error())
		return
	}

	response.Success(c, nil, "插件添加成功")
}

// 停用/启动插件
func TogglePluginHanlder(c *gin.Context) {
	param := struct {
		UUID   string `json:"uuid"`
		Enable int    `json:"enable"`
	}{}

	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	logger.Info("toggle plugin:%s to enable %d", param.UUID, param.Enable)
	// TODO:
	service.TogglePlugin(param.UUID, param.Enable)
	response.Success(c, nil, "插件信息更新成功")
}

// 卸载插件
func UnloadPluginHanlder(c *gin.Context) {
	param := struct {
		Id int `json:"id"`
	}{}
	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}
	logger.Info("unload plugin:%d", param.Id)
	service.DeletePlugin(param.Id)

	response.Success(c, nil, "插件信息更新成功")
}
