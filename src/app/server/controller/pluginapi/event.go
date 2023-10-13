//

package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/service/eventbus"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
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

func PublishEventHandler(c *gin.Context) {
	msg := &common.EventMessage{}
	if err := c.ShouldBindQuery(msg); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	eventbus.PublishEvent(msg)
}
