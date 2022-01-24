package main

import (
	"fmt"
	"os"

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
	go func() {
		server.Run("192.168.160.128:8879")
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
