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

// 远程获取agent端的网络连接信息
func (a *Agent) GetNetWorkConnectInfo() (*map[string]string, string, error) {
	info := &map[string]string{}
	responseMessage, err := a.SendMessageWrapper(protocol.GetNetWorkConnectInfo, struct{}{}, "failed to run script on agent", -1, info, "GetNetWorkConnectInfo")
	return info, responseMessage.(*protocol.Message).Error, err
}

// 获取agent的基础网络配置
func (a *Agent) GetNetWorkConnInfo() (*common.NetworkConfig, string, error) {
	info := &common.NetworkConfig{}
	responseMessage, err := a.SendMessageWrapper(protocol.GetNetWorkConnInfo, struct{}{}, "failed to run script on agent", -1, info, "GetNetWorkConnInfo")
	return info, responseMessage.(*protocol.Message).Error, err
}

// 获取网卡名字
func (a *Agent) GetNICName() (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.GetNICName, struct{}{}, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), responseMessage.(*protocol.Message).Error, err
}

// 重启网卡配置
func (a *Agent) RestartNetWork(NIC string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.RestartNetWork, NIC, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Error, err
}
