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
 * Description: agent config file struct
 ******************************************************************************/
package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

type Files struct {
	ID          int `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   time.Time
	SourcePath  string `json:"path"`
	FileName    string `json:"name"`
	Description string `json:"description"`
	File        string `gorm:"type:text" json:"file"`
}

type HistoryFiles struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	IP       string `json:"ip"`
	IPDept   string `json:"ipDept"`
	UUID     string `json:"uuid"`
	Path     string `json:"path"`
	FileName string `json:"name"`
	File     string `gorm:"type:text" json:"file"`
}

type DeleteFiles struct {
	FileIDs []int `json:"ids"`
}

type SearchFile struct {
	Search string `json:"search"`
}

type FileBroadcast struct {
	UUID     []string `json:"uuids"`
	Path     string   `json:"path"`
	FileName string   `json:"name"`
	User     string   `json:"user"`
	UserDept string   `json:"userDept"`
	Text     string   `json:"file"`
}

func (f *Files) AllFiles(q *PaginationQ) (list *[]Files, tx *gorm.DB) {
	list = &[]Files{}
	tx = mysqlmanager.DB.Order("created_at desc").Find(&list)
	return
}

func (f *SearchFile) FileSearch(q *PaginationQ, search string) (list *[]Files, tx *gorm.DB) {
	list = &[]Files{}
	tx = mysqlmanager.DB.Order("created_at desc").Where("source_path LIKE ?", "%"+search+"%").Find(&list)
	if len(*list) == 0 {
		tx = mysqlmanager.DB.Order("created_at desc").Where("file_name LIKE ?", "%"+search+"%").Find(&list)
	}
	return
}

func (f *HistoryFiles) AllHistoryFiles(q *PaginationQ) (list *[]HistoryFiles, tx *gorm.DB) {
	list = &[]HistoryFiles{}
	tx = mysqlmanager.DB.Order("id desc").Find(&list)
	return
}
func (f *SearchFile) LastFileSearch(q *PaginationQ, search string) (list *[]HistoryFiles, tx *gorm.DB) {
	list = &[]HistoryFiles{}
	tx = mysqlmanager.DB.Order("id desc").Where("path LIKE ?", "%"+search+"%").Find(&list)
	if len(*list) == 0 {
		tx = mysqlmanager.DB.Order("id desc").Where("ip LIKE ?", "%"+search+"%").Find(&list)
		if len(*list) == 0 {
			tx = mysqlmanager.DB.Order("id desc").Where("file_name LIKE ?", "%"+search+"%").Find(&list)
		}
	}
	return
}
