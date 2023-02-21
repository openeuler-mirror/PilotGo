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
 * LastEditTime: 2023-02-21 19:03:23
 * Description: agent main
 ******************************************************************************/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	aconfig "openeuler.org/PilotGo/PilotGo/pkg/app/agent/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/filemonitor"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/localstorage"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/register"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func main() {
	// 加载系统配置
	err := aconfig.Init()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
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
		logger.Error("local storage init failed: %s", err)
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
