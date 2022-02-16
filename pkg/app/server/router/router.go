package router

/**
 * @Author: zhang han
 * @Date: 2021/10/28 09:58
 * @Description: 接口路由转发
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network/handlers"
	"openeluer.org/PilotGo/PilotGo/pkg/common/middleware"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// TODO: 此处绑定 http api handler

	group := router.Group("/api")
	{
		group.GET("/agent_info", handlers.AgentInfoHandler)
		group.GET("/agent_list", handlers.AgentListHandler)
		group.GET("/run_script", handlers.RunScript)
		group.GET("/os_info", handlers.OSInfoHandler)
		group.GET("/cpu_info", handlers.CPUInfoHandler)
		group.GET("/memory_info", handlers.MemoryInfoHandler)
		group.GET("/sysctl_info", handlers.SysInfoHandler)
		group.GET("/sysctl_change", handlers.SysctlChangeHandler)
		group.GET("/sysctl_view", handlers.SysctlViewHandler)
		group.GET("/service_list", handlers.ServiceListHandler)
		group.GET("/service_status", handlers.ServiceStatusHandler)
		group.GET("/service_stop", handlers.ServiceStopHandler)
		group.GET("/service_start", handlers.ServiceStartHandler)
		group.GET("/service_restart", handlers.ServiceRestartHandler)
		group.GET("/rpm_all", handlers.AllRpmHandler)
		group.GET("/rpm_source", handlers.RpmSourceHandler)
		group.GET("/rpm_info", handlers.RpmInfoHandler)
		group.GET("/rpm_install", handlers.InstallRpmHandler)
		group.GET("/rpm_remove", handlers.RemoveRpmHandler)
		group.GET("/disk_use", handlers.DiskUsageHandler)
		group.GET("/disk_info", handlers.DiskInfoHandler)
		group.GET("/disk_path", handlers.DiskCreatPathHandler)
		group.GET("/disk_mount", handlers.DiskMountHandler)
		group.GET("/disk_umount", handlers.DiskUMountHandler)
		group.GET("/disk_format", handlers.DiskFormatHandler)
		group.GET("/net_tcp", handlers.NetTCPHandler)
		group.GET("/net_udp", handlers.NetUDPHandler)
		group.GET("/net_io", handlers.NetIOCounterHandler)
		group.GET("/net_nic", handlers.NetNICConfigHandler)
		group.GET("/user_info", handlers.CurrentUserInfoHandler)
		group.GET("/user_all", handlers.AllUserInfoHandler)
		group.GET("/user_add", handlers.AddLinuxUserHandler)
		group.GET("/user_del", handlers.DelUserHandler)
		group.GET("/user_ower", handlers.ChangeFileOwnerHandler)
		group.GET("/user_per", handlers.ChangePermissionHandler)

	}

	user := router.Group("user")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
		user.GET("/info", middleware.AuthMiddleware(), controller.Info)
		user.GET("/searchAll", controller.UserAll)
		user.GET("/reset", controller.ResetPassword)
		user.POST("/delete", controller.DeleteUser)
		user.POST("/update", controller.UpdateUser)
		user.POST("/import", controller.ImportUser)
	}
	machinemanager := router.Group("machinemanager")
	{
		machinemanager.POST("/adddepart", controller.AddDepart)
		machinemanager.POST("/addmachine", controller.AddMachine)
		machinemanager.GET("/departinfo", controller.DepartInfo)
		machinemanager.GET("/machineinfo", controller.MachineInfo)
		machinemanager.POST("/deletemachinedata", controller.Deletemachinedata)
		machinemanager.GET("/t", controller.Deletedepartdata)
		machinemanager.GET("/depart", controller.Dep)
		machinemanager.GET("/updatedepart", controller.UpdateDepart)
		machinemanager.GET("/test", controller.AddIP)
	}
	batchmanager := router.Group("batchmanager")
	{
		batchmanager.POST("/createbatch", controller.CreateBatch)
		batchmanager.POST("/deletebatch", controller.DeleteBatch)
		batchmanager.GET("/batchinfo", controller.BatchInform)
	}
	// TODO: 此处绑定前端静态资源handler
	router.Static("/static", "./dist/static")
	router.StaticFile("/", "./dist/index.html")

	// firewall := router.Group("firewall")
	// {
	// 	firewall.POST("/config", controller.Config)
	// 	firewall.POST("/stop", controller.Stop)
	// 	firewall.POST("/restart", controller.Restart)
	// 	firewall.POST("/reload", controller.Reload)
	// 	firewall.POST("/addzp", controller.AddZonePort)
	// 	firewall.POST("/delzp", controller.DelZonePort)
	// 	firewall.POST("/addzpp", controller.AddZonePortPermanent)
	// }

	//router.LoadHTMLFiles("./static/index.html")
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	//router.POST("/login", net.MakeHandler("getLogin", net.GetLogin))
	//router.Static("/static", "./static")
	//router.GET("/", func(context *gin.Context) {
	//	context.HTML(http.StatusOK, "index.html", nil)
	//})
	////注册session校验中间件
	////r.Use(checkSession)
	//
	//// PilotGo server端plugin处理
	//router.GET("/plugin", net.MakeHandler("pluginOpsHandler", net.PluginOpsHandler))
	//router.DELETE("/plugin", net.MakeHandler("pluginDeleteHandler", net.PluginDeleteHandler))
	//router.POST("/plugin", net.MakeHandler("pluginPutHandler", net.PluginAddHandler))
	//
	//// 转发到plugin server端处理
	//router.GET("/plugin/*any", net.PluginHandler)
	////获取机器列表
	//router.GET("/hosts", net.MakeHandler("hostGetHandler", net.HostsGetHandler))
	//router.POST("/hosts", net.MakeHandler("hostPutHandler", net.HostAddHandler))
	//router.DELETE("/hosts", net.MakeHandler("hostDeleteHandler", net.HostDeleteHandler))
	//router.GET("/overview", net.MakeHandler("overview", net.HostsOverview))

	return router
}
