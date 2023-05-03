package pluginapi

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func PluginList(c *gin.Context) {
	data := service.GetPlugins()

	response.Success(c, data, "")

}
