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
 * LastEditTime: 2023-07-10 14:49:50
 * Description: socket server
 ******************************************************************************/
package agentmanager

import (
	"net"
	"sync"

	machineservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
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

func AddandRunAgent(c net.Conn) error {
	agent, err := NewAgent(c)
	if err != nil {
		logger.Warn("create agent from conn error, error:%s , remote addr is:%s",
			err.Error(), c.RemoteAddr().String())
		return err
	}

	AddAgent(agent)
	logger.Info("Add new agent from:%s", c.RemoteAddr().String())
	AddAgents2DB(agent)

	return nil
}

func StopAgentManager() {
	// TODO
}

func AddAgents2DB(a *Agent) {
	agent_uuid := GetAgent(a.UUID)
	if agent_uuid == nil {
		logger.Error("获取uuid失败!")
		return
	}
	// TODO: 沒有对message对象status、Error字段进行判断，决定后续步骤是否执行
	agent_os, err := agent_uuid.GetAgentOSInfo()
	if err != nil {
		logger.Error("初始化系统信息失败: %s", err.Error())
		return
	}
	//查找此uuid是否已存入数据库
	node, err := machineservice.MachineInfoByUUID(a.UUID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//如果uuid已存入，则修改ip和状态等信息
	if node.ID != 0 {
		logger.Warn("机器%s已经存在!", agent_os.IP)
		Ma := &machineservice.MachineNode{
			IP:        agent_os.IP,
			RunStatus: machineservice.OnlineStatus,
		}
		if node.Departid != global.UncateloguedDepartId {
			Ma.MaintStatus = machineservice.NormalStatus
		} else {
			Ma.MaintStatus = machineservice.MaintenanceStatus
		}
		err := machineservice.UpdateMachine(a.UUID, Ma)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	//新添加一台机器信息的时候自动分配到根节点部门，并设为在线状态和维护状态
	agent_list := &machineservice.MachineNode{
		IP:          agent_os.IP,
		MachineUUID: a.UUID,
		DepartId:    global.UncateloguedDepartId,
		Systeminfo:  agent_os.PrettyName,
		CPU:         agent_os.ModelName,
		RunStatus:   machineservice.OnlineStatus,
		MaintStatus: machineservice.MaintenanceStatus,
	}
	err = machineservice.AddMachine(agent_list)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
