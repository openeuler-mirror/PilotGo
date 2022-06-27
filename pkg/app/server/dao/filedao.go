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

	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
)

func IsExistId(id int) bool {
	var file model.Files
	global.PILOTGO_DB.Where("id=?", id).Find(&file)
	return file.ID != 0
}

func IsExistFile(filename string) bool {
	var file model.Files
	global.PILOTGO_DB.Where("file_name = ?", filename).Find(&file)
	return file.ID != 0
}

func IsExistFileLatest(fileId int) (bool, int, string) {
	var files []model.HistoryFiles
	global.PILOTGO_DB.Order("id desc").Where("file_id = ?", fileId).Find(&files)
	for _, file := range files {
		if ok := strings.Contains(file.FileName, "latest"); ok {
			return true, file.ID, file.FileName
		}
	}
	return false, 0, ""
}

func SaveHistoryFile(id int) {
	var file model.Files
	global.PILOTGO_DB.Where("id=?", id).Find(&file)

	lastversion := model.HistoryFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName,
		Description: file.Description,
		File:        file.File,
	}
	global.PILOTGO_DB.Save(&lastversion)
}

func SaveLatestFile(id int) {
	var file model.Files
	global.PILOTGO_DB.Where("id = ?", id).Find(&file)

	lastversion := model.HistoryFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName + "-latest",
		Description: file.Description,
		File:        file.File,
	}
	global.PILOTGO_DB.Save(&lastversion)
}

func UpdateFile(id int, f model.Files) {
	var file model.Files
	global.PILOTGO_DB.Model(&file).Where("id = ?", id).Update(&f)
}

func UpdateLastFile(id int, f model.HistoryFiles) {
	var file model.HistoryFiles
	global.PILOTGO_DB.Model(&file).Where("id = ?", id).Update(&f)
}

func RollBackFile(id int, text string) {
	var file model.Files
	fd := model.Files{
		File: text,
	}
	global.PILOTGO_DB.Model(&file).Where("id = ?", id).Update(&fd)
}
func DeleteFile(id int) {
	var file model.Files
	global.PILOTGO_DB.Where("id = ?", id).Unscoped().Delete(file)
}

func DeleteHistoryFile(filePId int) {
	var file model.HistoryFiles
	global.PILOTGO_DB.Where("file_id = ?", filePId).Unscoped().Delete(file)
}

func SaveFile(file model.Files) {
	global.PILOTGO_DB.Save(&file)
}

func FileText(id int) (text string) {
	file := model.Files{}
	global.PILOTGO_DB.Where("id = ?", id).Find(&file)
	return file.File
}
func LastFileText(id int) (text string) {
	file := model.HistoryFiles{}
	global.PILOTGO_DB.Where("id = ?", id).Find(&file)
	return file.File
}
func FindLastVersionFile(uuid, filename string) []model.HistoryFiles {
	var files []model.HistoryFiles
	var lastfiles []model.HistoryFiles

	global.PILOTGO_DB.Where("uuid = ? ", uuid).Find(&files)
	for _, file := range files {
		if ok := strings.Contains(file.FileName, filename); ok {
			lastfiles = append(lastfiles, file)
		}
	}
	return lastfiles
}
