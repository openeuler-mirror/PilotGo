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
	"errors"
	"fmt"

	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	sdkcommon "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/google/uuid"
)

// 查看配置文件内容
func (a *Agent) ReadFilePattern(filepath, pattern string) ([]sdkcommon.File, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ReadFilePattern,
		Data: sdkcommon.File{Path: filepath, Name: pattern},
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	data, ok := resp_message.Data.([]interface{})
	if !ok {
		logger.Error("failed to get msg data on agent: %s", resp_message.Error)
		return nil, resp_message.Error, errors.New("failed to get msg data")
	}

	var files []sdkcommon.File
	for _, item := range data {
		if fileMap, ok := item.(map[string]interface{}); ok {
			f := sdkcommon.File{
				Path:    fileMap["path"].(string),
				Name:    fileMap["name"].(string),
				Content: base64.StdEncoding.EncodeToString([]byte(fileMap["content"].(string))),
			}
			files = append(files, f)
		} else {
			logger.Error("failed to get file from data")
		}
	}
	return files, resp_message.Error, nil
}

// 更新配置文件
func (a *Agent) UpdateFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		Path: filepath,
		Name: filename,
		Text: text,
	}
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.EditFile,
		Data: updatefile,
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := &common.UpdateFile{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind UpdateFile data error:%s", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}

// 存储配置文件
func (a *Agent) SaveFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		Path: filepath,
		Name: filename,
		Text: text,
	}
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SaveFile,
		Data: updatefile,
	}

	resp_message, err := a.sendMessage(msg, true)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}

	if resp_message.Status == -1 || resp_message.Error != "" {
		logger.Error("failed to run script on agent: %s", resp_message.Error)
		return nil, resp_message.Error, fmt.Errorf(resp_message.Error)
	}

	info := &common.UpdateFile{}
	err = resp_message.BindData(info)
	if err != nil {
		logger.Error("bind UpdateFile data error:%s", err)
		return nil, resp_message.Error, err
	}
	return info, resp_message.Error, nil
}
