package network

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network/handlers"
)

func HttpServerStart(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// TODO: 此处绑定 http api handler
	group := engine.Group("/api")
	{
		group.GET("/agent_info", handlers.AgentInfoHandler)
		group.GET("/agent_list", handlers.AgentListHandler)
		group.GET("/run_script", handlers.RunScript)

		group.GET("/os_info", handlers.OSInfoHandler)
	}

	// TODO: 此处绑定前端静态资源handler
	engine.Static("/static", "./dist/static")
	engine.StaticFile("/", "./dist/index.html")

	return engine.Run(addr)
}
