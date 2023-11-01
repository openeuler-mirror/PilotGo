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
 * Date: 2022-06-16 10:25:52
 * LastEditTime: 2022-06-16 16:16:10
 * Description: file info handler
 ******************************************************************************/
package controller

import (
	"strconv"

	"gitee.com/openeuler/PilotGo/app/server/service/common"
	config "gitee.com/openeuler/PilotGo/app/server/service/configmanage"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func SaveConfigFileToDatabaseHandler(c *gin.Context) {
	var file config.ConfigFiles
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := config.SaveToDatabase(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "文件保存成功")
}

func DeleteConfigFileHandler(c *gin.Context) {
	var files config.DeleteConfigFiles
	if err := c.Bind(&files); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	fileids := files.FileIDs
	err := config.DeleteConfig(fileids)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "储存的文件已从数据库中删除")
}

func UpdateConfigFileHandler(c *gin.Context) {
	var file config.ConfigFiles
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := config.UpdateConfig(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "配置文件修改成功")
}

func AllConfigFiles(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := config.GetConfigFilesPaged(num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
}

func ConfigFileSearchHandler(c *gin.Context) {
	var file config.SearchConfigFile
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	search := file.Search

	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	list, tx := file.ConfigFileSearch(search)

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, list, total, query)
}

func HistoryConfigFilesHandler(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	fileId := c.Query("id")
	FileId, err := strconv.Atoi(fileId)
	if err != nil {
		response.Fail(c, nil, "文件ID输入格式有误")
		return
	}

	files := config.HistoryConfigFiles{}
	list, tx := files.HistoryConfigFiles(FileId)

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, list, total, query)
}

func LastConfigFileRollBackHandler(c *gin.Context) {
	var file config.RollBackConfigFiles
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := config.LastConfigFileRollBack(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "已回退到历史版本")
}
