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
 * Date: 2021-11-18 10:25:52
 * LastEditTime: 2022-04-18 15:16:10
 * Description: server main
 ******************************************************************************/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/initialization"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/webSocket"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

func main() {
	err := sconfig.Init()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	if err := logger.Init(&sconfig.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo!")

	// redis db初始化
	if err := initialization.RedisdbInit(&sconfig.Config().RedisDBinfo); err != nil {
		logger.Error("redis db init failed, please check again: %s", err)
		os.Exit(-1)
	}

	// mysql db初始化
	if err := initialization.MysqldbInit(&sconfig.Config().MysqlDBinfo); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		os.Exit(-1)
	}

	// 鉴权模块初始化
	global.PILOTGO_E = service.Casbin(&sconfig.Config().MysqlDBinfo)

	// 监控初始化
	if err := initialization.MonitorInit(&sconfig.Config().Monitor); err != nil {
		logger.Error("monitor init failed: %s", err)
		os.Exit(-1)
	}

	// 启动agent socket server
	if err := initialization.SocketServerInit(&sconfig.Config().SocketServer); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	//此处启动前端及REST http server
	serverAddr, err := initialization.HttpServerInit(&sconfig.Config().HttpServer)
	controller.ServerAddr = serverAddr
	if err != nil {
		logger.Error("socket server init failed, error:%v", err)
		os.Exit(-1)
	}

	// 初始化插件组件
	if err = service.PluginServiceInit(); err != nil {
		logger.Error("plugin service init failed, error:%v", err)
		os.Exit(-1)
	}

	logger.Info("start to serve.")

	// 前端推送告警
	agentmanager.WARN_MSG = make(chan interface{}, 1)
	go webSocket.SendWarnMsgToWeb()

	// 信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			global.PILOTGO_REDIS.Close()

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}
