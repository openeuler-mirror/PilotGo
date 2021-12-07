package main

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/network"
)

func main() {
	// conf, err := config.Load()
	// if err != nil {
	// 	fmt.Println("failed to load configure, exit..", err)
	// 	os.Exit(-1)
	// }

	// logger.Init(conf)
	// logger.Info("Thanks to choose PilotGo!")

	server := &network.SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}

	// agentmanager := agentmanager.GetAgentManager()

	// 启动agent socket server
	go func() {
		server.Run("localhost:8879")
	}()

	// 此处启动前端及REST http server
	go func() {
		network.HttpServerStart("localhost:8080")
	}()

	select {}
}
