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
	c.JSON(http.StatusOK, BaseInfo)
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
