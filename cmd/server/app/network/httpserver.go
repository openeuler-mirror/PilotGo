/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package network

import (
	"context"
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/controller"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/controller/agentcontroller"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/controller/pluginapi"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/middleware"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/websocket"
	"gitee.com/openeuler/PilotGo/cmd/server/app/resource"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/klog/v2"
)

func HttpServerInit(conf *options.HttpServer, stopCh <-chan struct{}) error {
	if err := SessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := setupRouter()
		// 启动websocket服务
		go websocket.CliManager.Start(stopCh)
		shutdownCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		srv := &http.Server{
			Addr:    conf.Addr,
			Handler: r,
		}
		go func() {
			<-stopCh
			klog.Warningln("httpserver prepare stop")
			_ = srv.Shutdown(shutdownCtx)
		}()
		// 启动http server服务
		if conf.UseHttps {
			if conf.CertFile == "" || conf.KeyFile == "" {
				logger.Error("https cert or key not configd")
				return
			}

			logger.Info("start http service on: https://%s", conf.Addr)

			if err := srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil {
				if err != http.ErrServerClosed {
					logger.Error("ListenAndServeTLS start http server failed:%v", err)
					return
				}
			}
		} else {
			logger.Info("start http service on: http://%s", conf.Addr)
			if err := srv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					logger.Error("ListenAndServe start http server failed:%v", err)

				}

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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(logger.RequestLogger([]string{
		"/api/v1/pluginapi/heartbeat",
	}))
	router.Use(middleware.Recover)

	// 绑定 http api handler
	registerAPIs(router)

	// 绑定插件接口反向代理handler
	registerPluginGateway(router)

	// 全局通用接口
	router.GET("/ws", controller.WS)
	router.GET("/event", controller.PushAlarmHandler)

	// 绑定前端静态资源handler
	resource.StaticRouter(router)

	return router
}

