package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func PluginList(c *gin.Context) {
	pluginlists := service.PluginLists()
	response.JSON(c, http.StatusOK, http.StatusOK, pluginlists, "插件列表获取成功")
}

func LoadPlugin(c *gin.Context) {
	var pluginInfo model.LoadPlugin
	c.Bind(&pluginInfo)
	err := service.LoadPlugin(pluginInfo)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, nil, "插件加载成功")
}

func UnLoadPlugin(c *gin.Context) {
	var pluginInfo model.UnLoadPlugin
	c.Bind(&pluginInfo)
	err := service.UnLoadPlugin(pluginInfo)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.Success(c, nil, "插件卸载成功")
}
