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
 * Description: agent config file dao
 ******************************************************************************/
package dao

import (
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

func IsExistId(id int) (bool, error) {
	var file model.Files
	err := global.PILOTGO_DB.Where("id=?", id).Find(&file).Error
	if err != nil {
		return file.ID != 0, err
	}
	return file.ID != 0, nil
}

func IsExistFile(filename string) (bool, error) {
	var file model.Files
	err := global.PILOTGO_DB.Where("file_name = ?", filename).Find(&file).Error
	if err != nil {
		return file.ID != 0, err
	}
	return file.ID != 0, nil
}

func IsExistFileLatest(fileId int) (bool, int, string, error) {
	var files []model.HistoryFiles
	err := global.PILOTGO_DB.Order("id desc").Where("file_id = ?", fileId).Find(&files).Error
	if err != nil {
		return false, 0, "", err
	}
	for _, file := range files {
		if ok := strings.Contains(file.FileName, "latest"); ok {
			return true, file.ID, file.FileName, nil
		}
	}
	return false, 0, "", nil
}

func SaveHistoryFile(id int) error {
	var file model.Files
	err := global.PILOTGO_DB.Where("id=?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := model.HistoryFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName,
		Description: file.Description,
		File:        file.File,
	}
	return global.PILOTGO_DB.Save(&lastversion).Error
}

func SaveLatestFile(id int) error {
	var file model.Files
	err := global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := model.HistoryFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName + "-latest",
		Description: file.Description,
		File:        file.File,
	}
	return global.PILOTGO_DB.Save(&lastversion).Error
}

func UpdateFile(id int, f model.Files) error {
	var file model.Files
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&f).Error
}

func UpdateLastFile(id int, f model.HistoryFiles) error {
	var file model.HistoryFiles
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&f).Error
}

func RollBackFile(id int, text string) error {
	var file model.Files
	fd := model.Files{
		File: text,
	}
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&fd).Error
}
func DeleteFile(id int) error {
	var file model.Files
	return global.PILOTGO_DB.Where("id = ?", id).Unscoped().Delete(file).Error
}

func DeleteHistoryFile(filePId int) error {
	var file model.HistoryFiles
	return global.PILOTGO_DB.Where("file_id = ?", filePId).Unscoped().Delete(file).Error
}

func SaveFile(file model.Files) error {
	return global.PILOTGO_DB.Save(&file).Error
}

func FileText(id int) (text string, err error) {
	file := model.Files{}
	err = global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return file.File, err
	}
	return file.File, nil
}

func LastFileText(id int) (text string, err error) {
	file := model.HistoryFiles{}
	err = global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return file.File, err
	}
	return file.File, nil
}

func FindLastVersionFile(uuid, filename string) ([]model.HistoryFiles, error) {
	var files []model.HistoryFiles
	var lastfiles []model.HistoryFiles

	err := global.PILOTGO_DB.Where("uuid = ? ", uuid).Find(&files).Error
	if err != nil {
		return lastfiles, err
	}
	for _, file := range files {
		if ok := strings.Contains(file.FileName, filename); ok {
			lastfiles = append(lastfiles, file)
		}
	}
	return lastfiles, nil
}
