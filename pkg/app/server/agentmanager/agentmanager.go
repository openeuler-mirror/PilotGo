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
 * Date: 2021-01-18 02:33:55
 * LastEditTime: 2022-04-11 16:27:35
 * Description: socket server
 ******************************************************************************/
package agentmanager

import (
	"net"
	"strings"
	"sync"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
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
	logger.Info("Add new agent from:%s", c.RemoteAddr().String())
	AddAgents2DB(agent)
}

func StopAgentManager() {

}

func AddAgents2DB(a *Agent) {
	agent_uuid := GetAgent(a.UUID)
	if agent_uuid == nil {
		logger.Error("获取uuid失败!")
		return
	}
	agent_OS, err := agent_uuid.GetAgentOSInfo()
	if err != nil {
		logger.Error("初始化系统信息失败!")
		return
	}
	agentOS := strings.Split(agent_OS.(string), ";")
	UUIDExistbool, err := dao.IsUUIDExist(a.UUID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if UUIDExistbool {
		logger.Warn("机器%s已经存在!", agentOS[0])
		departId, err := dao.UUIDForDepartId(a.UUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		if departId != global.UncateloguedDepartId {
			err := dao.MachineStatusToNormal(a.UUID, agentOS[0])
			if err != nil {
				logger.Error(err.Error())
				return
			}
		} else {
			err := dao.MachineStatusToFree(a.UUID, agentOS[0])
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}
		return
	}

	agent_list := model.MachineNode{
		IP:          agentOS[0],
		MachineUUID: a.UUID,
		DepartId:    global.UncateloguedDepartId,
		Systeminfo:  agentOS[1] + " " + agentOS[2],
		CPU:         agentOS[3],
		State:       global.Free,
	}
	err = dao.AddNewMachine(agent_list)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
