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
 * LastEditTime: 2022-03-01 13:12:29
 * Description: server main
 ******************************************************************************/
package main

import (
	"fmt"
	"os"
	"strconv"

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network"
	"openeluer.org/PilotGo/PilotGo/pkg/cmd"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func main() {

	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	logger.Init(conf)
	logger.Info("Thanks to choose PilotGo!")

	server := &network.SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}

	// agentmanager := agentmanager.GetAgentManager()

	// 启动agent socket server
	url := conf.S.ServerIP + ":" + strconv.Itoa(conf.SocketPort)
	go func() {
		server.Run(url)
	}()

	// 此处启动前端及REST http server
	go func() {
		// 连接数据库及启动router
		err = cmd.Start(conf)
		if err != nil {
			logger.Info("server start failed:%s", err.Error())
			os.Exit(-1)
		}
		// network.HttpServerStart("192.168.160.128:8084")

	}()
	select {}
}
