/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/common"
	config "gitee.com/openeuler/PilotGo/cmd/server/app/service/configmanage"
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

	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := file.ConfigFileSearchPaged(file.Search, num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
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
	num := query.Size * (query.CurrentPageNum - 1)
	total, data, err := files.HistoryConfigFilesPaged(FileId, num, query.Size)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, data, total, query)
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
