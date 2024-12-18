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

// 获取防火墙配置
func (a *Agent) FirewalldConfig() (*common.FireWalldConfig, string, error) {
	info := &common.FireWalldConfig{}
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldConfig, struct{}{}, "failed to run script on agent", -1, info, "FirewalldConfig")
	return info, responseMessage.(protocol.Message).Error, err
}

// 更改防火墙默认区域
func (a *Agent) FirewalldSetDefaultZone(zone string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldDefaultZone, zone, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.(protocol.Message).Error, err
}

// 查看防火墙指定区域配置
func (a *Agent) FirewalldZoneConfig(zone string) (*common.FirewalldCMDList, string, error) {
	info := &common.FirewalldCMDList{}
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldZoneConfig, zone, "failed to run script on agent", -1, info, "FirewalldConfig")
	return info, responseMessage.(protocol.Message).Error, err
}

// 添加防火墙服务
func (a *Agent) FirewalldServiceAdd(zone, service string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldServiceAdd, zone+","+service, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Error, err
}

// 移除防火墙服务
func (a *Agent) FirewalldServiceRemove(zone, service string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldServiceRemove, zone+","+service, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Error, err
}

// 防火墙添加允许来源地址
func (a *Agent) FirewalldSourceAdd(zone, source string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldSourceAdd, zone+","+source, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Error, err
}

// 防火墙移除允许来源地址
func (a *Agent) FirewalldSourceRemove(zone, source string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldSourceRemove, zone+","+source, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Error, err
}

// 重启防火墙
func (a *Agent) FirewalldRestart() (bool, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldRestart, struct{}{}, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(bool), responseMessage.(protocol.Message).Error, err
}

// 关闭防火墙
func (a *Agent) FirewalldStop() (bool, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldStop, struct{}{}, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(bool), responseMessage.(protocol.Message).Error, err
}

// 防火墙指定区域添加端口
func (a *Agent) FirewalldZonePortAdd(zone, port, proto string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldZonePortAdd, zone+","+port+","+proto, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.(protocol.Message).Error, err
}

// 防火墙指定区域删除端口
func (a *Agent) FirewalldZonePortDel(zone, port, proto string) (string, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.FirewalldZonePortDel, zone+","+port+","+proto, "failed to run script on agent", -1, nil, "")
	return responseMessage.(protocol.Message).Data.(string), responseMessage.(protocol.Message).Error, err
}
