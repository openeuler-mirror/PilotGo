package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginList(c *gin.Context) {
	data := plugin.GetPlugins()

	response.Success(c, data, "")

}
