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
 * Description: agent config file handler
 ******************************************************************************/

package agentcontroller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func ReadFile(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	filepath := c.Query("file")
	result, Err, err := agent.ReadFile(filepath)
	if err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, gin.H{"file": result}, "Success")
}

func GetAgentFiles(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	var datas []interface{}
	// 获取repo文件
	repos, Err, err := agent.GetRepoFile()
	if err != nil {
		response.Fail(c, nil, Err)
		return
	}
	repo := controller.InterfaceToSlice(repos)
	datas = append(datas, repo...)

	// 获取网络配置文件
	network, Err, err := agent.GetNetWorkFile()
	if err != nil {
		response.Fail(c, nil, Err)
		return
	}
	datas = append(datas, network)

	data, err := controller.DataPaging(query, datas, len(datas))
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	controller.JsonPagination(c, data, len(datas), query)
}

func SaveFileToDatabase(c *gin.Context) {
	var file model.Files
	c.Bind(&file)

	path := file.SourcePath
	filename := file.FileName
	if len(filename) == 0 {
		response.Fail(c, nil, "请输入配置文件名字")
		return
	}
	description := file.Description
	if len(description) == 0 {
		response.Fail(c, nil, "请添加文件描述")
		return
	}
	text := file.File
	if len(text) == 0 {
		response.Fail(c, nil, "请重新检查文件内容")
		return
	}

	fd := model.Files{
		SourcePath:  path,
		FileName:    filename,
		Description: description,
		File:        text,
	}
	dao.SaveFile(fd)
	response.Success(c, nil, "文件保存成功")
}

func UpdateAgentFile(c *gin.Context) {
	var file model.HistoryFiles
	c.Bind(&file)

	ip := file.IP
	uuid := file.UUID
	path := file.Path
	ipdept := file.IPDept
	if len(path) == 0 {
		response.Fail(c, nil, "请检查配置文件路径")
		return
	}
	filename := file.FileName
	if len(filename) == 0 {
		response.Fail(c, nil, "请检查配置文件名字")
		return
	}
	text := file.File
	if len(text) == 0 {
		response.Fail(c, nil, "请重新检查文件内容")
		return
	}

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "server端获取uuid失败!")
		return
	}

	result, Err, err := agent.UpdateFile(path, filename, text)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}

	// 保存历史版本
	time := service.NowTime()
	fd := model.HistoryFiles{
		IP:       ip,
		IPDept:   ipdept,
		UUID:     uuid,
		Path:     path,
		FileName: filename + "-" + time,
		File:     result.(string),
	}
	dao.SaveHistoryFile(fd)

	response.JSON(c, http.StatusOK, http.StatusOK, result, "配置文件已更新")
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
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, err := controller.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	controller.JsonPagination(c, list, total, query)
}

