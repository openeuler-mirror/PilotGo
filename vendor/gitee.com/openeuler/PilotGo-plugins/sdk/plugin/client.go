package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gitee.com/openeuler/PilotGo-plugins/sdk/utils"
	"github.com/gin-gonic/gin"
)

type Client struct {
	HttpEngine *gin.Engine
	Server     string
	PluginName string
}

var BaseInfo *PluginInfo

func InfoHandler(c *gin.Context) {

	c.JSON(http.StatusOK, BaseInfo)
}

func ReverseProxyHandler(c *gin.Context) {
	remote := c.GetString("__internal__reverse_dest")
	if remote == "" {
		fmt.Println("get reverse dest failed!")
		return
	}

	target, err := url.Parse(remote)
	if err != nil {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/plugin/grafana", "", 1) //请求API

	proxy.ServeHTTP(c.Writer, c.Request)
}

func DefaultClient(desc *PluginInfo) *Client {
	BaseInfo = desc
	dest := desc.ReverseDest

	router := gin.Default()
	mg := router.Group("plugin_manage/")
	{
		mg.GET("/info", InfoHandler)
	}

	pg := router.Group("/plugin/" + desc.Name)
	{
		pg.Any("/*any", func(c *gin.Context) {
			c.Set("__internal__reverse_dest", dest)
			ReverseProxyHandler(c)
		})
	}

	return &Client{
		HttpEngine: router,
	}
}

func (c *Client) Serve(url ...string) {
	// TODO: 启动http服务
	c.HttpEngine.Run(url...)
}

type CmdResult struct {
	MachineUUID string
	MachineIP   string
	Code        int
	Stdout      string
	Stderr      string
}

func (c *Client) RunScript(batch []string, cmd string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_script"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := []*CmdResult{}
	if err := json.Unmarshal(bs, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type MachineNode struct {
	UUID       string
	Department string
	IP         string
	CPUArch    string
	OSInfo     string
	State      int
}

func (c *Client) MachineList() ([]*MachineNode, error) {
	url := c.Server + "/api/v1/pluginapi/machine_list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := []*MachineNode{}
	if err := json.Unmarshal(bs, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type Event struct {
	ID       int
	MetaData interface{}
}

type EventCallback func(e *Event)

func (c *Client) ListenEvent(event Event, callback EventCallback) error {
	url := c.Server + "/api/v1/pluginapi/listener"
	data, err := utils.Request("PUT", url)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	// TODO: register event handler here
	return nil
}

func (c *Client) UnListenEvent(listenerID string) error {
	url := c.Server + "/api/v1/pluginapi/listener"
	data, err := utils.Request("DELETE", url)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	// TODO: unregister event handler here
	return nil
}

func (c *Client) ApplyConfig(batch []string, path, content string) error {
	url := c.Server + "/api/v1/pluginapi/apply_config"
	data, err := utils.Request("PUT", url)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	return nil
}

func (c *Client) GetPluginInfo(name string) (*PluginInfo, error) {
	url := c.Server + "/api/v1/pluginapi/plugins"
	data, err := utils.Request("GET", url)
	if err != nil {
		return nil, err
	}

	resp := &PluginInfo{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
