/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentmanager

import (
	"encoding/base64"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
)

// 远程在agent上运行shell命令
func (a *Agent) RunCommand(cmd string) (*utils.CmdResult, error) {
	data := struct {
		Command string
	}{
		Command: cmd,
	}
	responseMessage, err := a.SendMessageWrapper(protocol.RunCommand, data, "failed to run command on agent", 0, nil, "")
	return responseMessage.(*utils.CmdResult), err
}

// 远程在agent上运行脚本文件
func (a *Agent) RunScript(script string, params []string) (*utils.CmdResult, error) {
	encoded_script := base64.StdEncoding.EncodeToString([]byte(script))
	data := struct {
		Script string
		Params []string
	}{
		Script: encoded_script,
		Params: params,
	}
	responseMessage, err := a.SendMessageWrapper(protocol.RunScript, data, "failed to run script on agent", 0, nil, "")
	result, ok := responseMessage.(*utils.CmdResult)
	if !ok {
		return &utils.CmdResult{}, err
	}
	return result, err
}

// chmod [-R] 权限值 文件名
func (a *Agent) ChangePermission(permission, file string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ChangePermission, permission+","+file, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), err
}

// chown [-R] 所有者 文件或目录
func (a *Agent) ChangeFileOwner(user, file string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ChangeFileOwner, user+","+file, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), err
}

// 临时修改agent端系统参数
func (a *Agent) ChangeSysctl(args string) (string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.SysctlChange, args, "failed to run script on agent", -1, nil, "")
	return responseMessage.(*protocol.Message).Data.(string), err
}
