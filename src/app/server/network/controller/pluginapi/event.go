//

package pluginapi

import (
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/eventbus"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func RegisterListenerHandler(c *gin.Context) {
	p := client.PluginInfo{}
	if err := c.ShouldBind(&p); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	l := &eventbus.Listener{
		Name: p.Name,
		URL:  p.Url,
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	log_s := auditlog.New_sub(auditlog.LogTypePlugin, " Register Listener", "", "", auditlog.StatusSuccess, user.ID)
	auditlog.Add(log_s)
	eventbus.AddListener(l)

	eventtypes := strings.Split(c.Query("eventTypes"), ",")
	for _, v := range eventtypes {
		eventtype, err := strconv.Atoi(v)
		if err != nil {
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}
		log_s := auditlog.New_sub(auditlog.LogTypePlugin, " Register eventType", "", "", auditlog.StatusSuccess, user.ID)
		auditlog.Add(log_s)
		eventbus.AddEventMap(eventtype, l)
	}
	response.Success(c, gin.H{"status": "ok"}, "注册eventType成功")
}

func UnregisterListenerHandler(c *gin.Context) {
	p := client.PluginInfo{}
	if err := c.ShouldBind(&p); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	l := &eventbus.Listener{
		Name: p.Name,
		URL:  p.Url,
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	eventtypes := strings.Split(c.Query("eventTypes"), ",")
	for _, v := range eventtypes {
		eventtype, err := strconv.Atoi(v)
		if err != nil {
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}
		log_s := auditlog.New_sub(auditlog.LogTypePlugin, "delete eventType", "", "", auditlog.StatusSuccess, user.ID)
		auditlog.Add(log_s)
		eventbus.RemoveEventMap(eventtype, l)
	}

	if !eventbus.IsExitEventMap(l) {
		log_s := auditlog.New_sub(auditlog.LogTypePlugin, "delete Listener", "", "", auditlog.StatusSuccess, user.ID)
		auditlog.Add(log_s)
		eventbus.RemoveListener(l)
	}
	response.Success(c, gin.H{"status": "ok"}, "删除eventType成功")
}

func PublishEventHandler(c *gin.Context) {
	msg := &common.EventMessage{}
	if err := c.ShouldBind(msg); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log_s := auditlog.New_sub(auditlog.LogTypePlugin, "publishEvent", "", "", auditlog.StatusSuccess, user.ID)
	auditlog.Add(log_s)

	eventbus.PublishEvent(msg)
	response.Success(c, gin.H{"status": "ok"}, "publishEvent成功")
}
