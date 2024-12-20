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

// 获取当前用户信息
func (a *Agent) CurrentUser() (*common.CurrentUser, error) {
	info := &common.CurrentUser{}
	_, err := a.SendMessageWrapper(protocol.CurrentUser, struct{}{}, "failed to run script on agent", -1, info, "")
	return info, err
}

// 获取所有用户的信息
func (a *Agent) AllUser() ([]*common.AllUserInfo, error) {
	info := &[]*common.AllUserInfo{}
	_, err := a.SendMessageWrapper(protocol.AllUser, struct{}{}, "failed to run script on agent", -1, info, "AllUser")
	return *info, err
}

// 创建新的用户，并新建家目录
func (a *Agent) AddLinuxUser(username, password string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.AddLinuxUser, username+","+password, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), err
}

// 删除用户
func (a *Agent) DelUser(username string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.DelUser, username, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), responseMessage.(*protocol.Message).Error, err
}
