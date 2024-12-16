/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentmanager

import (
	"fmt"
	"net"
	"strconv"
	"time"

	eventSDK "gitee.com/openeuler/PilotGo-plugins/event/sdk"
	"gitee.com/openeuler/PilotGo/cmd/agent/app/global"
	configservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/configfile"
	machineservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	pnet "gitee.com/openeuler/PilotGo/pkg/utils/message/net"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	commonSDK "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
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
			// 发布“平台主机离线”事件
			msgData := commonSDK.MessageData{
				MsgType:     eventSDK.MsgHostOffline,
				MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgHostOffline),
				TimeStamp:   time.Now(),
				Data: eventSDK.MDHostChange{
					IP:     a.IP,
					Status: machineservice.OfflineStatus,
				},
			}
			msgDataString, _ := msgData.ToMessageDataString()
			ms := commonSDK.EventMessage{
				MessageType: eventSDK.MsgHostOffline,
				MessageData: msgDataString,
			}
			plugin.PublishEvent(ms)
			// 消息推送到前端
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

func (a *Agent) sendMessage(msg *protocol.Message, wait bool) (*protocol.Message, error) {
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

func (a *Agent) SendMessageWrapper(protocolType int, msgData interface{}, errorMsg string, statusType int, info interface{}, bindErrorString string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocolType,
		Data: msgData,
	}

	responseMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error(errorMsg)
		return "", err
	}
	switch statusType {
	case -1:
		if responseMessage.Status == -1 || responseMessage.Error != "" {
			logger.Error(errorMsg+": %s", responseMessage.Error)
			return "", fmt.Errorf(responseMessage.Error)
		}
	case 0:
		if responseMessage.Status == 0 {
			//当状态为0时，表示命令执行成功，可以解析返回的数据。状态为-1的时候不会有数据
			result := &utils.CmdResult{}
			err = responseMessage.BindData(result)
			if err != nil {
				return nil, fmt.Errorf("failed to bind command result: %v", err)
			}
			return result, nil
		}
	}

	if info != nil {
		err = responseMessage.BindData(info)
		if err != nil {
			logger.Error("bind "+bindErrorString+" data error: %v", err)
			return nil, err
		}
	}

	return responseMessage, nil
}

// 心跳
func (a *Agent) HeartBeat() (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.Heartbeat, "connection is normal", "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}

// 开启定时任务
func (a *Agent) CronStart(id int, spec string, command string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.CronStart, strconv.Itoa(id)+","+spec+","+command, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.Error, err
}

// 暂停定时任务
func (a *Agent) CronStopAndDel(id int) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.CronStopAndDel, strconv.Itoa(id), "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), err
}

// 监控配置文件
func (a *Agent) ConfigfileInfo(ConMess global.ConfigMessage) error {
	_, err := a.SendMessageWrapper(protocol.AgentConfig, ConMess, "failed to config on agent", -1, nil, "")
	return err
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
