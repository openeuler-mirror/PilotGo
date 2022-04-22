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
 * LastEditTime: 2022-04-12 14:10:23
 * Description: Interface routing forwarding
 ******************************************************************************/
package router

import (
	"net/http"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager/agentcontroller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service/middleware"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.LoggerToFile())
	router.Use(gin.Recovery())

	// TODO: 此处绑定 http api handler

	group := router.Group("/api")
	{
		group.GET("/agent_info", agentcontroller.AgentInfoHandler)
		group.GET("/agent_list", agentcontroller.AgentListHandler)
		group.GET("/run_script", agentcontroller.RunScript)
		group.GET("/os_info", agentcontroller.OSInfoHandler)
		group.GET("/cpu_info", agentcontroller.CPUInfoHandler)
		group.GET("/memory_info", agentcontroller.MemoryInfoHandler)
		group.GET("/sysctl_info", agentcontroller.SysInfoHandler)
		group.GET("/sysctl_view", agentcontroller.SysctlViewHandler)
		group.GET("/service_list", agentcontroller.ServiceListHandler)
		group.GET("/service_status", agentcontroller.ServiceStatusHandler)
		group.GET("/rpm_all", agentcontroller.AllRpmHandler)
		group.GET("/rpm_source", agentcontroller.RpmSourceHandler)
		group.GET("/rpm_info", agentcontroller.RpmInfoHandler)
		group.GET("/disk_use", agentcontroller.DiskUsageHandler)
		group.GET("/disk_info", agentcontroller.DiskInfoHandler)
		group.GET("/net_tcp", agentcontroller.NetTCPHandler)
		group.GET("/net_udp", agentcontroller.NetUDPHandler)
		group.GET("/net_io", agentcontroller.NetIOCounterHandler)
		group.GET("/net_nic", agentcontroller.NetNICConfigHandler)
		group.GET("/user_info", agentcontroller.CurrentUserInfoHandler)
		group.GET("/user_all", agentcontroller.AllUserInfoHandler)
		group.GET("/os_basic", agentcontroller.OsBasic)
		group.GET("/firewall_config", agentcontroller.FirewalldConfig)
		group.GET("/firewall_restart", agentcontroller.FirewalldRestart)
		group.GET("/firewall_stop", agentcontroller.FirewalldStop)
		group.POST("/firewall_addzp", agentcontroller.FirewalldZonePortAdd)
		group.POST("/firewall_delzp", agentcontroller.FirewalldZonePortDel)
	}
	cluster := router.Group("/cluster")
	{
		cluster.GET("/info", controller.ClusterInfo)
		cluster.GET("/depart_info", controller.DepartClusterInfo)
	}
	agent := router.Group("/agent")
	{
		agent.GET("/sysctl_change", agentcontroller.SysctlChangeHandler)
		agent.POST("/service_stop", agentcontroller.ServiceStopHandler)
		agent.POST("/service_start", agentcontroller.ServiceStartHandler)
		agent.POST("/service_restart", agentcontroller.ServiceRestartHandler)
		agent.POST("/rpm_install", agentcontroller.InstallRpmHandler)
		agent.POST("/rpm_remove", agentcontroller.RemoveRpmHandler)
		agent.GET("/disk_path", agentcontroller.DiskCreatPathHandler)
		agent.GET("/disk_mount", agentcontroller.DiskMountHandler)
		agent.GET("/disk_umount", agentcontroller.DiskUMountHandler)
		agent.GET("/disk_format", agentcontroller.DiskFormatHandler)
		agent.GET("/user_add", agentcontroller.AddLinuxUserHandler)
		agent.GET("/user_del", agentcontroller.DelUserHandler)
		agent.GET("/user_ower", agentcontroller.ChangeFileOwnerHandler)
		agent.GET("/user_per", agentcontroller.ChangePermissionHandler)
		agent.GET("/log_all", controller.LogAll)
		agent.GET("/logs", controller.AgentLogs)
		agent.POST("/delete", controller.DeleteLog)
	}
	user := router.Group("user")
	{
		user.POST("/login", controller.Login)
		user.GET("/logout", controller.Logout)
		user.GET("/searchAll", controller.UserAll)
		user.POST("/userSearch", controller.UserSearch)
		user.GET("/info", middleware.AuthMiddleware(), controller.Info)
		user.POST("/permission", controller.GetLoginUserPermission)
		user.GET("/roles", controller.GetRoles)
		user.GET("/role", controller.GetUserRole)
		user.POST("/addRole", controller.AddUserType)
		user.POST("/delRole", controller.DeleteUserRole)
		user.POST("/updateRole", controller.UpdateUserRole)
		user.POST("/roleChange", controller.RolePermissionChange)
	}
	machinemanager := router.Group("machinemanager")
	{
		machinemanager.GET("/departinfo", controller.DepartInfo)
		machinemanager.GET("/machineinfo", controller.MachineInfo)
		machinemanager.GET("/depart", controller.Dep)
		machinemanager.GET("/test", controller.AddIP)
		machinemanager.GET("/machinealldata", controller.MachineAllData)
		machinemanager.POST("/modifydepart", controller.ModifyMachineDepart)
		machinemanager.GET("/sourcepool", controller.FreeMachineSource)
	}
	batchmanager := router.Group("batchmanager")
	{
		batchmanager.POST("/createbatch", controller.CreateBatch)
		batchmanager.GET("/batchinfo", controller.BatchInform)
		batchmanager.GET("/batchmachineinfo", controller.Batchmachineinfo)
	}
	prometheus := router.Group("prometheus")
	{
		prometheus.POST("/queryrange", controller.Queryrange)
		prometheus.POST("/query", controller.Query)
		prometheus.GET("/alert", controller.ListenALert)
		prometheus.POST("/alertmanager", controller.AlertMessageConfig)
	}
	policy := router.Group("casbin")
	{
		policy.GET("/get", controller.GetPolicy)
		policy.POST("/delete", controller.PolicyDelete)
		policy.POST("/add", controller.PolicyAdd)
	}
	a := gormadapter.NewAdapter("mysql", mysqlmanager.Url, true)
	service.E = casbin.NewEnforcer("./rbac_models.conf", a)
	service.E.LoadPolicy()
	Level := router.Group("")
	Level.Use(service.CasbinHandler())
	{
		user.POST("/register", controller.Register)
		user.POST("/reset", controller.ResetPassword)
		user.POST("/delete", controller.DeleteUser)
		user.POST("/update", controller.UpdateUser)
		user.POST("/import", controller.ImportUser)
		machinemanager.POST("/deletedepartdata", controller.Deletedepartdata)
		machinemanager.POST("/deletemachinedata", controller.Deletemachinedata)
		machinemanager.POST("/adddepart", controller.AddDepart)
		machinemanager.POST("/updatedepart", controller.UpdateDepart)
		batchmanager.POST("/updatebatch", controller.UpdateBatch)
		batchmanager.POST("/deletebatch", controller.DeleteBatch)
	}
	// TODO: 此处绑定前端静态资源handler
	router.Static("/static", "./dist/static")
	router.StaticFile("/", "./dist/index.html")

	// 关键点【解决页面刷新404的问题】
	router.NoRoute(func(c *gin.Context) {
		url := c.Request.RequestURI
		c.Redirect(http.StatusFound, url)
		router.StaticFile(url, "./dist/index.html")
	})

	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

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
