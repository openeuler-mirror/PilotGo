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
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
)

func IsExistId(id int) bool {
	var file model.Files
	mysqlmanager.DB.Where("id=?", id).Find(&file)
	return file.ID != 0
}

func UpdateFile(id int, path string, filename string, descrip string, text string) {
	var file model.Files
	f := model.Files{
		SourcePath:  path,
		FileName:    filename,
		Description: descrip,
		File:        text,
	}
	mysqlmanager.DB.Model(&file).Where("id = ?", id).Update(&f)
}

func DeleteFile(id int) {
	var file model.Files
	mysqlmanager.DB.Where("id = ?", id).Unscoped().Delete(file)
}

func SaveFile(file model.Files) {
	mysqlmanager.DB.Save(&file)
}

func SaveHistoryFile(file model.HistoryFiles) {
	mysqlmanager.DB.Save(&file)
}
