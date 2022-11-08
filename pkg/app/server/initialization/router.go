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
package initialization

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager/agentcontroller"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/middleware"
	"openeuler.org/PilotGo/PilotGo/resource"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.LoggerDebug())
	router.Use(middleware.Recover)

	// TODO: 此处绑定 http api handler
	api := router.Group("/api/v1")

	overview := api.Group("/overview")
	{
		overview.GET("/info", controller.ClusterInfoHandler)
		overview.GET("/depart_info", controller.DepartClusterInfoHandler)
	}

	macList := api.Group("cluster/macList")
	{
		macList.POST("/script_save", controller.AddScriptHandler)
		macList.POST("/deletemachine", controller.DeleteMachineHandler)
		macList.GET("/depart", controller.DepartHandler)
		macList.GET("/selectmachine", controller.MachineListHandler)
		macList.POST("/createbatch", controller.CreateBatchHandler)
		macList.GET("/machineinfo", controller.MachineInfoHandler)
		macList.POST("/modifydepart", controller.ModifyMachineDepartHandler)
		macList.GET("/sourcepool", controller.FreeMachineSource)
	}

	macDetails := api.Group("cluster/macList/api")
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
		macDetails.GET("/firewall_zone", agentcontroller.FirewalldZoneConfig)
		macDetails.GET("/repos", agentcontroller.GetAgentRepo)
		macDetails.GET("/net", agentcontroller.GetAgentNetworkConnect)
	}

	macBasicModify := api.Group("cluster/macList/agent")
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
		macBasicModify.GET("/firewall_restart", agentcontroller.FirewalldRestart)
		macBasicModify.GET("/firewall_stop", agentcontroller.FirewalldStop)
		macBasicModify.POST("/firewall_addzp", agentcontroller.FirewalldZonePortAdd)
		macBasicModify.POST("/firewall_delzp", agentcontroller.FirewalldZonePortDel)
		macBasicModify.POST("/firewall_default", agentcontroller.FirewalldSetDefaultZone)
		macBasicModify.POST("/firewall_serviceAdd", agentcontroller.FirewalldServiceAdd)
		macBasicModify.POST("/firewall_serviceRemove", agentcontroller.FirewalldServiceRemove)
		macBasicModify.POST("/firewall_sourceAdd", agentcontroller.FirewalldSourceAdd)
		macBasicModify.POST("/firewall_sourceRemove", agentcontroller.FirewalldSourceRemove)
		macBasicModify.POST("/network", agentcontroller.ConfigNetworkConnect)
	}

	monitor := api.Group("prometheus")
	{
		monitor.GET("/queryrange", controller.QueryRange)
		monitor.GET("/query", controller.Query)
		monitor.GET("/alert", controller.ListenALert)
		monitor.POST("/alertmanager", controller.AlertMessageConfigHandler)
	}

	batchmanager := api.Group("batchmanager")
	{
		batchmanager.GET("/batchinfo", controller.BatchInfoHandler)
		batchmanager.GET("/batchmachineinfo", controller.BatchMachineInfoHandler)
	}

	user := api.Group("user")
	{
		user.POST("/login", controller.LoginHandler)
		user.GET("/logout", controller.Logout)
		user.GET("/searchAll", controller.UserAll)
		user.POST("/userSearch", controller.UserSearchHandler)
		user.GET("/info", middleware.AuthMiddleware(), controller.Info)
		user.POST("/permission", controller.GetLoginUserPermissionHandler)
		user.GET("/roles", controller.GetRolesHandler)
		user.GET("/role", controller.GetUserRoleHandler)
		user.POST("/addRole", controller.AddUserRoleHandler)
		user.POST("/delRole", controller.DeleteUserRoleHandler)
		user.POST("/updateRole", controller.UpdateUserRoleHandler)
		user.POST("/roleChange", controller.RolePermissionChangeHandler)
	}

	configmanager := api.Group("config")
	{
		configmanager.GET("/read_file", agentcontroller.ReadFile)
		configmanager.POST("/fileSaveAdd", controller.SaveFileToDatabaseHandler)
		configmanager.GET("/file_all", controller.AllFiles)
		configmanager.POST("/file_search", controller.FileSearchHandler)
		configmanager.POST("/file_update", controller.UpdateFileHandler)
		configmanager.POST("/file_delete", controller.DeleteFileHandler)
		configmanager.GET("/lastfile_all", controller.HistoryFilesHandler)
		configmanager.POST("/lastfile_rollback", controller.LastFileRollBackHandler)
		configmanager.POST("/file_broadcast", agentcontroller.FileBroadcastToAgents)
	}

	userLog := api.Group("log")
	{
		userLog.GET("/log_all", controller.LogAllHandler)
		userLog.GET("/logs", controller.AgentLogsHandler)
		userLog.POST("/delete", controller.DeleteLogHandler)
	}

	// 此处绑定casbin过滤规则
	policy := api.Group("casbin")
	{
		policy.GET("/get", controller.GetPolicy)
		policy.POST("/delete", controller.PolicyDelete)
		policy.POST("/add", controller.PolicyAdd)
	}

	Level := api.Group("")
	Level.Use(middleware.CasbinHandler())
	{
		user.POST("/register", controller.RegisterHandler)
		user.POST("/reset", controller.ResetPasswordHandler)
		user.POST("/delete", controller.DeleteUserHandler)
		user.POST("/update", controller.UpdateUserHandler)
		user.POST("/import", controller.ImportUser)
		macList.POST("/deletedepartdata", controller.DeleteDepartDataHandler)
		macList.POST("/adddepart", controller.AddDepartHandler)
		macList.POST("/updatedepart", controller.UpdateDepartHandler)
		batchmanager.POST("/updatebatch", controller.UpdateBatchHandler)
		batchmanager.POST("/deletebatch", controller.DeleteBatchHandler)
	}

	plugin := api.Group("plugins")
	{
		plugin.GET("", controller.GetPluginsHanlder)
		plugin.PUT("", controller.AddPluginHanlder)
		plugin.POST("/:uuid", controller.TogglePluginHanlder)
		plugin.DELETE("/:uuid", controller.UnloadPluginHanlder)
	}

	// TODO: 此处绑定前端静态资源handler
	resource.StaticRouter(router)

	// 全局通用接口
	router.GET("/ws", controller.ShellWs)
	router.GET("/macList/machinealldata", controller.MachineAllDataHandler)
	router.GET("/macList/departinfo", controller.DepartInfoHandler)
	router.GET("/macList/depart", controller.DepartHandler)
	// TODO: 不知道用途
	router.GET("/batchmanager/selectbatch", controller.SelectBatchHandler)
	router.GET("/event", controller.PushAlarmHandler)

	return router
}
