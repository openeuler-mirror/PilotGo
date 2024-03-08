package plugin

import (
	"time"

	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
)

func CheckPluginHeartbeats() {
	for {
		time.Sleep(client.HeartbeatInterval)
		checkAndRebind()
	}
}

func checkAndRebind() {
	plugins, err := GetPlugins()
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
			err := Handshake(p.Url, p)
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
