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
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	fileservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/file"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func SaveFileToDatabaseHandler(c *gin.Context) {
	var file fileservice.Files
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := fileservice.SaveToDatabase(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "文件保存成功")
}

func DeleteFileHandler(c *gin.Context) {
	var files fileservice.DeleteFiles
	if err := c.Bind(&files); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	fileids := files.FileIDs
	err := fileservice.Delete(fileids)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "储存的文件已从数据库中删除")
}

func UpdateFileHandler(c *gin.Context) {
	var file fileservice.Files
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := fileservice.Update(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "配置文件修改成功")
}

func AllFiles(c *gin.Context) {
	query := &common.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	files := fileservice.Files{}
	list, tx := files.AllFiles()

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	var filetype []string
	filetype = append(filetype, global.ConfigRepo)

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"ok":    true,
		"data":  list,
		"total": total,
		"page":  query.CurrentPageNum,
		"size":  query.Size,
		"type":  filetype})
}

func FileSearchHandler(c *gin.Context) {
	var file fileservice.SearchFile
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

	list, tx := file.FileSearch(search)

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, list, total, query)
}

func HistoryFilesHandler(c *gin.Context) {
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

	files := fileservice.HistoryFiles{}
	list, tx := files.HistoryFiles(FileId)

	total, err := common.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	common.JsonPagination(c, list, total, query)
}

func LastFileRollBackHandler(c *gin.Context) {
	var file fileservice.RollBackFiles
	if err := c.Bind(&file); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	err := fileservice.LastFileRollBack(&file)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "已回退到历史版本")
}
