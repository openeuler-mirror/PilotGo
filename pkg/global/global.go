/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Mon Jun 27 17:13:49 2022 +0800
 */
package global

import "gitee.com/openeuler/PilotGo/sdk/logger"

const (
	// 新注册机器添加到部门根节点
	UncateloguedDepartId = 1
	// 是否为部门根节点
	Departroot   = 0
	DepartUnroot = 1

	// 集群规模
	ClusterSize = 10
)

var WARN_MSG chan *WebsocketSendMsg = make(chan *WebsocketSendMsg, 100)

type WebsocketSendMsgType int

const (
	MachineSendMsg WebsocketSendMsgType = iota
	PluginSendMsg
	ServerSendMsg
)

type WebsocketSendMsg struct {
	MsgType WebsocketSendMsgType `json:"msgtype"`
	Msg     string               `json:"msg"`
}

/*
	WebsocketSendMsgType: pkg/global/global.go
*/
func SendRemindMsg(msg_type WebsocketSendMsgType, msg string) {
	defer func ()  {
		if r := recover(); r != nil {
			logger.Error("send remind message to closed channel WARN_MSG: %+v", r)	
		}	
	}()

	WARN_MSG <- &WebsocketSendMsg{
		MsgType: msg_type,
		Msg:     msg,
	}
}