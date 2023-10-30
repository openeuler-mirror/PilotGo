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
	"github.com/google/uuid"
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Register Listener",
	}
	auditlog.Add(log)
	eventbus.AddListener(l)

	eventtypes := strings.Split(c.Query("eventTypes"), ",")
	for _, v := range eventtypes {
		log_s := &auditlog.AuditLog{
			LogUUID:    uuid.New().String(),
			ParentUUID: log.LogUUID,
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "Register Listener",
		}
		auditlog.Add(log_s)

		eventtype, err := strconv.Atoi(v)
		if err != nil {
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	eventtypes := strings.Split(c.Query("eventTypes"), ",")
	for _, v := range eventtypes {
		log_s := &auditlog.AuditLog{
			LogUUID:    uuid.New().String(),
			ParentUUID: "",
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "delete eventType",
		}
		auditlog.Add(log_s)

		eventtype, err := strconv.Atoi(v)
		if err != nil {
			auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}
		eventbus.RemoveEventMap(eventtype, l)
	}

	if !eventbus.IsExitEventMap(l) {
		log := &auditlog.AuditLog{
			LogUUID:    uuid.New().String(),
			ParentUUID: "",
			Module:     auditlog.ModulePlugin,
			Status:     auditlog.StatusOK,
			UserID:     u.ID,
			Action:     "delete Listener",
		}
		auditlog.Add(log)
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuid.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "publish Event",
	}
	auditlog.Add(log)

	eventbus.PublishEvent(msg)
	response.Success(c, gin.H{"status": "ok"}, "publishEvent成功")
}
