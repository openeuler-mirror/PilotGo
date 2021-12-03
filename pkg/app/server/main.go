package main

import (
  "openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
  "openeluer.org/PilotGo/PilotGo/pkg/app/server/network"
)

func main() {

	server := &network.SocketServer{
		// MessageProcesser: protocol.NewMessageProcesser(),
		OnAccept: agentmanager.AddandRunAgent,
		OnStop:   agentmanager.StopAgentManager,
	}

	// agentmanager := agentmanager.GetAgentManager()

	go func() {
		server.Run("localhost:8879")
	}()

	select {}

}

