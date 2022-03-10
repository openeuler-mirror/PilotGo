/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-11-18 13:03:16
 * LastEditTime: 2022-03-10 13:39:14
 * Description: Interface routing forwarding
 ******************************************************************************/
package router

import (
	"net/http"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network/handlers"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/middleware"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.LoggerToFile())
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
		group.GET("/sysctl_view", handlers.SysctlViewHandler)
		group.GET("/service_list", handlers.ServiceListHandler)
		group.GET("/service_status", handlers.ServiceStatusHandler)
		group.GET("/rpm_all", handlers.AllRpmHandler)
		group.GET("/rpm_source", handlers.RpmSourceHandler)
		group.GET("/rpm_info", handlers.RpmInfoHandler)
		group.GET("/disk_use", handlers.DiskUsageHandler)
		group.GET("/disk_info", handlers.DiskInfoHandler)
		group.GET("/net_tcp", handlers.NetTCPHandler)
		group.GET("/net_udp", handlers.NetUDPHandler)
		group.GET("/net_io", handlers.NetIOCounterHandler)
		group.GET("/net_nic", handlers.NetNICConfigHandler)
		group.GET("/user_info", handlers.CurrentUserInfoHandler)
		group.GET("/user_all", handlers.AllUserInfoHandler)
	}
	agent := router.Group("/agent")
	{
		agent.GET("/sysctl_change", handlers.SysctlChangeHandler)
		agent.GET("/service_stop", handlers.ServiceStopHandler)
		agent.GET("/service_start", handlers.ServiceStartHandler)
		agent.GET("/service_restart", handlers.ServiceRestartHandler)
		agent.POST("/rpm_install", handlers.InstallRpmHandler)
		agent.POST("/rpm_remove", handlers.RemoveRpmHandler)
		agent.GET("/disk_path", handlers.DiskCreatPathHandler)
		agent.GET("/disk_mount", handlers.DiskMountHandler)
		agent.GET("/disk_umount", handlers.DiskUMountHandler)
		agent.GET("/disk_format", handlers.DiskFormatHandler)
		agent.GET("/user_add", handlers.AddLinuxUserHandler)
		agent.GET("/user_del", handlers.DelUserHandler)
		agent.GET("/user_ower", handlers.ChangeFileOwnerHandler)
		agent.GET("/user_per", handlers.ChangePermissionHandler)
		agent.GET("/log_all", controller.LogAll)
		agent.GET("/logs", controller.AgentLogs)
		agent.POST("/delete", controller.DeleteLog)
	}
	user := router.Group("user")
	{
		user.POST("/login", controller.Login)
		user.GET("/logout", controller.Logout)
		user.GET("/info", middleware.AuthMiddleware(), controller.Info)
	}
	machinemanager := router.Group("machinemanager")
	{
		machinemanager.GET("/departinfo", controller.DepartInfo)
		machinemanager.GET("/machineinfo", controller.MachineInfo)
		machinemanager.GET("/depart", controller.Dep)
		machinemanager.GET("/test", controller.AddIP)
	}
	batchmanager := router.Group("batchmanager")
	{
		batchmanager.POST("/createbatch", controller.CreateBatch)
		batchmanager.GET("/batchinfo", controller.BatchInform)
		batchmanager.POST("/batchmachineinfo", controller.Batchmachineinfo)
	}
	prometheus := router.Group("prometheus")
	{
		prometheus.POST("/queryrange", controller.Queryrange)
		prometheus.POST("/query", controller.Query)
	}
	a := gormadapter.NewAdapter("mysql", mysqlmanager.Url, true)
	common.E = casbin.NewEnforcer("./rbac_models.conf", a)
	common.E.LoadPolicy()
	Level := router.Group("")
	Level.Use(common.CasbinHandler())
	{
		user.POST("/register", controller.Register)
		user.GET("/searchAll", controller.UserAll)
		user.POST("/userSearch", controller.UserSearch)
		user.GET("/reset", controller.ResetPassword)
		user.POST("/delete", controller.DeleteUser)
		user.POST("/update", controller.UpdateUser)
		user.POST("/import", controller.ImportUser)
		machinemanager.GET("/t", controller.Deletedepartdata)
		machinemanager.POST("/deletemachinedata", controller.Deletemachinedata)
		machinemanager.POST("/adddepart", controller.AddDepart)
		machinemanager.GET("/updatedepart", controller.UpdateDepart)
		batchmanager.POST("/updatebatch", controller.UpdateBatch)
		batchmanager.POST("/deletebatch", controller.DeleteBatch)
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
