/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
)

type ConfigFiles struct {
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

type HistoryConfigFiles struct {
	ID          int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	FileID      int `json:"filePId"`
	UpdatedAt   time.Time
	UserUpdate  string `json:"user"`
	UserDept    string `json:"userDept"`
	FileName    string `json:"name"`
	Description string `json:"description"`
	File        string `gorm:"type:text" json:"file"`
}

type SearchConfigFile struct {
	Search string `json:"search"`
}

// 分页查询
func GetConfigFilesPaged(offset, size int) (int64, []ConfigFiles, error) {
	var count int64
	var configFiles []ConfigFiles
	err := mysqlmanager.MySQL().Model(ConfigFiles{}).Order("id desc").Offset(offset).Limit(size).Find(&configFiles).Offset(-1).Limit(-1).Count(&count).Error
	return count, configFiles, err
}

func (f *SearchConfigFile) ConfigFileSearchPaged(search string, offset, size int) (int64, []ConfigFiles, error) {
	var count int64
	var configFiles []ConfigFiles
	err := mysqlmanager.MySQL().Model(ConfigFiles{}).Order("id desc").Where("file_name LIKE ?", "%"+search+"%").Offset(offset).Limit(size).Find(&configFiles).Offset(-1).Limit(-1).Count(&count).Error
	return count, configFiles, err
}

func (f *HistoryConfigFiles) HistoryConfigFilesPaged(fileId, offset, size int) (int64, []HistoryConfigFiles, error) {
	var count int64
	var historyConfigFiles []HistoryConfigFiles
	err := mysqlmanager.MySQL().Model(HistoryConfigFiles{}).Order("id desc").Where("file_id=?", fileId).Offset(offset).Limit(size).Find(&historyConfigFiles).Offset(-1).Limit(-1).Count(&count).Error
	return count, historyConfigFiles, err
}

func IsExistId(id int) (bool, error) {
	var file ConfigFiles
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&file).Error
	return file.ID != 0, err
}

func IsExistConfigFile(filename string) (bool, error) {
	var file ConfigFiles
	err := mysqlmanager.MySQL().Where("file_name = ?", filename).Find(&file).Error
	return file.ID != 0, err
}

func IsExistConfigFileLatest(fileId int) (bool, int, string, error) {
	var files []HistoryConfigFiles
	err := mysqlmanager.MySQL().Order("id desc").Where("file_id = ?", fileId).Find(&files).Error
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

func SaveHistoryConfigFile(id int) error {
	var file ConfigFiles
	err := mysqlmanager.MySQL().Where("id=?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := HistoryConfigFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName,
		Description: file.Description,
		File:        file.File,
	}
	return mysqlmanager.MySQL().Save(&lastversion).Error
}

func SaveLatestConfigFile(id int) error {
	var file ConfigFiles
	err := mysqlmanager.MySQL().Where("id = ?", id).Find(&file).Error
	if err != nil {
		return err
	}
	lastversion := HistoryConfigFiles{
		FileID:      id,
		UserUpdate:  file.UserUpdate,
		UserDept:    file.UserDept,
		FileName:    file.FileName + "-latest",
		Description: file.Description,
		File:        file.File,
	}
	return mysqlmanager.MySQL().Save(&lastversion).Error
}

func UpdateConfigFile(id int, f ConfigFiles) error {
	var file ConfigFiles
	return mysqlmanager.MySQL().Model(&file).Where("id = ?", id).Updates(&f).Error
}

func UpdateLastConfigFile(id int, f HistoryConfigFiles) error {
	var file HistoryConfigFiles
	return mysqlmanager.MySQL().Model(&file).Where("id = ?", id).Updates(&f).Error
}

func RollBackFile(id int, text string) error {
	var file ConfigFiles
	fd := ConfigFiles{
		File: text,
	}
	return mysqlmanager.MySQL().Model(&file).Where("id = ?", id).Updates(&fd).Error
}
func DeleteConfigFile(id int) error {
	var file ConfigFiles
	return mysqlmanager.MySQL().Where("id = ?", id).Unscoped().Delete(file).Error
}

func DeleteHistoryConfigFile(filePId int) error {
	var file HistoryConfigFiles
	return mysqlmanager.MySQL().Where("file_id = ?", filePId).Unscoped().Delete(file).Error
}

func SaveConfigFile(file ConfigFiles) error {
	return mysqlmanager.MySQL().Save(&file).Error
}

func ConfigFileText(id int) (text string, err error) {
	file := ConfigFiles{}
	err = mysqlmanager.MySQL().Where("id = ?", id).Find(&file).Error
	return file.File, err
}

func LastConfigFileText(id int) (text string, err error) {
	file := HistoryConfigFiles{}
	err = mysqlmanager.MySQL().Where("id = ?", id).Find(&file).Error
	return file.File, err
}

func FindLastVersionConfigFile(uuid, filename string) ([]HistoryConfigFiles, error) {
	var files []HistoryConfigFiles
	var lastfiles []HistoryConfigFiles

	err := mysqlmanager.MySQL().Where("uuid = ? ", uuid).Find(&files).Error
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
