package agentmanager

import (
	"fmt"
	"net"
	"sync"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

// 用于管理server连接的agent
type AgentManager struct {
	agentMap sync.Map
}

var globalAgentManager = AgentManager{}

// func GetAgentManager() *AgentManager {
// 	return &globalAgentManager
// }

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

func AddAgent(a *Agent) {
	globalAgentManager.agentMap.Store(a.UUID, a)
}

func GetAgent(uuid string) *Agent {
	agent, ok := globalAgentManager.agentMap.Load(uuid)
	if ok {
		return agent.(*Agent)
	}
	return nil
}

func GetAgentList() []map[string]string {

	agentList := []map[string]string{}

	globalAgentManager.agentMap.Range(
		func(uuid interface{}, agent interface{}) bool {
			agentInfo := map[string]string{}
			agentInfo["agent_version"] = agent.(*Agent).Version
			agentInfo["agent_uuid"] = agent.(*Agent).UUID

			agentList = append(agentList, agentInfo)
			return true
		},
	)

	return agentList
}

func DeleteAgent(uuid string) {
	if _, ok := globalAgentManager.agentMap.LoadAndDelete(uuid); !ok {
		logger.Warn("delete known agent:%s", uuid)
	}
}

func AddandRunAgent(c net.Conn) {
	agent, err := NewAgent(c)
	if err != nil {
		logger.Warn("create agent from conn error, error:%s , remote addr is:%s",
			err.Error(), c.RemoteAddr().String())
	}

	AddAgent(agent)
	fmt.Println("add new agent from:", c.RemoteAddr().String())
}

func StopAgentManager() {

}
