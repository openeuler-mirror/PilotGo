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

	overview := router.Group("/overview")
	{
		overview.GET("/info", controller.ClusterInfo)
		overview.GET("/depart_info", controller.DepartClusterInfo)
	}

	macList := router.Group("cluster/macList")
	{
		macList.GET("/selectmachine", controller.MachineList)
		macList.POST("/createbatch", controller.CreateBatch)
		macList.GET("/departinfo", controller.DepartInfo)
		macList.GET("/machineinfo", controller.MachineInfo)
		macList.GET("/depart", controller.Dept)
		macList.GET("/machinealldata", controller.MachineAllData)
		macList.POST("/modifydepart", controller.ModifyMachineDepart)
		macList.GET("/sourcepool", controller.FreeMachineSource)
	}
	macDetails := router.Group("cluster/macList/api")
	{
		macDetails.GET("/agent_info", agentcontroller.AgentInfoHandler)
		macDetails.GET("/agent_list", agentcontroller.AgentListHandler)
		macDetails.GET("/run_script", agentcontroller.RunScript)
		macDetails.GET("/os_info", agentcontroller.OSInfoHandler)
		macDetails.GET("/cpu_info", agentcontroller.CPUInfoHandler)
		macDetails.GET("/memory_info", agentcontroller.MemoryInfoHandler)
		macDetails.GET("/sysctl_info", agentcontroller.SysInfoHandler)
		macDetails.GET("/sysctl_view", agentcontroller.SysctlViewHandler)
		macDetails.GET("/service_list", agentcontroller.ServiceListHandler)
		macDetails.GET("/service_status", agentcontroller.ServiceStatusHandler)
		macDetails.GET("/rpm_all", agentcontroller.AllRpmHandler)
		macDetails.GET("/rpm_source", agentcontroller.RpmSourceHandler)
		macDetails.GET("/rpm_info", agentcontroller.RpmInfoHandler)
		macDetails.GET("/disk_use", agentcontroller.DiskUsageHandler)
		macDetails.GET("/disk_info", agentcontroller.DiskInfoHandler)
		macDetails.GET("/net_tcp", agentcontroller.NetTCPHandler)
		macDetails.GET("/net_udp", agentcontroller.NetUDPHandler)
		macDetails.GET("/net_io", agentcontroller.NetIOCounterHandler)
		macDetails.GET("/net_nic", agentcontroller.NetNICConfigHandler)
		macDetails.GET("/user_info", agentcontroller.CurrentUserInfoHandler)
		macDetails.GET("/user_all", agentcontroller.AllUserInfoHandler)
		macDetails.GET("/os_basic", agentcontroller.OsBasic)
		macDetails.GET("/firewall_config", agentcontroller.FirewalldConfig)
		macDetails.GET("/firewall_restart", agentcontroller.FirewalldRestart)
		macDetails.GET("/firewall_stop", agentcontroller.FirewalldStop)
		macDetails.POST("/firewall_addzp", agentcontroller.FirewalldZonePortAdd)
		macDetails.POST("/firewall_delzp", agentcontroller.FirewalldZonePortDel)
	}

	macBasicModify := router.Group("cluster/macList/agent")
	{
		macBasicModify.GET("/sysctl_change", agentcontroller.SysctlChangeHandler)
		macBasicModify.POST("/service_stop", agentcontroller.ServiceStopHandler)
		macBasicModify.POST("/service_start", agentcontroller.ServiceStartHandler)
		macBasicModify.POST("/service_restart", agentcontroller.ServiceRestartHandler)
		macBasicModify.POST("/rpm_install", agentcontroller.InstallRpmHandler)
		macBasicModify.POST("/rpm_remove", agentcontroller.RemoveRpmHandler)
		macBasicModify.GET("/disk_path", agentcontroller.DiskCreatPathHandler)
		macBasicModify.GET("/disk_mount", agentcontroller.DiskMountHandler)
		macBasicModify.GET("/disk_umount", agentcontroller.DiskUMountHandler)
		macBasicModify.GET("/disk_format", agentcontroller.DiskFormatHandler)
		macBasicModify.GET("/user_add", agentcontroller.AddLinuxUserHandler)
		macBasicModify.GET("/user_del", agentcontroller.DelUserHandler)
		macBasicModify.GET("/user_ower", agentcontroller.ChangeFileOwnerHandler)
		macBasicModify.GET("/user_per", agentcontroller.ChangePermissionHandler)
		macBasicModify.POST("cron_new", agentcontroller.CreatCron)
		macBasicModify.POST("/cron_del", agentcontroller.DeleteCronTask)
		macBasicModify.POST("/cron_update", agentcontroller.UpdateCron)
		macBasicModify.POST("/cron_status", agentcontroller.CronTaskStatus)
		macBasicModify.GET("/cron_list", agentcontroller.CronTaskList)
	}

	monitor := router.Group("prometheus")
	{
		monitor.POST("/queryrange", controller.Queryrange)
		monitor.POST("/query", controller.Query)
		monitor.GET("/alert", controller.ListenALert)
		monitor.POST("/alertmanager", controller.AlertMessageConfig)
	}

	batchmanager := router.Group("batchmanager")
	{
		batchmanager.GET("/batchinfo", controller.BatchInfo)
		batchmanager.GET("/batchmachineinfo", controller.BatchMachineInfo)
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
		user.POST("/addRole", controller.AddUserRole)
		user.POST("/delRole", controller.DeleteUserRole)
		user.POST("/updateRole", controller.UpdateUserRole)
		user.POST("/roleChange", controller.RolePermissionChange)
	}
	userLog := router.Group("log")
	{
		userLog.GET("/log_all", controller.LogAll)
		userLog.GET("/logs", controller.AgentLogs)
		userLog.POST("/delete", controller.DeleteLog)
	}

	// 此处绑定casbin过滤规则
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
		macList.POST("/deletedepartdata", controller.Deletedepartdata)
		macList.POST("/adddepart", controller.AddDepart)
		macList.POST("/updatedepart", controller.UpdateDepart)
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
	router.GET("/ws", controller.ShellWs) // 终端
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	return router
}
