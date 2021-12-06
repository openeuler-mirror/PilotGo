package network

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network/handlers"
)

func HttpServerStart(addr string) error {
	// TODO: 此处绑定 http api handler
	r := gin.Default()
	r.GET("/api/agent_info", handlers.AgentInfoHandler)

	// TODO: 此处绑定前端静态资源handler

	return r.Run(addr)
}
