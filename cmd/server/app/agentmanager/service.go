/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentmanager

import (
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
)

// 查看服务列表
func (a *Agent) ServiceList() ([]*common.ListService, error) {
	info := &[]*common.ListService{}
	_, err := a.SendMessageWrapper(protocol.ServiceList, struct{}{}, "failed to run script on agent", -1, info, "ServiceList")
	return *info, err
}

// 查看某个服务
func (a *Agent) GetService(service string) (*common.ServiceInfo, error) {
	serviceInfo := &common.ServiceInfo{}
	_, err := a.SendMessageWrapper(protocol.GetService, struct{}{}, "failed to run script on agent", -1, serviceInfo, "GetServiceInfo")
	return serviceInfo, err
}

// 重启服务
func (a *Agent) ServiceRestart(service string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ServiceRestart, service, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), responseMessage.(*protocol.Message).Error, err
}

// 关闭服务
func (a *Agent) ServiceStop(service string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ServiceStop, service, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), responseMessage.(*protocol.Message).Error, err
}

// 启动服务
func (a *Agent) ServiceStart(service string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ServiceStart, service, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), responseMessage.(*protocol.Message).Error, err
}
