/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	aconfig "gitee.com/openeuler/PilotGo/cmd/agent/app/config"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/filemonitor"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/localstorage"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/register"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func main() {
	// 加载系统配置
	err := aconfig.Init()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	// 检查os环境
	_, err = common.InitOSName()
	if err != nil {
		logger.Error("os detect failed: %s", err)
		os.Exit(-1)
	}

	// 初始化日志
	if err := logger.Init(&aconfig.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		os.Exit(-1)
	}
	logger.Info("Start PilotGo agent.")

	// 定时任务初始化
	if err := common.CronInit(); err != nil {
		logger.Error("cron init failed: %s", err)
		os.Exit(-1)
	}

	// init agent info
	if err := localstorage.Init(); err != nil {
		logger.ErrorStack("local storage init failed", err)
		os.Exit(-1)
	}
	logger.Info("agent uuid is:%s", localstorage.AgentUUID())

	go func(conf *aconfig.Server) {
		// 与server握手
		client := network.NewSocketClient()
		register.RegitsterHandler(client)
		go filemonitor.FileMonitor(client)

		for {
			logger.Info("start to connect server")
			err = client.Run(&aconfig.Config().Server)
			if err != nil {
				logger.Error("socket client exit, error:%s", err.Error())
			}

			// 延迟5s+5s内随机时间
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			delayTime := time.Second*5 + time.Duration(r.Uint32()%5000*uint32(time.Millisecond))
			time.Sleep(delayTime)
		}
	}(&aconfig.Config().Server)

	// 文件监控初始化
	if err := filemonitor.FileMonitorInit(); err != nil {
		logger.Error("config file monitor init failed: %s", err)
		os.Exit(-1)
	}

	// 信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}
