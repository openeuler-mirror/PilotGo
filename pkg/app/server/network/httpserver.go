package network

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network/handlers"
)

func HttpServerStart(addr string) error {
	r := gin.Default()
	// TODO: 此处绑定 http api handler
	r.GET("/api/agent_info", handlers.AgentInfoHandler)

	// TODO: 此处绑定前端静态资源handler
	r.Static("/static", "./dist/static")
	r.StaticFile("/", "./dist/index.html")

	return r.Run(addr)
}
