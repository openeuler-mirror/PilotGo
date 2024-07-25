package agentmanager

import (
	"errors"
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

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run command on agent")
		return nil, err
	}

	if resp_message.Status == 0 {
		//状态为-1的时候不会有数据
		result := &utils.CmdResult{}
		err = resp_message.BindData(result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New(resp_message.Error)
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

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}

	if resp_message.Status == 0 {
		result := &utils.CmdResult{}
		err = resp_message.BindData(result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New(resp_message.Error)
}

// chmod [-R] 权限值 文件名
func (a *Agent) ChangePermission(permission, file string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangePermission,
		Data: permission + "," + file,
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// chown [-R] 所有者 文件或目录
func (a *Agent) ChangeFileOwner(user, file string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangeFileOwner,
		Data: user + "," + file,
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}

// 临时修改agent端系统参数
func (a *Agent) ChangeSysctl(args string) (string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlChange,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return "", fmt.Errorf(resp_message.Error)
	}

	return resp_message.Data.(string), nil
}