func registerAPIs(router *gin.Engine) {
	noAuthenApis := router.Group("/api/v1")
	{
		noAuthenApis.GET("/version", controller.VersionHandler)

		noAuthenApis.POST("/user/login", controller.LoginHandler)
		noAuthenApis.GET("/user/logout", controller.Logout)
		noAuthenApis.POST("/user/permission", controller.GetLoginUserPermissionHandler)
		noAuthenApis.GET("/plugins", controller.GetPluginsHandler)
		noAuthenApis.GET("/plugins_paged", controller.GetPluginsPagedHandler)
		noAuthenApis.GET("/download/:filename", controller.Download)
	}

	authenApi := router.Group("/api/v1")
	{
		{
			macBasicModify := authenApi.Group("/agent")
			macBasicModify.POST("/rpm_install", middleware.NeedPermission("rpm_install", "button"), agentcontroller.InstallRpmHandler)
			macBasicModify.POST("/rpm_remove", middleware.NeedPermission("rpm_uninstall", "button"), agentcontroller.RemoveRpmHandler)
		}
		{
			batchmanager := authenApi.Group("batchmanager")
			batchmanager.POST("/updatebatch", middleware.NeedPermission("batch_update", "button"), controller.UpdateBatchHandler)
			batchmanager.POST("/deletebatch", middleware.NeedPermission("batch_delete", "button"), controller.DeleteBatchHandler)
			batchmanager.POST("/createbatch", middleware.NeedPermission("batch_create", "button"), controller.CreateBatchHandler)
		}
		{
			user := authenApi.Group("user")
			user.POST("/register", middleware.NeedPermission("user_add", "button"), controller.RegisterHandler)
			user.POST("/reset", middleware.NeedPermission("user_reset", "button"), controller.ResetPasswordHandler)
			user.POST("/delete", middleware.NeedPermission("user_del", "button"), controller.DeleteUserHandler)
			user.POST("/update", middleware.NeedPermission("user_edit", "button"), controller.UpdateUserHandler)
			user.POST("/import", middleware.NeedPermission("user_import", "button"), controller.ImportUser)

			user.POST("/addRole", middleware.NeedPermission("role_add", "button"), controller.AddRoleHandler)
			user.POST("/delRole", middleware.NeedPermission("role_delete", "button"), controller.DeleteRoleHandler)
			user.POST("/updateRole", middleware.NeedPermission("role_update", "button"), controller.UpdateRoleInfoHandler)
			user.POST("/roleChange", middleware.NeedPermission("role_modify", "button"), controller.RolePermissionChangeHandler)
		}
		{
			macList := authenApi.Group("/macList")
			macList.POST("/deletedepartdata", middleware.NeedPermission("dept_delete", "button"), controller.DeleteDepartDataHandler)
			macList.POST("/adddepart", middleware.NeedPermission("dept_add", "button"), controller.AddDepartHandler)
			macList.POST("/updatedepart", middleware.NeedPermission("dept_update", "button"), controller.UpdateDepartHandler)
		}
		/*
			{
				configmanager := authenApi.Group("config")
				configmanager.POST("/file_broadcast", middleware.NeedPermission("config_install", "button"), agentcontroller.ConfigFileBroadcastToAgents)
			}
		*/
	}

	api := router.Group("/api/v1")
	api.Use(middleware.TokenCheckMiddleware)
	overview := api.Group("/overview") // 机器概览
	{
		overview.GET("/info", controller.ClusterInfoHandler)
		overview.GET("/depart_info", controller.DepartClusterInfoHandler)
	}

	macList := api.Group("/macList") // 机器管理
	{
		macList.POST("/script_save", controller.AddScriptHandler)
		macList.POST("/deletemachine", middleware.NeedPermission("machine_delete", "button"), controller.DeleteMachineHandler)
		macList.GET("/depart", controller.DepartHandler)
		macList.GET("/selectmachine", controller.MachineListHandler)
		macList.GET("/machineinfo", controller.MachineInfoHandler)
		macList.POST("/modifydepart", middleware.NeedPermission("dept_change", "button"), controller.ModifyMachineDepartHandler)
		macList.GET("/sourcepool", controller.FreeMachineSource)
		macList.POST("/gettags", pluginapi.GetTagHandler)
	}

	macDetails := api.Group("/api") // 机器详情
	{
		macDetails.GET("/agent_overview", agentcontroller.AgentOverviewHandler)
		macDetails.GET("/agent_list", agentcontroller.AgentListHandler)
		macDetails.GET("/run_command", agentcontroller.RunCmd)
		macDetails.GET("/run_script", agentcontroller.RunScriptWithBooleanCheck)
		macDetails.GET("/os_info", agentcontroller.OSInfoHandler)
		macDetails.GET("/cpu_info", agentcontroller.CPUInfoHandler)
		macDetails.GET("/memory_info", agentcontroller.MemoryInfoHandler)
		macDetails.GET("/sysctl_info", agentcontroller.SysInfoHandler)
		macDetails.GET("/sysctl_view", agentcontroller.SysctlViewHandler)
		macDetails.GET("/service_list", agentcontroller.ServiceListHandler)
		macDetails.GET("/service_status", agentcontroller.ServiceStatusHandler)
		macDetails.GET("/rpm_all", agentcontroller.AllRpmHandler)
		macDetails.GET("/rpm_source", agentcontroller.RpmSourceHandler)
		macDetails.GET("/repos", agentcontroller.GetAgentRepo)
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
		macDetails.GET("/net", agentcontroller.GetAgentNetworkConnect)
	}
	/*
		macBasicModify := api.Group("/agent") // 机器配置
		{
			macBasicModify.GET("/sysctl_change", agentcontroller.SysctlChangeHandler)
			macBasicModify.POST("/service_stop", agentcontroller.ServiceStopHandler)
			macBasicModify.POST("/service_start", agentcontroller.ServiceStartHandler)
			macBasicModify.POST("/service_restart", agentcontroller.ServiceRestartHandler)
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
	*/
	batchmanager := api.Group("batchmanager") // 批次
	{
		batchmanager.GET("/batchinfo", controller.BatchInfoHandler)
		batchmanager.GET("/batchmachineinfo", controller.BatchMachineInfoHandler)
	}

	user := api.Group("user") // 用户管理
	{
		user.POST("/updatepwd", controller.UpdatePasswordHandler)
		// user.GET("/logout", controller.Logout)
		user.GET("/searchAll", controller.UserAll)
		user.POST("/userSearch", controller.UserSearchHandler)
		user.GET("/info", controller.Info)
		// user.POST("/permission", controller.GetLoginUserPermissionHandler)
		user.GET("/roles", controller.GetRolesHandler)
		user.GET("/roles_paged", controller.GetRolesPagedHandler)
	}

	/*
		configmanager := api.Group("config") // 配置管理
		{
			configmanager.GET("/read_file", agentcontroller.ReadFile)
			configmanager.POST("/fileSaveAdd", controller.SaveConfigFileToDatabaseHandler)
			configmanager.GET("/file_all", controller.AllConfigFiles)
			configmanager.POST("/file_search", controller.ConfigFileSearchHandler)
			configmanager.POST("/file_update", controller.UpdateConfigFileHandler)
			configmanager.POST("/file_delete", controller.DeleteConfigFileHandler)
			configmanager.GET("/lastfile_all", controller.HistoryConfigFilesHandler)
			configmanager.POST("/lastfile_rollback", controller.LastConfigFileRollBackHandler)
		}
	*/

	userLog := api.Group("log") // 日志管理
	{
		userLog.GET("/log_all", middleware.NeedPermission("audit", "menu"), controller.LogAllHandler)
		userLog.GET("/log_child", controller.GetAuditLogByIdHandler)
	}

	plugin := api.Group("plugins") // 插件
	{
		// 添加插件
		plugin.PUT("", controller.AddPluginHandler)
		// 启用/停用plugin
		plugin.POST("/:uuid", controller.TogglePluginHandler)
		// 删除插件
		plugin.DELETE("/:uuid", controller.UnloadPluginHandler)
	}

	// 对插件提供的api接口
	registerPluginApi(router)

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
}

