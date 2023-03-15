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
	"time"

	"gorm.io/gorm"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

type Files struct {
	ID              int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	FileName        string `json:"name"`
	FilePath        string `json:"path"`
	Type            string `json:"type"`
	Description     string `json:"description"`
	UserUpdate      string `json:"user"`
	UserDept        string `json:"userDept"`
	UpdatedAt       time.Time
	ControlledBatch string `json:"batchId"`
	TakeEffect      string `json:"activeMode"`
	File            string `gorm:"type:text" json:"file"`
}

type HistoryFiles struct {
	ID          int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	FileID      int `json:"filePId"`
	UpdatedAt   time.Time
	UserUpdate  string `json:"user"`
	UserDept    string `json:"userDept"`
	FileName    string `json:"name"`
	Description string `json:"description"`
	File        string `gorm:"type:text" json:"file"`
}

type RollBackFiles struct {
	HistoryFileID int    `json:"id"`
	FileID        int    `json:"filePId"`
	UserUpdate    string `json:"user"`
	UserDept      string `json:"userDept"`
}

type DeleteFiles struct {
	FileIDs []int `json:"ids"`
}

type SearchFile struct {
	Search string `json:"search"`
}

type FileBroadcast struct {
	BatchId  []int  `json:"batches"`
	Path     string `json:"path"`
	FileName string `json:"name"`
	User     string `json:"user"`
	UserDept string `json:"userDept"`
	Text     string `json:"file"`
}

func (f *Files) AllFiles(q *PaginationQ) (list *[]Files, tx *gorm.DB) {
	list = &[]Files{}
	tx = global.PILOTGO_DB.Order("id desc").Find(&list)
	return
}

func (f *SearchFile) FileSearch(q *PaginationQ, search string) (list *[]Files, tx *gorm.DB) {
	list = &[]Files{}
	tx = global.PILOTGO_DB.Order("id desc").Where("type LIKE ?", "%"+search+"%").Find(&list)
	if len(*list) == 0 {
		tx = global.PILOTGO_DB.Order("id desc").Where("file_name LIKE ?", "%"+search+"%").Find(&list)
	}
	return
}

func (f *HistoryFiles) HistoryFiles(q *PaginationQ, fileId int) (list *[]HistoryFiles, tx *gorm.DB) {
	list = &[]HistoryFiles{}
	tx = global.PILOTGO_DB.Order("id desc").Where("file_id=?", fileId).Find(&list)
	return
}

func IsExistId(id int) (bool, error) {
	var file Files
	err := global.PILOTGO_DB.Where("id=?", id).Find(&file).Error
	if err != nil {
		return file.ID != 0, err
	}
	return file.ID != 0, nil
}

func IsExistFile(filename string) (bool, error) {
	var file Files
	err := global.PILOTGO_DB.Where("file_name = ?", filename).Find(&file).Error
	if err != nil {
		return file.ID != 0, err
	}
	return file.ID != 0, nil
}

func IsExistFileLatest(fileId int) (bool, int, string, error) {
	var files []HistoryFiles
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
	var file Files
	err := global.PILOTGO_DB.Where("id=?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := HistoryFiles{
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
	var file Files
	err := global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := HistoryFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName + "-latest",
		Description: file.Description,
		File:        file.File,
	}
	return global.PILOTGO_DB.Save(&lastversion).Error
}

func UpdateFile(id int, f Files) error {
	var file Files
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&f).Error
}

func UpdateLastFile(id int, f HistoryFiles) error {
	var file HistoryFiles
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&f).Error
}

func RollBackFile(id int, text string) error {
	var file Files
	fd := Files{
		File: text,
	}
	return global.PILOTGO_DB.Model(&file).Where("id = ?", id).Updates(&fd).Error
}
func DeleteFile(id int) error {
	var file Files
	return global.PILOTGO_DB.Where("id = ?", id).Unscoped().Delete(file).Error
}

func DeleteHistoryFile(filePId int) error {
	var file HistoryFiles
	return global.PILOTGO_DB.Where("file_id = ?", filePId).Unscoped().Delete(file).Error
}

func SaveFile(file Files) error {
	return global.PILOTGO_DB.Save(&file).Error
}

func FileText(id int) (text string, err error) {
	file := Files{}
	err = global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return file.File, err
	}
	return file.File, nil
}

func LastFileText(id int) (text string, err error) {
	file := HistoryFiles{}
	err = global.PILOTGO_DB.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return file.File, err
	}
	return file.File, nil
}

func FindLastVersionFile(uuid, filename string) ([]HistoryFiles, error) {
	var files []HistoryFiles
	var lastfiles []HistoryFiles

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
