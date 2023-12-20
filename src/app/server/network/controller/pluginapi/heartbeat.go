package pluginapi

import (
	"strings"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

var mu sync.Mutex

func PluginHeartbeat(c *gin.Context) {
	ClientID := &struct {
		ClientID string `json:"clientID"`
	}{}
	if err := c.ShouldBindJSON(&ClientID); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	// 更新心跳时间
	key := client.HeartbeatKey + ClientID.ClientID
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
	mu.Lock()
	defer mu.Unlock()

	pluginkeys := redismanager.Scan(client.HeartbeatKey + "*")
	for _, pluginkey := range pluginkeys {
		url_name := strings.Split(pluginkey, client.HeartbeatKey)[1]
		url := strings.Split(url_name, "+")[0]
		var valueObj client.PluginStatus
		plugin_status, err := redismanager.Get(pluginkey, &valueObj)
		if err != nil {
			logger.Error("Error getting %v last heartbeat: %v", url, err)
			continue
		}

		if !plugin_status.(*client.PluginStatus).Connected || time.Since(plugin_status.(*client.PluginStatus).LastConnect) > client.HeartbeatInterval+1 {
			err := plugin.Handshake(url)
			if err != nil {
				logger.Error("rebind plugin and pilotgo server failed:%v", err.Error())
				value := client.PluginStatus{
					Connected:   false,
					LastConnect: plugin_status.(*client.PluginStatus).LastConnect,
				}
				redismanager.Set(pluginkey, value)
			} else {
				value := client.PluginStatus{
					Connected:   true,
					LastConnect: time.Now(),
				}
				redismanager.Set(pluginkey, value)
			}
		}
	}
}
