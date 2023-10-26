package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 查询插件清单
func GetPluginsHandler(c *gin.Context) {
	plugins, err := plugin.GetPlugins()
	if err != nil {
		response.Fail(c, nil, "查询插件错误："+err.Error())
		return
	}

	logger.Info("find %d plugins", len(plugins))
	response.Success(c, plugins, "插件查询成功")
}

// 添加插件
func AddPluginHandler(c *gin.Context) {
	param := &plugin.AddPluginParam{}
	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypePlugin, "Add Plugin", user.ID)
	auditlog.Add(log)

	if err := plugin.AddPlugin(param); err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFail)
		response.Fail(c, nil, "add plugin failed:"+err.Error())
		return
	}
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "插件添加成功")
}

// 停用/启动插件
func TogglePluginHandler(c *gin.Context) {
	param := struct {
		UUID   string `json:"uuid"`
		Enable int    `json:"enable"`
	}{}

	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypePlugin, "Toggle Plugin", user.ID)
	auditlog.Add(log)

	logger.Info("toggle plugin:%s to enable %d", param.UUID, param.Enable)
	plugin.TogglePlugin(param.UUID, param.Enable)
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "插件信息更新成功")
}

// 卸载插件
func UnloadPluginHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "undefined" {
		response.Fail(c, nil, "参数错误")
		return
	}
	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypePlugin, "Unload Plugin", user.ID)
	auditlog.Add(log)

	logger.Info("unload plugin:%s", uuid)
	plugin.DeletePlugin(uuid)
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	response.Success(c, nil, "插件信息更新成功")
}

func PluginGatewayHandler(c *gin.Context) {
	// TODO
	name := c.Param("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		c.String(http.StatusNotFound, "plugin not found: "+err.Error())
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := auditlog.New(auditlog.LogTypePlugin, "parse plugin", user.ID)
	auditlog.Add(log)

	s := strings.Replace(p.Url, "/plugin/"+name, "", 1)
	target, err := url.Parse(s)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	logger.Debug("proxy plugin request to: %s", target)
	c.Request.Host = target.Host

	proxy := httputil.NewSingleHostReverseProxy(target)
	auditlog.UpdateStatus(log, auditlog.StatusSuccess)
	proxy.ServeHTTP(c.Writer, c.Request)
}
