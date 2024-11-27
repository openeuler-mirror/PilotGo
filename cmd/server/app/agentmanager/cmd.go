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

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/google/uuid"
)

// 远程在agent上运行shell命令
func (a *Agent) RunCommand(cmd string) (*utils.CmdResult, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RunCommand,
		Data: struct {
			Command string
		}{
			Command: cmd,
		},
	}

	respMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run command on agent: %v", err)
		return nil, err
	}

	if respMessage.Status == 0 {
		//当状态为0时，表示命令执行成功，可以解析返回的数据。状态为-1的时候不会有数据
		result := &utils.CmdResult{}
		err = respMessage.BindData(result)
		if err != nil {
			return nil, fmt.Errorf("failed to bind command result: %v", err)
		}
		return result, nil
	}
	return nil, fmt.Errorf("agent returned error: %s", respMessage.Error)
}

// 远程在agent上运行脚本文件
func (a *Agent) RunScript(script string, params []string) (*utils.CmdResult, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RunScript,
		Data: struct {
			Script string
			Params []string
		}{
			Script: script,
			Params: params,
		},
	}

	respMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent: %v", err)
		return nil, err
	}

	if respMessage.Status == 0 {
		result := &utils.CmdResult{}
		err = respMessage.BindData(result)
		if err != nil {
			return nil, fmt.Errorf("failed to bind command result: %v", err)
		}
		return result, nil
	}
	return nil, fmt.Errorf("agent returned error: %s", respMessage.Error)
}

// chmod [-R] 权限值 文件名
func (a *Agent) ChangePermission(permission, file string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangePermission,
		Data: permission + "," + file,
	}

	respMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if respMessage.Status == -1 || respMessage.Error != "" {
		logger.Error("failed to run script on agent: %s", respMessage.Error)
		return "", fmt.Errorf(respMessage.Error)
	}

	return respMessage.Data.(string), nil
}

// chown [-R] 所有者 文件或目录
func (a *Agent) ChangeFileOwner(user, file string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangeFileOwner,
		Data: user + "," + file,
	}

	respMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if respMessage.Status == -1 || respMessage.Error != "" {
		logger.Error("failed to run script on agent: %s", respMessage.Error)
		return "", fmt.Errorf(respMessage.Error)
	}

	return respMessage.Data.(string), nil
}

// 临时修改agent端系统参数
func (a *Agent) ChangeSysctl(args string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlChange,
		Data: args,
	}

	respMessage, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if respMessage.Status == -1 || respMessage.Error != "" {
		logger.Error("failed to run script on agent: %s", respMessage.Error)
		return "", fmt.Errorf(respMessage.Error)
	}

	return respMessage.Data.(string), nil
}
