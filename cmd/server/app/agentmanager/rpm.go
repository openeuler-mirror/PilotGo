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

	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
)

// 获取全部安装的rpm包列表
func (a *Agent) AllRpm() ([]string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.AllRpm, struct{}{}, "failed to run script on agent", -1, nil, "")

	if v, ok := responseMessage.(protocol.Message).Data.([]interface{}); ok {
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			}
		}
		return result, err
	}
	return nil, fmt.Errorf("failed to convert interface{} in allrpm")
}

// 获取源软件包名以及源
func (a *Agent) RpmSource(rpm string) (*common.RpmSrc, error) {
	info := &common.RpmSrc{}
	_, err := a.SendMessageWrapper(protocol.RpmSource, rpm, "failed to run script on agent", -1, info, "RpmSource")
	return info, err
}

// 获取软件包信息
func (a *Agent) RpmInfo(rpm string) (*common.RpmInfo, string, error) {
	info := &common.RpmInfo{}
	responseMessage, err := a.SendMessageWrapper(protocol.RpmInfo, rpm, "failed to run script on agent", -1, info, "RpmInfo")
	return info, responseMessage.(protocol.Message).Error, err
}

// 安装软件包
func (a *Agent) InstallRpm(rpm string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.InstallRpm, rpm, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.(protocol.Message).Error, err
}

// 卸载软件包
func (a *Agent) RemoveRpm(rpm string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.RemoveRpm, rpm, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.(protocol.Message).Error, err
}
