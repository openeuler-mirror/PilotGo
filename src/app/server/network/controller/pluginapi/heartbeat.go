package pluginapi

import (
	"time"

	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginHeartbeat(c *gin.Context) {
	PluginID := &struct {
		PluginUrl string `json:"clientID"`
	}{}
	if err := c.ShouldBindJSON(&PluginID); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 更新心跳时间
	key := client.HeartbeatKey + PluginID.PluginUrl
	value := client.PluginStatus{
		Connected:   true,
		LastConnect: time.Now(),
	}
	err := redismanager.Set(key, value)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "Heartbeat received")
}

func CheckPluginHeartbeats() {
	for {
		time.Sleep(client.HeartbeatInterval)
		checkAndRebind()
	}
}

func checkAndRebind() {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		logger.Error("get plugins failed:%v", err.Error())
	}
	for _, p := range plugins {
		key := client.HeartbeatKey + p.Url
		plugin_status, err := redismanager.Get(key, &client.PluginStatus{})
		if err != nil {
			logger.Error("Error getting %v last heartbeat: %v", p.Url, err)
			continue
		}
		if !plugin_status.(*client.PluginStatus).Connected || time.Since(plugin_status.(*client.PluginStatus).LastConnect) > client.HeartbeatInterval+1*time.Second {
			err := plugin.Handshake(p.Url)
			if err != nil {
				logger.Error("rebind plugin and pilotgo server failed:%v", err.Error())
				value := client.PluginStatus{
					Connected:   false,
					LastConnect: plugin_status.(*client.PluginStatus).LastConnect,
				}
				redismanager.Set(key, value)
			} else {
				value := client.PluginStatus{
					Connected:   true,
					LastConnect: time.Now(),
				}
				redismanager.Set(key, value)
			}
		}
	}
}