func AllHistoryFiles(c *gin.Context) {
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	files := model.HistoryFiles{}
	list, tx := files.AllHistoryFiles(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, err := controller.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	controller.JsonPagination(c, list, total, query)
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

	total, err := controller.CrudAll(query, tx, list)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	controller.JsonPagination(c, list, total, query)
}

func LastFileSearch(c *gin.Context) {
	var file model.SearchFile
	c.Bind(&file)
	search := file.Search

	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	list, tx := file.LastFileSearch(query, search)

	total, err := controller.CrudAll(query, tx, list)
	if err != nil {
		response.Response(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}
	controller.JsonPagination(c, list, total, query)
}

func UpdateFile(c *gin.Context) {
	var file model.Files
	c.Bind(&file)

	path := file.SourcePath
	filename := file.FileName
	description := file.Description
	text := file.File
	if !dao.IsExistId(file.ID) {
		response.Fail(c, nil, "id有误,请重新确认该文件是否存在")
		return
	}
	dao.UpdateFile(file.ID, path, filename, description, text)

	response.Success(c, nil, "配置文件修改成功")
}

func DeleteFile(c *gin.Context) {
	var files model.DeleteFiles
	c.Bind(&files)

	for _, fileId := range files.FileIDs {
		dao.DeleteFile(fileId)
	}
	response.Success(c, nil, "储存的文件已从数据库中删除")
}

func FindLastVersionFile(c *gin.Context) {
	uuid := c.Query("uuid")
	filename := c.Query("name")
	lastfiles := dao.FindLastVersionFile(uuid, filename)
	response.Success(c, gin.H{"oldfiles": lastfiles}, "获取该文件的历史版本")
}

func LastFileRollBack(c *gin.Context) {
	var file model.HistoryFiles
	c.Bind(&file)

	ip := file.IP
	uuid := file.UUID
	path := file.Path
	ipdept := file.IPDept
	if len(path) == 0 {
		response.Fail(c, nil, "请检查配置文件路径")
		return
	}
	text := file.File
	if len(text) == 0 {
		response.Fail(c, nil, "请重新检查文件内容")
		return
	}

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "server端获取uuid失败!")
		return
	}

	filename := file.FileName
	if len(filename) == 0 {
		response.Fail(c, nil, "请检查配置文件名字")
		return
	}
	name := strings.Split(filename, "-")[0]
	if ok := dao.IsFileLatest(name, uuid); !ok {
		LatestVersion, Err, err := agent.UpdateFile(path, name, text)
		if len(Err) != 0 || err != nil {
			response.Fail(c, nil, Err)
			return
		}
		//  保存最新版本
		fd := model.HistoryFiles{
			IP:       ip,
			IPDept:   ipdept,
			UUID:     uuid,
			Path:     path,
			FileName: filename + "-latest",
			File:     LatestVersion.(string),
		}
		dao.SaveHistoryFile(fd)
		response.JSON(c, http.StatusOK, http.StatusOK, LatestVersion, "已回退到历史版本")
		return
	}

	_, Err, err := agent.UpdateFile(path, name, text)
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.JSON(c, http.StatusOK, http.StatusOK, nil, "已回退到历史版本")
}

func FileBroadcastToAgents(c *gin.Context) {
	var fb model.FileBroadcast
	c.Bind(&fb)

	path := fb.Path
	filename := fb.FileName
	text := fb.Text

	if len(path) == 0 {
		response.Fail(c, nil, "路径为空，请检查配置文件路径")
		return
	}
	if len(filename) == 0 {
		response.Fail(c, nil, "文件名为空，请检查配置文件名字")
		return
	}
	if len(text) == 0 {
		response.Fail(c, nil, "文件内容为空，请重新检查文件内容")
		return
	}

	logParent := model.AgentLogParent{
		UserName:   fb.User,
		DepartName: fb.UserDept,
		Type:       model.LogTypeBroadcast,
	}
	logParentId := dao.ParentAgentLog(logParent)

	StatusCodes := make([]string, 0)

	for _, uuid := range fb.UUID {
		agent := agentmanager.GetAgent(uuid)
		if agent == nil {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: filename,
				Action:          model.BroadcastFile,
				StatusCode:      http.StatusBadRequest,
				Message:         "获取uuid失败",
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		}

		result, Err, err := agent.UpdateFile(path, filename, text)
		if len(Err) != 0 || err != nil {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: filename,
				Action:          model.BroadcastFile,
				StatusCode:      http.StatusBadRequest,
				Message:         Err,
			}
			dao.AgentLog(log)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusBadRequest))
			continue
		} else {
			log := model.AgentLog{
				LogParentID:     logParentId,
				IP:              dao.UUID2MacIP(uuid),
				OperationObject: filename,
				Action:          model.BroadcastFile,
				StatusCode:      http.StatusOK,
				Message:         "配置文件下发成功",
			}
			dao.AgentLog(log)

			// 保存历史版本
			time := service.NowTime()
			ip, _, ipdept := dao.MachineBasic(uuid)
			fd := model.HistoryFiles{
				IP:       ip,
				IPDept:   ipdept,
				UUID:     uuid,
				Path:     path,
				FileName: filename + "-" + time,
				File:     result.(string),
			}
			dao.SaveHistoryFile(fd)

			StatusCodes = append(StatusCodes, strconv.Itoa(http.StatusOK))
		}
	}
	status := service.BatchActionStatus(StatusCodes)
	dao.UpdateParentAgentLog(logParentId, status)
	response.Success(c, nil, "配置文件下发完成!")
}
