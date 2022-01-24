package net

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"openeluer.org/PilotGo/PilotGo/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/plugin"
)

type errorResp struct {
	Status    string `json:"status"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
}

type resultResp struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type pluginGet struct {
	Plugin      string `json:"plugin"`
	Host        string `json:"host"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type pluginAdd struct {
	Plugin      string `json:"plugin"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type hostPut struct {
	Ip               string `json:"ip"`
	SystemInfo       string `json:"system_info"`
	SystemVersion    string `json:"system_version"`
	Arch             string `json:"arch"`
	InstallationTime string `json:"installation_time"`
	MachineType      int    `json:"machine_type"`
}

type hostDelete struct {
	Id []string `json:"id"`
}

var sessionManage SessionManage
var sqlManager *mysqlmanager.MysqlManager

func MakeHandler(name string, f func(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		errResp := &errorResp{
			Status: "error",
		}

		rResp := &resultResp{
			Status: "success",
		}

		code, err := f(c, errResp, rResp)
		if err != nil {
			logger.Error("%s\n", err.Error())
			c.JSON(code, errResp)
			return
		}
		c.JSON(code, rResp)
	}
}

func PluginHandler(c *gin.Context) {
	path := c.Param("any")
	pluginName := strings.Split(path, "/")[1]
	fmt.Println("process plugin request:", path)
	fmt.Println("process plugin request:", pluginName)
	if plugin.GetManager().Check(pluginName) {

		pluginurl, port, protocol, err := plugin.GetManager().Get(pluginName)
		if err != nil {
			c.JSON(http.StatusBadRequest, &errorResp{
				Status:    "error",
				ErrorType: "plugin_error",
				Error:     err.Error(),
			})
			return
		}

		rp := httputil.ReverseProxy{
			Director: func(r *http.Request) {
				// pluginBaseUrl := "http://localhost:8080"
				pluginBaseUrl := protocol + "://" + pluginurl + ":" + port

				newUrl := pluginBaseUrl + strings.Replace(c.Request.URL.Path, "/plugin/"+pluginName, "", -1) + "?" + c.Request.URL.RawQuery

				logger.Trace("new url:", newUrl)
				u, _ := url.Parse(newUrl)
				r.URL = u
			},
		}

		rp.ServeHTTP(c.Writer, c.Request)
	} else {
		c.JSON(http.StatusBadRequest, &errorResp{
			Status:    "error",
			ErrorType: "plugin_error",
			Error:     "plugin not registered",
		})
	}
}

func PluginOpsHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	plugins := plugin.GetManager().GetAll()
	if plugins == nil {
		rResp.Data = "{}"
		return http.StatusOK, nil
	}

	pluginLen := len(plugins)
	pluginget := make([]pluginGet, pluginLen)
	for i := 0; i < pluginLen; i++ {
		pluginget[i].Plugin = plugins[i].Name
		pluginget[i].Host = plugins[i].Url
		pluginget[i].Version = plugins[i].Version
		pluginget[i].Description = plugins[i].Description
		pluginget[i].Status = plugins[i].Status
	}

	rResp.Data = pluginget
	return http.StatusOK, nil
}

func PluginDeleteHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	type pluginName struct {
		Plugin []string `json:"plugin"`
	}

	var plugName pluginName
	err := c.ShouldBindBodyWith(&plugName, binding.JSON)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	err = sqlManager.Delete("name", plugName.Plugin, mysqlmanager.PluginInfo{})
	if err != nil {
		errResp.ErrorType = "db_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	for _, value := range plugName.Plugin {
		plugin.GetManager().Remove(value)
	}

	rResp.Data = "{}"
	return http.StatusOK, nil
}

func PluginAddHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	var pluAdd pluginAdd
	err := c.ShouldBindBodyWith(&pluAdd, binding.JSON)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	err = getPluginAbout(&pluAdd)
	if err != nil {
		errResp.ErrorType = "plugin_about_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	err = plugin.GetManager().Regist(&plugin.Plugin{
		Name:        pluAdd.Plugin,
		Version:     pluAdd.Version,
		Description: pluAdd.Description,
		Url:         pluAdd.Host,
		Port:        pluAdd.Port,
		Protocol:    pluAdd.Protocol,
		Status:      "registered",
	})
	if err != nil {
		errResp.ErrorType = "regist_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, errors.New(fmt.Sprintf("can not regster %s,error:%s", pluAdd.Plugin, err.Error()))
	}

	err = sqlManager.Insert(&mysqlmanager.PluginInfo{
		Name:        pluAdd.Plugin,
		Description: pluAdd.Description,
		Url:         pluAdd.Host,
		Port:        pluAdd.Port,
		Protocol:    pluAdd.Protocol,
		Version:     pluAdd.Version,
	}, 1)

	if err != nil {
		plugin.GetManager().Remove(pluAdd.Plugin)
		errResp.ErrorType = "db_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	rResp.Data = "{}"
	return http.StatusOK, nil
}

func HostsGetHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	mi, err := mysqlmanager.GetMachInfo(sqlManager)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	rResp.Data = mi
	return http.StatusOK, nil
}

func HostsOverview(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	type HostView struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	mi, err := mysqlmanager.GetMachInfo(sqlManager)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}
	arrayLen := len(mi)
	i := 0
	for _, m := range mi {
		if m.MachineType == 0 {
			i++
		}
	}

	hostView := [2]HostView{
		{Name: "physics", Value: i},
		{Name: "virtual", Value: arrayLen - i},
	}
	rResp.Data = hostView
	return http.StatusOK, nil
}
func HostAddHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	var info hostPut
	err := c.ShouldBindBodyWith(&info, binding.JSON)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	instalTime, err := time.ParseInLocation("2006-01-02 15:04:05", info.InstallationTime, time.Local)
	if err != nil {
		errResp.ErrorType = "time_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	err = sqlManager.Insert(&mysqlmanager.MachInfo{
		Ip:               info.Ip,
		SystemStatus:     0,
		SystemInfo:       info.SystemInfo,
		SystemVersion:    info.SystemVersion,
		Arch:             info.Arch,
		InstallationTime: instalTime,
		MachineType:      info.MachineType,
	}, 1)
	if err != nil {
		errResp.ErrorType = "insert_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, errors.New(fmt.Sprintf("insert error:%s", err.Error()))
	}

	rResp.Data = "{}"
	return http.StatusOK, nil
}

func HostDeleteHandler(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {
	var value hostDelete
	err := c.ShouldBindBodyWith(&value, binding.JSON)
	if err != nil {
		errResp.ErrorType = "json_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, err
	}

	for _, val := range value.Id {
		if val == "0" {
			err = errors.New(fmt.Sprintf("invaild value:%d", value.Id))
			errResp.ErrorType = "value_error"
			errResp.Error = err.Error()
			return http.StatusBadRequest, err
		}
	}

	err = sqlManager.Delete("id", value.Id, mysqlmanager.MachInfo{})
	if err != nil {
		errResp.ErrorType = "delete_error"
		errResp.Error = err.Error()
		return http.StatusBadRequest, errors.New(fmt.Sprintf("delete error:%s", err.Error()))
	}

	rResp.Data = "{}"
	return http.StatusOK, nil
}

func GetLogin(c *gin.Context, errResp *errorResp, rResp *resultResp) (int, error) {

	u, p, ok := c.Request.BasicAuth()
	if !ok {
		errResp.ErrorType = "BasicAuth_error"
		errResp.Error = "BasicAuth_error"
		return http.StatusUnauthorized, errors.New("BasicAuth failed")
	}
	ok = CheckAuth(u, p)
	if !ok {
		errResp.ErrorType = "checkauth_error"
		errResp.Error = "checkauth_error"
		return http.StatusUnauthorized, errors.New(fmt.Sprintf("CheckAuth %s:%s failed", u, p))
	}

	id, err := c.Cookie("pilotgoSession")
	isfind := false
	if len(id) > 0 {
		isfind = sessionManage.FindAndFlush(id)
	}

	//创建session id
	if !isfind || err != nil {
		sessionId := CreateSessionId()
		err = sessionManage.Set(sessionId, &SessionInfo{
			sessionTime: time.Now(),
		})
		if err != nil {
			errResp.ErrorType = "session_create_error"
			errResp.Error = err.Error()
			return http.StatusUnauthorized, errors.New("session_create_error")
		}
		c.SetCookie("pilotgoSession", sessionId, sessionManage.maxAge, "/", "", false, true)
	}

	rResp.Data = "{}"
	return http.StatusOK, nil
}

func checkSession(c *gin.Context) {
	if c == nil {
		logger.Error("c == nil")
		return
	}

	errResp := &errorResp{
		Status: "error",
	}

	defer func() {
		if errResp.Status == "error" {
			c.JSON(http.StatusUnauthorized, errResp)
			c.Abort()
		}
	}()

	id, err := c.Cookie("pilotgoSession")
	if err != nil {
		logger.Error("get cookie failed!")
		errResp.ErrorType = "cookie_error"
		errResp.Error = "cookie_error"
		return
	}

	if isfind := sessionManage.FindAndFlush(id); !isfind {
		logger.Error("find session %s failed!", id)
		errResp.ErrorType = "session_error"
		errResp.Error = "session_error"
		return
	}
	errResp.Status = ""
	c.Next()
}

func Start(conf *config.Configure) (err error) {
	sqlManager, err = mysqlmanager.Init(conf.Dbinfo.HostName, conf.Dbinfo.UserName, conf.Dbinfo.Password, conf.Dbinfo.DataBase, conf.Dbinfo.Port)
	if err != nil {
		return err
	}

	sessionManage.Init(conf.MaxAge, conf.SessionCount)
	go func() {
		for true {
			time.Sleep(time.Second * 10)
			//每10秒读取一次数据库，并更改数据库状态
			mi, err := mysqlmanager.GetMachInfo(sqlManager)
			if err != nil {
				continue
			}

			for _, m := range mi {
				status := m.CheckStatus()
				if m.SystemStatus != status {
					m1 := mysqlmanager.MachInfo{
						Id:           m.Id,
						SystemStatus: status,
					}
					sqlManager.Update(&m1)
				}
			}
		}
	}()

	pi, err := mysqlmanager.GetPluginInfo(sqlManager)
	if err != nil {
		return err
	}

	for _, value := range pi {
		plugin.GetManager().Regist(&plugin.Plugin{
			Name:        value.Name,
			Version:     value.Version,
			Description: value.Description,
			Url:         value.Url,
			Port:        value.Port,
			Protocol:    value.Protocol,
		})
	}

	r := gin.Default()
	r.LoadHTMLFiles("./static/index.html")
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("/login", MakeHandler("getLogin", GetLogin))
	r.Static("/static", "./static")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	//注册session校验中间件
	//r.Use(checkSession)

	// PilotGo server端plugin处理
	r.GET("/plugin", MakeHandler("pluginOpsHandler", PluginOpsHandler))
	r.DELETE("/plugin", MakeHandler("pluginDeleteHandler", PluginDeleteHandler))
	r.POST("/plugin", MakeHandler("pluginPutHandler", PluginAddHandler))

	// 转发到plugin server端处理
	r.GET("/plugin/*any", PluginHandler)
	//获取机器列表
	r.GET("/hosts", MakeHandler("hostGetHandler", HostsGetHandler))
	r.POST("/hosts", MakeHandler("hostPutHandler", HostAddHandler))
	r.DELETE("/hosts", MakeHandler("hostDeleteHandler", HostDeleteHandler))
	r.GET("/overview", MakeHandler("overview", HostsOverview))
	server_url := ":" + strconv.Itoa(conf.S.ServerPort)
	e := r.RunTLS(server_url, "server.crt", "server.key")
	return e
}

func getPluginAbout(pluAdd *pluginAdd) error {
	temp := pluginAdd{}
	//从插件中获取信息
	url := pluAdd.Protocol + "://" + pluAdd.Host + ":" + pluAdd.Port + "/plugin/about"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return err
	}

	pluAdd.Plugin = temp.Plugin
	pluAdd.Version = temp.Version
	pluAdd.Description = temp.Description

	return nil
}
