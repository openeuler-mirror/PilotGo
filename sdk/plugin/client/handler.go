package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

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

func InfoHandler(c *gin.Context) {
	v, ok := c.Get("__internal__client_instance")
	if !ok {
		response.Fail(c, gin.H{"status": false}, "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		response.Fail(c, gin.H{"status": false}, "client信息错误")
		return
	}

	info := &PluginFullInfo{
		PluginInfo: *client.PluginInfo,
		Extentions: client.extentions,
	}

	c.JSON(http.StatusOK, info)
}

func BindHandler(c *gin.Context) {
	v, ok := c.Get("__internal__client_instance")
	if !ok {
		response.Fail(c, gin.H{"status": false}, "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		response.Fail(c, gin.H{"status": false}, "client信息错误")
		return
	}
	client.server = c.Request.RemoteAddr

	response.Success(c, nil, "bind server success")
}

func EventHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		return
	}
	var msg common.EventMessage
	if err := json.Unmarshal(j, &msg); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		return
	}

	v, ok := c.Get("__internal__client_instance")
	if !ok {
		return
	}
	client, ok := v.(*Client)
	if !ok {
		return
	}

	client.ProcessEvent(&msg)
}

func ExtentionsHandler(c *gin.Context) {
	v, ok := c.Get("__internal__client_instance")
	if !ok {
		response.Fail(c, gin.H{"status": false}, "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		response.Fail(c, gin.H{"status": false}, "client信息错误")
		return
	}

	response.Success(c, client.extentions, "")
}

func CommandResultHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		return
	}
	var result common.AsyncCmdResult
	if err := json.Unmarshal(j, &result); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		return
	}

	v, ok := c.Get("__internal__client_instance")
	if !ok {
		logger.Error("%v", "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		logger.Error("%v", "client获取失败")
		return
	}

	client.ProcessCommandResult(&result)

}

func TagsHandler(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body) // 接收数据
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		response.Fail(c, gin.H{"status": false}, "没获取到："+err.Error())
		return
	}
	uuidTags := &struct {
		UUIDS []string `json:"uuids"`
	}{}
	if err := json.Unmarshal(j, &uuidTags); err != nil {
		logger.Error("反序列化结果失败%s", err.Error())
		response.Fail(c, gin.H{"status": false}, "反序列化结果失败："+err.Error())
		return
	}

	v, ok := c.Get("__internal__client_instance")
	if !ok {
		logger.Error("%v", "未获取到client值信息")
		response.Fail(c, gin.H{"status": false}, "未获取到client值信息")
		return
	}
	client, ok := v.(*Client)
	if !ok {
		logger.Error("%v", "client获取失败")
		response.Fail(c, gin.H{"status": false}, "client获取失败")
		return
	}

	if client.getTagsCallback != nil {
		result := client.getTagsCallback(uuidTags.UUIDS)
		response.Success(c, result, "")
	} else {
		logger.Error("get tags callback not set")
		response.Fail(c, gin.H{"status": false}, "get tags callback not set")
	}
}
