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
 * Date: 2022-02-18 02:33:55
 * LastEditTime: 2023-07-11 16:52:58
 * Description: socket server's agentmanager
 ******************************************************************************/
package agentmanager

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"gitee.com/openeuler/PilotGo/app/agent/global"
	configservice "gitee.com/openeuler/PilotGo/app/server/service/configfile"
	machineservice "gitee.com/openeuler/PilotGo/app/server/service/machine"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	pnet "gitee.com/openeuler/PilotGo/utils/message/net"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	"github.com/google/uuid"
)

type AgentMessageHandler func(*Agent, *protocol.Message) error

var WARN_MSG chan interface{}

type Agent struct {
	UUID             string
	Version          string
	IP               string
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
	messageChan      chan *protocol.Message
}

// 通过给定的conn连接初始化一个agent并启动监听
func NewAgent(conn net.Conn) (*Agent, error) {
	agent := &Agent{
		UUID:             "agent",
		conn:             conn,
		MessageProcesser: protocol.NewMessageProcesser(),
		messageChan:      make(chan *protocol.Message, 50),
	}

	go func(agent *Agent) {
		for {
			msg := <-agent.messageChan
			logger.Debug("send message:%s", msg.String())
			pnet.SendBytes(agent.conn, protocol.TlvEncode(msg.Encode()))
		}
	}(agent)

	go func(agent *Agent) {
		agent.startListen()
	}(agent)

	if err := agent.Init(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *Agent) bindHandler(t int, f AgentMessageHandler) {
	a.MessageProcesser.BindHandler(t, func(c protocol.MessageContext, msg *protocol.Message) error {
		return f(c.(*Agent), msg)
	})
}

func (a *Agent) startListen() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("server processor panic error:%s", err.(error).Error())
			a.conn.Close()
		}
	}()

	readBuff := []byte{}
	for {
		buff := make([]byte, 1024)
		n, err := a.conn.Read(buff)
		if err != nil {
			err := machineservice.UpdateMachine(a.UUID, &machineservice.MachineNode{RunStatus: machineservice.OfflineStatus})
			if err != nil {
				logger.Error("update machine status failed: %s", err.Error())
			}
			DeleteAgent(a.UUID)
			str := "agent机器" + a.IP + "已断开连接"
			logger.Warn("agent %s disconnected, ip:%s", a.UUID, a.IP)
			WARN_MSG <- str
			return
		}
		readBuff = append(readBuff, buff[:n]...)

		// 切割frame
		i, f := protocol.TlvDecode(&readBuff)
		if i != 0 {
			readBuff = readBuff[i:]
			go func(a *Agent, f *[]byte) {
				msg := protocol.ParseMessage(*f)
				a.MessageProcesser.ProcessMessage(a, msg)
			}(a, f)
		}
	}
}

// 远程获取agent端的信息进行初始化
func (a *Agent) Init() error {
	// TODO: 此处绑定所有的消息处理函数
	a.bindHandler(protocol.Heartbeat, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})
	a.bindHandler(protocol.FileMonitor, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process file monitor from processor:%s", msg.String())
		WARN_MSG <- msg.Data.(string)
		return nil
	})

	a.bindHandler(protocol.AgentInfo, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})

	a.bindHandler(protocol.ConfigFileMonitor, func(a *Agent, msg *protocol.Message) error {
		logger.Info("remote addr:%s,process config file monitor from processor:%s",
			a.conn.RemoteAddr().String(), msg.String())
		ConfigMessageInfo(msg.Data)
		return nil
	})

	data, err := a.AgentInfo()
	if err != nil {
		logger.Error("fail to get agent info, address:%s", a.conn.RemoteAddr().String())
		return err
	}

	a.UUID = data.AgentUUID
	a.IP = data.IP
	a.Version = data.AgentVersion

	return nil
}

func (a *Agent) sendMessage(msg *protocol.Message, wait bool, timeout time.Duration) (*protocol.Message, error) {
	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}
	if wait {
		waitChan := make(chan *protocol.Message)
		a.MessageProcesser.WaitMap.Store(msg.UUID, waitChan)

		// send message to data send channel
		a.messageChan <- msg

		// wail for response
		data := <-waitChan
		return data, nil
	}

	// just send message to channel
	a.messageChan <- msg
	return nil, nil
}

// 心跳
func (a *Agent) HeartBeat() (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.Heartbeat,
		Data: "connection is normal",
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// 开启定时任务
func (a *Agent) CronStart(id int, spec string, command string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CronStart,
		Data: strconv.Itoa(id) + "," + spec + "," + command,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), resp_message.Error, nil
}

// 暂停定时任务
func (a *Agent) CronStopAndDel(id int) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CronStopAndDel,
		Data: strconv.Itoa(id),
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// 监控配置文件
func (a *Agent) ConfigfileInfo(ConMess global.ConfigMessage) error {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentConfig,
		Data: ConMess,
	}
	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to config on agent")
		return err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to config on agent: %s", resp_message.Error)
		return fmt.Errorf(resp_message.Error)
	}

	return nil
}

// 监控文件信息回传
func ConfigMessageInfo(Data interface{}) {
	p, ok := Data.(map[string]interface{})
	if ok {
		cf := configservice.ConfigFile{
			MachineUUID: p["Machine_uuid"].(string),
			Content:     p["ConfigContent"].(string),
			Path:        p["ConfigName"].(string),
			UpdatedAt:   time.Time{},
		}
		err := configservice.AddConfigFile(cf)
		if err != nil {
			logger.Error("配置文件添加失败" + err.Error())
		}
	}
}
