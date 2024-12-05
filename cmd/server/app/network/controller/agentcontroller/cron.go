/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package agentcontroller

import (
	"fmt"
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/cron"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// CreatCron 创建定时任务
func CreatCron(c *gin.Context) {
	// 存入数据库
	var newCron cron.CrontabUpdate
	c.Bind(&newCron)

	uuid := newCron.MachineUUID
	TaskName := newCron.TaskName
	description := newCron.Description
	spec := newCron.CronSpec
	command := newCron.Command
	status := newCron.Status

	if len(TaskName) == 0 {
		response.Fail(c, nil, "任务名字不能为空")
		return
	}
	temp, err := cron.IsTaskNameExist(TaskName)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if temp {
		response.Fail(c, nil, "任务名称已存在!")
		return
	}
	if len(spec) == 0 {
		response.Fail(c, nil, "cron表达式不能为空")
		return
	}
	if len(command) == 0 {
		response.Fail(c, nil, "执行命令不能为空")
		return
	}
	newcron := cron.CrontabList{
		MachineUUID: uuid,
		TaskName:    TaskName,
		Description: description,
		CronSpec:    spec[:len(spec)-2],
		Command:     command,
		Status:      &status,
	}
	id, err := cron.NewCron(newcron)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if !status {
		response.Fail(c, nil, "定时任务已保存")
		return
	}

	// 远程命令执行
	cronSpec, Command, err := cron.Id2CronInfo(id)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}
	cron_start, err := cron.CronStart(uuid, id, cronSpec, Command)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}

	response.Success(c, gin.H{"data": newCron, "cron": cron_start}, "任务已生效")
}

// DeleteCronTask 是一个删除定时任务的函数
func DeleteCronTask(c *gin.Context) {
	var crons cron.DelCrons
	var cronIds string
	c.Bind(&crons)
	uuid := crons.MachineUUID
	for _, cronId := range crons.IDs {
		_, err := cron.StopAndDel(uuid, cronId)
		if err != nil {
			cronIds = strconv.Itoa(cronId) + ","
			continue
		}
		if err := cron.DeleteTask(cronId); err != nil {
			logger.Error(err.Error())
		}
	}
	if len(cronIds) != 0 {
		msg := fmt.Sprintf("以下任务编号未删除成功：%s", cronIds[:len(cronIds)-1])
		response.Fail(c, nil, msg)
		return
	}
	response.Success(c, nil, "任务删除成功!")
}

// UpdateCron 是一个处理更新定时任务的函数
func UpdateCron(c *gin.Context) {
	var Cron cron.CrontabUpdate
	c.Bind(&Cron)
	id := Cron.ID
	TaskName := Cron.TaskName
	description := Cron.Description
	spec := Cron.CronSpec
	command := Cron.Command
	uuid := Cron.MachineUUID
	status := Cron.Status
	UpdateCron := cron.CrontabList{
		TaskName:    TaskName,
		Description: description,
		CronSpec:    spec[:len(spec)-2],
		Command:     command,
		Status:      &status,
	}
	// 数据库内容修改
	if err := cron.UpdateTask(id, UpdateCron); err != nil {
		response.Fail(c, nil, err.Error())
	}
	if !status {
		response.Fail(c, nil, "定时任务已保存,未执行")
		return
	}

	// 更新agent任务
	_, err := cron.StopAndDel(uuid, id)
	if err != nil {
		msg := fmt.Sprintf("任务已保存,重启失败：%s", err)
		response.Fail(c, nil, msg)
		return
	}
	cron_start, err := cron.CronStart(uuid, id, spec[:len(spec)-2], command)
	if err != nil {
		msg := fmt.Sprintf("任务已保存,重启失败：%s", err)
		response.Fail(c, nil, msg)
		return
	}

	response.Success(c, gin.H{"cron": cron_start}, "任务更新成功,已开始执行")
}

// CronTaskStatus 处理任务状态更新请求
func CronTaskStatus(c *gin.Context) {
	var Cron cron.CrontabUpdate
	c.Bind(&Cron)
	id := Cron.ID
	uuid := Cron.MachineUUID
	status := Cron.Status
	temp, err := cron.IsTaskStatus(id, status)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if temp {
		response.Fail(c, nil, "请重新确认任务状态")
		return
	}

	if err := cron.CronTaskStatus(id, status); err != nil {
		response.Fail(c, nil, err.Error())
	}

	if status {
		cron_stop, err := cron.StopAndDel(uuid, id)
		if err != nil {
			response.Fail(c, nil, "任务暂停失败")
			return
		}
		response.Success(c, gin.H{"cron": cron_stop}, "任务已暂停")
		return
	}

	cronSpec, Command, err := cron.Id2CronInfo(id)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}
	cron_start, err := cron.CronStart(uuid, id, cronSpec, Command)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}

	response.Success(c, gin.H{"cron": cron_start}, "任务已开启")
}

// CronTaskList 获取定时任务列表
func CronTaskList(c *gin.Context) {
	uuid := c.Query("uuid")

	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	offset := query.PageSize * (query.Page - 1)
	total, data, err := cron.CronListPaged(uuid, offset, query.PageSize)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	response.DataPagination(c, data, int(total), query)
}
