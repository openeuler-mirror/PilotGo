package agentmanager

import (
	"fmt"
	"net"
	"sync"
)

// 用于管理server连接的agent
type AgentManager struct {
	agentMap sync.Map
}

var globalAgentManager = AgentManager{}

func GetAgentManager() *AgentManager {
	return &globalAgentManager
}

// func regitsterHandler(s *network.SocketServer) {
// 	s.BindHandler(protocol.Heartbeat, func(conn net.Conn, data []byte) error {
// 		fmt.Println("process heartbeat:", string(data))
// 		return nil
// 	})

// 	s.BindHandler(protocol.RunScript, func(conn net.Conn, data []byte) error {
// 		fmt.Println("process run script command:", string(data))
// 		return nil
// 	})
// }

func (am *AgentManager) Stop() {
	// stop server here
}

func (am *AgentManager) AddAgent(a *Agent) {
	am.agentMap.Store(a.UUID, a)
}

func (am *AgentManager) GetAgent(a *Agent) {
	am.agentMap.Store(a.UUID, a)
}

func (am *AgentManager) DeleteAgent(uuid string) {
	if _, ok := am.agentMap.LoadAndDelete(uuid); ok {
	}
}

func AddandRunAgent(c net.Conn) {
	agent, _ := NewAgent(c)
	agent.StartListen()
	GetAgentManager().AddAgent(agent)
	fmt.Println("add new agent from:", c.RemoteAddr().String())

}

func StopAgentManager() {

}
