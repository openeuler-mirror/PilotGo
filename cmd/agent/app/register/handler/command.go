/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package handler

import (
	"encoding/base64"
	"path"
	"strings"

	"github.com/google/uuid"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func RunCommandHandler(c *network.SocketClient, msg *protocol.Message) error {
	d := &struct {
		Command string
	}{}

	err := msg.BindData(d)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "parse data error:" + err.Error(),
		}
		return c.Send(resp_msg)
	}

	content, err := base64.StdEncoding.DecodeString(d.Command)
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "run command error:" + err.Error(),
		}
		return c.Send(resp_msg)
	}

	logger.Debug("process run command:%s", string(content))

	retCode, stdout, stderr, err := utils.RunCommand(string(content))
	if err != nil {
		resp_msg := &protocol.Message{
			UUID:   msg.UUID,
			Type:   msg.Type,
			Status: -1,
			Error:  "run command error:" + err.Error(),
		}
		return c.Send(resp_msg)
	}

	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
		Data: &utils.CmdResult{
			RetCode: retCode,
			Stdout:  stdout,
			Stderr:  stderr,
		},
	}
	return c.Send(resp_msg)
}

func RunScriptHandler(c *network.SocketClient, msg *protocol.Message) error {
	errorInfo := ""
	resp_msg := &protocol.Message{
		UUID:   msg.UUID,
		Type:   msg.Type,
		Status: 0,
	}

	var result *utils.CmdResult
	var err error
	var decoded_script []byte
	workDir := "/opt/PilotGo/agent/"
	fileName := uuid.New().String()
	filePath := path.Join(workDir, fileName+".sh")

	d := &struct {
		Script string
		Params []string
	}{}

	err = msg.BindData(d)
	if err != nil {
		errorInfo = "parse data error:" + err.Error()
		logger.Error("%s", errorInfo)
		goto ERROR
	}

	decoded_script, err = base64.StdEncoding.DecodeString(d.Script)
	if err != nil {
		errorInfo = "Err decoding base64: " + err.Error()
		logger.Error("%s", errorInfo)
		goto ERROR
	}

	logger.Debug("process run script command: %s %v", filePath+" ", d.Params)

	err = utils.FileSaveString(filePath, strings.Replace(string(decoded_script), "\r", "", -1))
	if err != nil {
		errorInfo = "Err running filesavestring:" + err.Error()
		logger.Error("%s", errorInfo)
		goto ERROR
	}

	result, err = utils.RunScript(filePath, d.Params)
	if err != nil {
		errorInfo = "run command error:" + err.Error()
		logger.Error("%s", errorInfo)
		goto ERROR
	}

	resp_msg.Status = 0
	resp_msg.Data = result
	return c.Send(resp_msg)

ERROR:
	resp_msg.Status = -1
	resp_msg.Error = errorInfo
	return c.Send(resp_msg)
}