func registerPluginApi(router *gin.Engine) {
	pluginAPI := router.Group("/api/v1/pluginapi")
	pluginAPI.Use(pluginapi.AuthCheck)
	{
		pluginAPI.POST("/upload", controller.Upload)

		pluginAPI.POST("/run_command_async", pluginapi.RunCommandAsyncHandler)
		pluginAPI.POST("/run_command", pluginapi.RunCommandHandler)
		pluginAPI.POST("/run_script", pluginapi.RunScriptHandler)

		pluginAPI.PUT("/install_package", pluginapi.InstallPackage)
		pluginAPI.PUT("/uninstall_package", pluginapi.UninstallPackage)

		pluginAPI.GET("/service/:name", pluginapi.Service)
		pluginAPI.PUT("/start_service", pluginapi.StartService)
		pluginAPI.PUT("/stop_service", pluginapi.StopService)

		pluginAPI.GET("/machine_list", pluginapi.MachineList)
		pluginAPI.GET("/machine_info", pluginapi.MachineInfoByUUID)
		pluginAPI.POST("/file_deploy", pluginapi.FileDeploy)

		pluginAPI.GET("/batch_list", pluginapi.BatchListHandler)
		pluginAPI.GET("/batch_uuid", pluginapi.MachineListOfBatch)
		pluginAPI.POST("/getnodefiles", pluginapi.GetNodeFiles)
	}
	// plugin
	{
		pluginAPI.GET("/plugins", pluginapi.PluginList)
		pluginAPI.POST("/heartbeat", pluginapi.PluginHeartbeat)
		pluginAPI.POST("/has_permission", pluginapi.HasPermission)
	}
}

func registerPluginGateway(router *gin.Engine) {
	gateway := router.Group("/plugin")
	logger.Info("gateway process plugin request")
	gateway.Any("/:plugin_name", controller.PluginGatewayHandler)
	gateway.Any("/:plugin_name/*action", controller.PluginGatewayHandler)

	gateway.GET("/ws/:plugin_name", controller.PluginWebsocketGatewayHandler)
}
