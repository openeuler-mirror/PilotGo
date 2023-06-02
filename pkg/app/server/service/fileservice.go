/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-05-26 10:25:52
 * LastEditTime: 2022-06-02 10:16:10
 * Description: agent config file service
 ******************************************************************************/
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

type Files = dao.Files
type SearchFile = dao.SearchFile
type HistoryFiles = dao.HistoryFiles

type DeleteFiles struct {
	FileIDs []int `json:"ids"`
}

type RollBackFiles struct {
	HistoryFileID int    `json:"id"`
	FileID        int    `json:"filePId"`
	UserUpdate    string `json:"user"`
	UserDept      string `json:"userDept"`
}

type FileBroadcast struct {
	BatchId  []int  `json:"batches"`
	Path     string `json:"path"`
	FileName string `json:"name"`
	User     string `json:"user"`
	UserDept string `json:"userDept"`
	Text     string `json:"file"`
}

// 获取时间的日期函数 => 20200426-17:36:04
func NowTime() string {
	time := time.Now()
	year := time.Year()
	month := time.Month()
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	nowtime := fmt.Sprintf("%d%02d%02d-%02d:%02d:%02d", year, month, day, hour, minute, second)
	return nowtime
}

func SaveFileToDatabase(file *Files) error {
	filename := file.FileName
	if len(filename) == 0 {
		return errors.New("请输入配置文件名字")
	}

	filepath := file.FilePath
	if len(filepath) == 0 {
		return errors.New("请输入下发文件路径")
	}
	temp, err := dao.IsExistFile(filename)
	if err != nil {
		return err
	}
	if temp {
		return errors.New("文件名字已存在，请重新输入")
	}

	filetype := file.Type
	if len(filetype) == 0 {
		return errors.New("请选择文件类型")
	}

	description := file.Description
	if len(description) == 0 {
		return errors.New("请添加文件描述")
	}

	batchId := file.ControlledBatch

	text := file.File
	if len(text) == 0 {
		return errors.New("请重新检查文件内容")
	}

	fd := Files{
		UserUpdate:      file.UserUpdate,
		UserDept:        file.UserDept,
		FileName:        filename,
		FilePath:        filepath,
		Type:            filetype,
		Description:     description,
		ControlledBatch: batchId,
		TakeEffect:      file.TakeEffect,
		File:            text,
	}
	return dao.SaveFile(fd)
}

func DeleteFile(fileIds []int) error {
	for _, fileId := range fileIds {
		err := dao.DeleteFile(fileId)
		if err != nil {
			logger.Error(err.Error())
		}
		err = dao.DeleteHistoryFile(fileId)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}
func UpdateFile(file *Files) error {
	id := file.ID
	err := dao.SaveHistoryFile(id)
	if err != nil {
		return err
	}
	user := file.UserUpdate
	userDept := file.UserDept
	filename := file.FileName
	description := file.Description
	batchId := file.ControlledBatch
	text := file.File
	ExistIdBool, err := dao.IsExistId(file.ID)
	if err != nil {
		return err
	}
	if !ExistIdBool {
		return errors.New("id有误,请重新确认该文件是否存在")
	}
	if ok, lastfileId, fileName, err := dao.IsExistFileLatest(id); ok {
		if err != nil {
			return err
		}
		fname := strings.Split(fileName, "-")
		f := dao.HistoryFiles{
			FileName: fname[0],
		}
		err = dao.UpdateLastFile(lastfileId, f)
		if err != nil {
			return err
		}
	}
	f := Files{
		Type:            file.Type,
		FileName:        filename,
		FilePath:        file.FilePath,
		Description:     description,
		UserUpdate:      user,
		UserDept:        userDept,
		ControlledBatch: batchId,
		TakeEffect:      file.TakeEffect,
		File:            text,
	}
	return dao.UpdateFile(id, f)
}

func LastFileRollBack(file *RollBackFiles) error {
	lastfileId := file.HistoryFileID
	fileId := file.FileID
	user := file.UserUpdate
	userDept := file.UserDept
	lastfileText, err := dao.LastFileText(lastfileId)
	if err != nil {
		return err
	}
	if ok, _, _, err := dao.IsExistFileLatest(fileId); !ok {
		if err != nil {
			return err
		}
		err := dao.SaveLatestFile(fileId)
		if err != nil {
			return err
		}
	}
	fd := Files{
		UserUpdate: user,
		UserDept:   userDept,
		File:       lastfileText,
	}
	return dao.UpdateFile(fileId, fd)
}
