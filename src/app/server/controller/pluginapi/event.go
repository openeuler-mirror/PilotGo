//

package pluginapi

import (
	"gitee.com/PilotGo/PilotGo/app/server/service/eventbus"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/sdk/plugin/client"
	"gitee.com/PilotGo/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func RegisterListenerHandler(c *gin.Context) {
	p := client.PluginInfo{}
	if err := c.ShouldBindQuery(p); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	eventbus.AddListener(&eventbus.Listener{
		Name: p.Name,
		URL:  p.Url,
	})

	logger.Info("")
}

func UnregisterListenerHandler(c *gin.Context) {
	p := client.PluginInfo{}
	if err := c.ShouldBindQuery(p); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	eventbus.RemoveListener(&eventbus.Listener{
		Name: p.Name,
		URL:  p.Url,
	})
}
