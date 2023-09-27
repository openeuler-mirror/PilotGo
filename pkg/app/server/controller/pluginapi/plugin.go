package pluginapi

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/plugin"
	"openeuler.org/PilotGo/PilotGo/sdk/response"
)

func PluginList(c *gin.Context) {
	data := plugin.GetPlugins()

	response.Success(c, data, "")

}
