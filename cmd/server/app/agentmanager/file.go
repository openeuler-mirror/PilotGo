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

	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	sdkcommon "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// 查看配置文件内容
func (a *Agent) ReadFilePattern(filepath, pattern string) ([]sdkcommon.File, string, error) {
	responseMessage, err := a.SendMessageWrapper(protocol.ReadFilePattern, sdkcommon.File{Path: filepath, Name: pattern}, "failed to run script on agent", -1, nil, "")
	data, ok := responseMessage.(protocol.Message).Data.([]interface{})
	if !ok {
		logger.Error("failed to get msg data on agent: %s", responseMessage.(protocol.Message).Error)
		return nil, responseMessage.(protocol.Message).Error, errors.New("failed to get msg data")
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
	return files, responseMessage.(protocol.Message).Error, err
}

// 更新配置文件
func (a *Agent) UpdateFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		Path: filepath,
		Name: filename,
		Text: text,
	}
	info := &common.UpdateFile{}
	responseMessage, err := a.SendMessageWrapper(protocol.EditFile, updatefile, "failed to run script on agent", -1, info, "UpdateFile")
	return info, responseMessage.(protocol.Message).Error, err
}

// 存储配置文件
func (a *Agent) SaveFile(filepath string, filename string, text string) (*common.UpdateFile, string, error) {
	updatefile := common.UpdateFile{
		Path: filepath,
		Name: filename,
		Text: text,
	}
	info := &common.UpdateFile{}
	responseMessage, err := a.SendMessageWrapper(protocol.SaveFile, updatefile, "failed to run script on agent", -1, info, "UpdateFile")
	return info, responseMessage.(protocol.Message).Error, err
}
