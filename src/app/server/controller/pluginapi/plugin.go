package pluginapi

import (
	"gitee.com/PilotGo/PilotGo/app/server/service/plugin"
	"gitee.com/PilotGo/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginList(c *gin.Context) {
	data := plugin.GetPlugins()

	response.Success(c, data, "")

}
