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

var (
	mu           sync.Mutex
	heartbeatKey = "heartbeat:"
)

type HeartbeatTime struct {
	HeartbeatTime string
}

func PluginHeartbeat(c *gin.Context) {
	ClientID := &struct {
		ClientID string `json:"clientID"`
	}{}
	if err := c.ShouldBindJSON(&ClientID); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	key := heartbeatKey + ClientID.ClientID
	value := HeartbeatTime{
		HeartbeatTime: time.Now().Format(time.RFC3339),
	}
	err := redismanager.Set(key, value)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	logger.Discard()
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

	pluginkeys := redismanager.Scan(heartbeatKey + "*")
	if len(pluginkeys) != 0 {
		for _, pluginkey := range pluginkeys {
			url_name := strings.Split(pluginkey, heartbeatKey)[1]
			url := strings.Split(url_name, "+")[0]
			var valueObj HeartbeatTime
			lastHeartbeatStr, err := redismanager.Get(pluginkey, &valueObj)
			if err != nil {
				logger.Error("Error getting %v last heartbeat: %v", url, err)
				continue
			}
			lastHeartbeat, err := time.Parse(time.RFC3339, lastHeartbeatStr.(*HeartbeatTime).HeartbeatTime)
			if err != nil {
				logger.Error("Error parse %v last heartbeat: %v", url, err)
				continue
			}

			if time.Since(lastHeartbeat) > client.HeartbeatInterval {
				err := plugin.Handshake(url)
				if err != nil {
					logger.Error("rebind plugin and pilotgo server failed:%v", err.Error())
				}
			}
		}
	}
}
