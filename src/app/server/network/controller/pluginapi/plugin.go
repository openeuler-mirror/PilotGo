package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginList(c *gin.Context) {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		response.Fail(c, nil, "查询插件错误："+err.Error())
		return
	}
	response.Success(c, plugins, "插件查询成功")
}
