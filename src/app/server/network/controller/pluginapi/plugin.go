package pluginapi

import (
	"time"

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func PluginList(c *gin.Context) {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		response.Fail(c, nil, "查询插件错误："+err.Error())
		return
	}
	response.Success(c, plugins, "插件查询成功")
}

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

func HasPermission(c *gin.Context) {
	p := &common.Permission{}
	if err := c.ShouldBind(&p); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	claims, err := jwt.ParseMyClaims(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	logger.Debug("request from %d, %s", claims.UserId, claims.UserName)
	ok, err := auth.CheckAuth(claims.UserName, p.Resource, p.Operate)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, ok, "include permission")
}
