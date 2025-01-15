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
		r := SetupRouter()
		// start websocket server
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
		// start http server
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
			// pprof
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

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(logger.RequestLogger([]string{
		"/api/v1/pluginapi/heartbeat",
		"/",
	}))
	router.Use(middleware.Recover)

	// 绑定 http api handler
	registerAPIs(router)

	// 对插件提供的api接口
	registerPluginApi(router)

	// 绑定插件接口反向代理handler
	registerPluginGateway(router)

	// 绑定前端静态资源handler
	resource.StaticRouter(router)

	return router
}

func registerAPIs(router *gin.Engine) {
	router.GET("/event", middleware.TokenCheckMiddleware, controller.PushAlarmHandler)

	api := router.Group("/api/v1")

	noTokenApi := api.Group("") // 首页登录、退出
	{
		noTokenApi.GET("/version", controller.VersionHandler)

		noTokenApi.POST("/user/login", controller.LoginHandler)
		noTokenApi.GET("/user/logout", controller.Logout)

		noTokenApi.GET("/download/:filename", controller.Download)

		noTokenApi.GET("/webterminal", controller.WebTerminal)
	}

	authenApi := api.Group("") // 按钮权限，是否显示
	{
		{
			macBasicModify := authenApi.Group("/agent")
			macBasicModify.POST("/rpm_install", middleware.NeedPermission("rpm_install", "button"), agentcontroller.InstallRpmHandler)
			macBasicModify.POST("/rpm_remove", middleware.NeedPermission("rpm_uninstall", "button"), agentcontroller.RemoveRpmHandler)
		}
		{
			batch := authenApi.Group("/batchmanager")
			batch.POST("/updatebatch", middleware.NeedPermission("batch_update", "button"), controller.UpdateBatchHandler)
			batch.POST("/deletebatch", middleware.NeedPermission("batch_delete", "button"), controller.DeleteBatchHandler)
			batch.POST("/createbatch", middleware.NeedPermission("batch_create", "button"), controller.CreateBatchHandler)
		}
		{
			user := authenApi.Group("/user")
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
			system := authenApi.Group("/macList")
			system.POST("/deletedepartdata", middleware.NeedPermission("dept_delete", "button"), controller.DeleteDepartDataHandler)
			system.POST("/adddepart", middleware.NeedPermission("dept_add", "button"), controller.AddDepartHandler)
			system.POST("/updatedepart", middleware.NeedPermission("dept_update", "button"), controller.UpdateDepartHandler)
			system.POST("/modifydepart", middleware.NeedPermission("dept_change", "button"), controller.ModifyMachineDepartHandler)
			system.POST("/deletemachine", middleware.NeedPermission("machine_delete", "button"), controller.DeleteMachineHandler)
		}
		{
			script := authenApi.Group("/script_auth")
			script.POST("/run", middleware.NeedPermission("run_script", "button"), controller.RunScriptHandler)
			script.PUT("/update_blacklist", middleware.NeedPermission("update_script_blacklist", "button"), controller.UpdateCommandsBlackListHandler)
		}
	}

	tokenApi := api.Group("") // web页面显示
	tokenApi.Use(middleware.TokenCheckMiddleware)

	overview := tokenApi.Group("/overview") // 概览
	{
		overview.GET("/info", controller.ClusterInfoHandler)
		overview.GET("/depart_info", controller.DepartClusterInfoHandler)
	}

	system := tokenApi.Group("") // 系统
	{
		{
			//machine list
			mac := system.Group("/macList")
			mac.GET("/machineinfo", controller.MachineInfoHandler)
			mac.POST("/gettags", pluginapi.GetTagHandler)
			mac.GET("/machinealldata", controller.MachineAllDataHandler)
		}
		{
			// depart manager
			depart := system.Group("/macList")
			depart.GET("/depart", controller.DepartHandler)
			depart.GET("/departinfo", controller.DepartInfoHandler)
		}
		{
			// batch related
			batch := system.Group("/macList")
			batch.GET("/selectmachine", controller.MachineListHandler)

			batchmanager := system.Group("/batchmanager")
			{
				batchmanager.GET("/batchinfo", controller.BatchInfoHandler)
				batchmanager.GET("/batchmachineinfo", controller.BatchMachineInfoHandler)
				batchmanager.GET("/selectbatch", controller.SelectBatchHandler)
			}
		}
		{
			// machine detail info
			macDetails := system.Group("/api")
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
		{
			// script manager
			script := system.Group("/script")
			script.POST("/create", controller.AddScriptHandler)
			script.PUT("/update", controller.UpdateScriptHandler)
			script.DELETE("/delete", controller.DeleteScriptHandler)
			script.GET("/list_all", controller.GetScriptListHandler)
			script.GET("/list_history", controller.GetScriptHistoryVersionHandler)
			script.GET("/blacklist", controller.GetDangerousCommandsList)
		}
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

	user := tokenApi.Group("/user") // 用户、角色管理
	{
		{
			u := user.Group("")
			u.GET("/info", controller.Info)
			u.POST("/updatepwd", controller.UpdatePasswordHandler)
			u.GET("/searchAll", controller.UserAll)
			u.POST("/userSearch", controller.UserSearchHandler)
			u.POST("/permission", controller.GetLoginUserPermissionHandler)
		}
		{
			role := user.Group("")
			role.GET("/roles", controller.GetRolesHandler)
			role.GET("/roles_paged", controller.GetRolesPagedHandler)
		}
	}

	AuditLog := tokenApi.Group("/log") // 审计日志
	{
		AuditLog.GET("/log_all", controller.LogAllHandler)
		// TODO: 界面未调用该接口
		AuditLog.GET("/log_child", controller.GetAuditLogByIdHandler)
	}

	plugin := tokenApi.Group("") // 插件管理
	{
		plugin.GET("/plugins_paged", controller.GetPluginsPagedHandler)
		p := plugin.Group("/plugins")
		{
			// 添加插件
			p.PUT("", controller.AddPluginHandler)
			// 启用/停用plugin
			p.POST("/:uuid", controller.TogglePluginHandler)
			// 删除插件
			p.DELETE("/:uuid", controller.UnloadPluginHandler)
			// 获取插件列表
			p.GET("/", controller.GetPluginsHandler)
		}
	}
}

func registerPluginApi(router *gin.Engine) {
	api := router.Group("/api/v1")
	pluginAPI := api.Group("/pluginapi")
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
