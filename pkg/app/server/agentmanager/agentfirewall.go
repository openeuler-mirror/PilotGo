package agentmanager

import (
	"fmt"

	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取防火墙配置
func (a *Agent) FirewalldConfig() (*common.FireWalldConfig, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldConfig,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := &common.FireWalldConfig{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind FirewalldConfig data error:%s", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}

// 更改防火墙默认区域
func (a *Agent) FirewalldSetDefaultZone(zone string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldDefaultZone,
		Data: zone,
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

// 查看防火墙指定区域配置
func (a *Agent) FirewalldZoneConfig(zone string) (*common.FirewalldCMDList, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldZoneConfig,
		Data: zone,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := &common.FirewalldCMDList{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind data error:%s", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}

// 添加防火墙服务
func (a *Agent) FirewalldServiceAdd(zone, service string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldServiceAdd,
		Data: zone + "," + service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Error, nil
}

// 移除防火墙服务
func (a *Agent) FirewalldServiceRemove(zone, service string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldServiceRemove,
		Data: zone + "," + service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Error, nil
}

// 防火墙添加允许来源地址
func (a *Agent) FirewalldSourceAdd(zone, source string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldSourceAdd,
		Data: zone + "," + source,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Error, nil
}

// 防火墙移除允许来源地址
func (a *Agent) FirewalldSourceRemove(zone, source string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldSourceRemove,
		Data: zone + "," + source,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Error, nil
}

// 重启防火墙
func (a *Agent) FirewalldRestart() (bool, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldRestart,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return false, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return false, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(bool), resp_message.Error, nil
}

// 关闭防火墙
func (a *Agent) FirewalldStop() (bool, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldStop,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return false, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return false, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(bool), resp_message.Error, nil
}

// 防火墙指定区域添加端口
func (a *Agent) FirewalldZonePortAdd(zone, port, proto string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldZonePortAdd,
		Data: zone + "," + port + "," + proto,
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

// 防火墙指定区域删除端口
func (a *Agent) FirewalldZonePortDel(zone, port, proto string) (string, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.FirewalldZonePortDel,
		Data: zone + "," + port + "," + proto,
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
