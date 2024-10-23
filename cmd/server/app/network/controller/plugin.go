package controller

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/common"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
	"github.com/gorilla/websocket"
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

// 分页查询插件清单
func GetPluginsPagedHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := plugin.GetPluginPaged(num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

// 添加插件
func AddPluginHandler(c *gin.Context) {
	param := plugin.PluginParam{}
	if err := c.BindJSON(&param); err != nil {
		response.Fail(c, nil, "参数错误")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Add Plugin",
	}
	auditlog.Add(log)

	if err := plugin.AddPlugin(&param); err != nil {
		auditlog.UpdateStatus(log, auditlog.StatusFailed)
		response.Fail(c, nil, "add plugin failed:"+err.Error())
		return
	}
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Toggle Plugin",
	}
	auditlog.Add(log)

	logger.Info("toggle plugin:%s to enable %d", param.UUID, param.Enable)
	if err = plugin.TogglePlugin(param.UUID, param.Enable); err != nil {
		response.Fail(c, nil, "toggle plugin error:"+err.Error())
		return
	}
	//根据uuid从pluginmanage中获取本插件信息，并返回插件扩展点信息
	plugin, err := plugin.GetPluginByUUID(param.UUID)
	if err != nil {
		response.Fail(c, nil, "get plugin by uuid error:"+err.Error())
	}
	response.Success(c, plugin.Extentions, "插件信息更新成功")
}

// 卸载插件
func UnloadPluginHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "undefined" {
		response.Fail(c, nil, "参数错误")
		return
	}
	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Unload Plugin",
	}
	auditlog.Add(log)

	logger.Info("unload plugin:%s", uuid)
	if err := plugin.DeletePlugin(uuid); err != nil {
		return
	}
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "parse Plugin",
	}
	auditlog.Add(log)

	s := strings.Replace(p.Url, "/plugin/"+name, "", 1)
	ishttp, err := httputils.ServerIsHttp(s)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	if ishttp && strings.Split(s, "://")[0] == "https" {
		s = "http://" + strings.Split(s, "://")[1]
	}
	if !ishttp && strings.Split(s, "://")[0] == "http" {
		s = "https://" + strings.Split(s, "://")[1]
	}

	target, err := url.Parse(s)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	logger.Debug("proxy plugin request to: %s", target)
	c.Request.Host = target.Host

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func PluginWebsocketGatewayHandler(c *gin.Context) {
	name := c.Param("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		c.String(http.StatusNotFound, "plugin not found: "+err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "parse Plugin",
	}
	auditlog.Add(log)

	target_addr := strings.Replace(strings.Split(p.Url, "//")[1], "/plugin/"+name, "", 1)
	targetURL_str := fmt.Sprintf("ws://%s/ws/proxy", target_addr)
	ishttp, err := httputils.ServerIsHttp("http://" + target_addr)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	if ishttp && strings.Split(targetURL_str, "://")[0] == "wss" {
		targetURL_str = "ws://" + strings.Split(targetURL_str, "://")[1]
	}
	if !ishttp && strings.Split(targetURL_str, "://")[0] == "ws" {
		targetURL_str = "wss://" + strings.Split(targetURL_str, "://")[1]
	}

	logger.Debug("proxy plugin request to: %s", targetURL_str)

	dialer := websocket.Dialer{
		Proxy: nil,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		HandshakeTimeout: 10 * time.Second,
	}
	target_wsconn, _, err := dialer.Dial(targetURL_str, nil)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("dial to target WebSocket failed: %s", err.Error()))
		return
	}
	defer target_wsconn.Close()

	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 10 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	client_wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("update client WebSocket failed: %s", err.Error()))
		return
	}
	defer client_wsconn.Close()

	transferMsg := func(_srcConn, _destConn *websocket.Conn, _err_ch chan error) {
		for {
			messageType, message, err := _srcConn.ReadMessage()
			if err != nil {
				_err_ch <- fmt.Errorf("error while reading message: %s", err.Error())
				return
			}

			err = _destConn.WriteMessage(messageType, message)
			if err != nil {
				_err_ch <- fmt.Errorf("error while writing message: %s", err.Error())
				return
			}
		}
	}
	client2TargetErrChan, target2ClientErrChan := make(chan error, 1), make(chan error, 1)

	go transferMsg(client_wsconn, target_wsconn, client2TargetErrChan)
	go transferMsg(target_wsconn, client_wsconn, target2ClientErrChan)
	select {
	case err := <-client2TargetErrChan:
		logger.Error(err.Error())
		c.String(http.StatusBadGateway, err.Error())
		return
	case err := <-target2ClientErrChan:
		logger.Error(err.Error())
		c.String(http.StatusBadGateway, err.Error())
		return
	}
}
