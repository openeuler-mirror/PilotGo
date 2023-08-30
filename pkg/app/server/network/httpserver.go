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
 * LastEditTime: 2023-08-30 10:54:41
 * Description: Interface routing forwarding
 ******************************************************************************/
package network

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/controller/agentcontroller"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/controller/pluginapi"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/network/websocket"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/resource"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auth"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

func HttpServerInit(conf *sconfig.HttpServer) error {
	if err := SessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := setupRouter()
		if conf.UseHttps {
			if conf.CertFile == "" || conf.KeyFile == "" {
				logger.Error("https cert or key not configd")
				return
			}

			logger.Info("start http service on: https://%s", conf.Addr)
			if err := r.RunTLS(conf.Addr, conf.CertFile, conf.KeyFile); err != nil {
				logger.Error("start http server failed:%v", err)
			}
		} else {
			logger.Info("start http service on: http://%s", conf.Addr)
			if err := r.Run(conf.Addr); err != nil {
				logger.Error("start http server failed:%v", err)
			}
		}
	}()

	if conf.Debug {
		go func() {
			// 分解字符串然后添加后缀6060
			portIndex := strings.Index(conf.Addr, ":")
			addr := conf.Addr[:portIndex] + ":6060"
			logger.Debug("start pprof service on: %s", addr)
			if conf.UseHttps {
				if conf.CertFile == "" || conf.KeyFile == "" {
					logger.Error("https cert or key not configd")
					return
				}

				err := http.ListenAndServeTLS(addr, conf.CertFile, conf.KeyFile, nil)
				if err != nil {
					logger.Error("failed to start pprof, error:%v", err)
				}
			} else {
				err := http.ListenAndServe(addr, nil)
				if err != nil {
					logger.Error("failed to start pprof, error:%v", err)
				}
			}
		}()
	}

	return nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(auth.LoggerDebug())
	router.Use(auth.Recover)

	// 绑定 http api handler
	registerAPIs(router)

	// 绑定前端静态资源handler
	resource.StaticRouter(router)

	// 绑定插件接口反向代理handler
	registerPluginGateway(router)

	// 全局通用接口
	router.GET("/ws", controller.WS)
	router.GET("/event", controller.PushAlarmHandler)

	return router
}

func registerAPIs(router *gin.Engine) {
	api := router.Group("/api/v1")
	overview := api.Group("/overview") // 机器概览
	{
		overview.GET("/info", controller.ClusterInfoHandler)
		overview.GET("/depart_info", controller.DepartClusterInfoHandler)
	}

	macList := api.Group("/macList") // 机器管理
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

	macDetails := api.Group("/api") // 机器详情
	{
		macDetails.GET("/agent_overview", agentcontroller.AgentOverviewHandler)
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

	macBasicModify := api.Group("/agent") // 机器配置
	{
		macBasicModify.GET("/sysctl_change", agentcontroller.SysctlChangeHandler)
		macBasicModify.POST("/service_stop", agentcontroller.ServiceStopHandler)
		macBasicModify.POST("/service_start", agentcontroller.ServiceStartHandler)
		macBasicModify.POST("/service_restart", agentcontroller.ServiceRestartHandler)
		macBasicModify.POST("/rpm_install", agentcontroller.InstallRpmHandler)
		macBasicModify.POST("/rpm_remove", agentcontroller.RemoveRpmHandler)
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

	batchmanager := api.Group("batchmanager") // 批次
	{
		batchmanager.GET("/batchinfo", controller.BatchInfoHandler)
		batchmanager.GET("/batchmachineinfo", controller.BatchMachineInfoHandler)
	}

	user := api.Group("user") // 用户管理
	{
		user.POST("/login", controller.LoginHandler)
		user.POST("/updatepwd", controller.UpdatePasswordHandler)
		user.GET("/logout", controller.Logout)
		user.GET("/searchAll", controller.UserAll)
		user.POST("/userSearch", controller.UserSearchHandler)
		user.GET("/info", auth.AuthMiddleware(), controller.Info)
		user.POST("/permission", controller.GetLoginUserPermissionHandler)
		user.GET("/roles", controller.GetRolesHandler)
		user.GET("/role", controller.GetUserRoleHandler)
		user.POST("/addRole", controller.AddUserRoleHandler)
		user.POST("/delRole", controller.DeleteUserRoleHandler)
		user.POST("/updateRole", controller.UpdateUserRoleHandler)
		user.POST("/roleChange", controller.RolePermissionChangeHandler)
	}

	configmanager := api.Group("config") // 配置管理
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

	userLog := api.Group("log") // 日志管理
	{
		userLog.GET("/log_all", controller.LogAllHandler)
		userLog.GET("/logs", controller.AgentLogsHandler)
	}

	// 此处绑定casbin过滤规则
	policy := api.Group("casbin")
	{
		policy.GET("/get", controller.GetPolicy)
		policy.POST("/delete", controller.PolicyDelete)
		policy.POST("/add", controller.PolicyAdd)
	}

	Level := api.Group("")
	Level.Use(auth.CasbinHandler())
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

	plugin := api.Group("plugins") // 插件
	{
		plugin.GET("", controller.GetPluginsHandler)
		plugin.PUT("", controller.AddPluginHandler)
		plugin.POST("/:uuid", controller.TogglePluginHandler)
		plugin.DELETE("/:uuid", controller.UnloadPluginHandler)
	}

	// 对插件提供的api接口
	registerPluginApi(api)

	other := api.Group("")
	{
		// 监控机器列表
		other.GET("/macList/machinealldata", controller.MachineAllDataHandler)
		// 未用到
		other.GET("/macList/departinfo", controller.DepartInfoHandler)
		// 配置批次下发
		other.GET("/batchmanager/selectbatch", controller.SelectBatchHandler)
		// 健康监测
		other.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	}
	go websocket.CliManager.Start()
}

func registerPluginApi(router *gin.RouterGroup) {
	pluginAPI := router.Group("/pluginapi")
	pluginAPI.Use(pluginapi.AuthCheck)
	{
		pluginAPI.POST("/run_command", pluginapi.RunCommandHandler)
		pluginAPI.POST("/run_script", pluginapi.RunScriptHandler)

		pluginAPI.PUT("/listener", pluginapi.RegisterListenerHandler)
		pluginAPI.DELETE("/listener", pluginapi.UnregisterListenerHandler)

		pluginAPI.PUT("/install_package", pluginapi.InstallPackage)
		pluginAPI.PUT("/uninstall_package", pluginapi.UninstallPackage)

		pluginAPI.GET("/service/:name", pluginapi.Service)
		pluginAPI.PUT("/start_service", pluginapi.StartService)
		pluginAPI.PUT("/stop_service", pluginapi.StopService)

		pluginAPI.GET("/machine_list", pluginapi.MachineList)
	}
	// plugin
	{
		pluginAPI.GET("/plugins", pluginapi.PluginList)
	}
}

func registerPluginGateway(router *gin.Engine) {
	gateway := router.Group("/plugin/:plugin_name")
	gateway.Any("", controller.PluginGatewayHandler)
}
