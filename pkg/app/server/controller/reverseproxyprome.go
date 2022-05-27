package controller

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func Query(c *gin.Context) {
	target := "http://" + sconfig.Config().Monitor.PrometheusAddr //转向的host
	remote, err := url.Parse(target)
	if err != nil {
		logger.Info("parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = "api/v1/query" //请求API
	proxy.ServeHTTP(c.Writer, c.Request)
}
func QueryRange(c *gin.Context) {
	target := "http://" + sconfig.Config().Monitor.PrometheusAddr //转向的host
	remote, err := url.Parse(target)
	if err != nil {
		logger.Info("parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = "api/v1/query_range" //请求API
	proxy.ServeHTTP(c.Writer, c.Request)
}
