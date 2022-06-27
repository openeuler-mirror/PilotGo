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
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func SaveFileToDatabase(c *gin.Context) {
	var file model.Files
	c.Bind(&file)

	filename := file.FileName
	if len(filename) == 0 {
		response.Fail(c, nil, "请输入配置文件名字")
		return
	}

	filepath := file.FilePath
	if len(filepath) == 0 {
		response.Fail(c, nil, "请输入下发文件路径")
		return
	}

	if !dao.IsExistFile(filename) {
		response.Fail(c, nil, "文件名字已存在，请重新输入")
		return
	}

	filetype := file.Type
	if len(filetype) == 0 {
		response.Fail(c, nil, "请选择文件类型")
		return
	}

	description := file.Description
	if len(description) == 0 {
		response.Fail(c, nil, "请添加文件描述")
		return
	}

	batchId := file.ControlledBatch

	text := file.File
	if len(text) == 0 {
		response.Fail(c, nil, "请重新检查文件内容")
		return
	}

	fd := model.Files{
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
	dao.SaveFile(fd)
	response.Success(c, nil, "文件保存成功")
}

func DeleteFile(c *gin.Context) {
	var files model.DeleteFiles
	c.Bind(&files)

	for _, fileId := range files.FileIDs {
		dao.DeleteFile(fileId)
		dao.DeleteHistoryFile(fileId)
	}
	response.Success(c, nil, "储存的文件已从数据库中删除")
}

func UpdateFile(c *gin.Context) {
	var file model.Files
	c.Bind(&file)

	id := file.ID
	dao.SaveHistoryFile(id)

	user := file.UserUpdate
	userDept := file.UserDept
	filename := file.FileName
	description := file.Description

	batchId := file.ControlledBatch
	text := file.File
	if !dao.IsExistId(file.ID) {
		response.Fail(c, nil, "id有误,请重新确认该文件是否存在")
		return
	}
	if ok, lastfileId, fileName := dao.IsExistFileLatest(id); ok {
		fname := strings.Split(fileName, "-")
		f := model.HistoryFiles{
			FileName: fname[0],
		}
		dao.UpdateLastFile(lastfileId, f)
	}
	f := model.Files{
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
	dao.UpdateFile(id, f)

	response.Success(c, nil, "配置文件修改成功")
}

func AllFiles(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	files := model.Files{}
	list, tx := files.AllFiles(query)

	total, err := CrudAll(query, tx, list)
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

func FileSearch(c *gin.Context) {
	var file model.SearchFile
	c.Bind(&file)
	search := file.Search

	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	list, tx := file.FileSearch(query, search)

	total, err := CrudAll(query, tx, list)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	JsonPagination(c, list, total, query)
}

func HistoryFiles(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	fileId := c.Query("id")
	FileId, err := strconv.Atoi(fileId)
	if err != nil {
		response.Fail(c, nil, "文件ID输入格式有误")
		return
	}

	files := model.HistoryFiles{}
	list, tx := files.HistoryFiles(query, FileId)

	total, err := CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	JsonPagination(c, list, total, query)
}

func LastFileRollBack(c *gin.Context) {
	var file model.RollBackFiles
	c.Bind(&file)

	lastfileId := file.HistoryFileID
	fileId := file.FileID
	user := file.UserUpdate
	userDept := file.UserDept

	lastfileText := dao.LastFileText(lastfileId)

	if ok, _, _ := dao.IsExistFileLatest(fileId); !ok {
		dao.SaveLatestFile(fileId)
	}

	fd := model.Files{
		UserUpdate: user,
		UserDept:   userDept,
		File:       lastfileText,
	}
	dao.UpdateFile(fileId, fd)
	response.JSON(c, http.StatusOK, http.StatusOK, nil, "已回退到历史版本")
}
