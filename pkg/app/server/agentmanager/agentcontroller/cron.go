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
 * Date: 2022-05-23 10:25:52
 * LastEditTime: 2022-05-23 15:16:10
 * Description: os scheduled task
 ******************************************************************************/
package agentcontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func CreatCron(c *gin.Context) {
	// 存入数据库
	var newCron model.CrontabUpdate
	c.Bind(&newCron)

	uuid := newCron.MachineUUID
	TaskName := newCron.TaskName
	description := newCron.Description
	spec := newCron.CronSpec
	command := newCron.Command
	status := newCron.Status

	if len(TaskName) == 0 {
		response.JSON(c, http.StatusOK, http.StatusBadRequest, nil, "任务名字不能为空")
		return
	}
	if dao.IsTaskNameExist(TaskName) {
		response.JSON(c, http.StatusOK, http.StatusBadRequest, nil, "任务名称已存在!")
		return
	}
	if len(spec) == 0 {
		response.JSON(c, http.StatusOK, http.StatusBadRequest, nil, "cron表达式不能为空")
		return
	}
	if len(command) == 0 {
		response.JSON(c, http.StatusOK, http.StatusBadRequest, nil, "执行命令不能为空")
		return
	}
	newcron := model.CrontabList{
		MachineUUID: uuid,
		TaskName:    TaskName,
		Description: description,
		CronSpec:    spec[:len(spec)-2],
		Command:     command,
		Status:      &status,
	}
	id := dao.NewCron(newcron)

	if !status {
		response.Fail(c, nil, "定时任务已保存")
		return
	}

	// 远程命令执行
	cronSpec, Command := dao.Id2CronInfo(id)
	cron_start, err := service.CronStart(uuid, id, cronSpec, Command)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}

	response.JSON(c, http.StatusOK, http.StatusOK, gin.H{"data": newCron, "cron": cron_start}, "任务已生效")
}

func DeleteCronTask(c *gin.Context) {
	var crons model.DelCrons
	var cronIds string
	c.Bind(&crons)
	uuid := crons.MachineUUID
	for _, cronId := range crons.IDs {
		_, err := service.StopAndDel(uuid, cronId)
		if err != nil {
			cronIds = strconv.Itoa(cronId) + ","
			continue
		}
		dao.DeleteTask(cronId)
	}
	if len(cronIds) != 0 {
		msg := fmt.Sprintf("以下任务编号未删除成功：%s", cronIds[:len(cronIds)-1])
		response.Fail(c, nil, msg)
		return
	}
	response.JSON(c, http.StatusOK, http.StatusOK, nil, "任务删除成功!")
}

func UpdateCron(c *gin.Context) {
	var Cron model.CrontabUpdate
	c.Bind(&Cron)
	id := Cron.ID
	TaskName := Cron.TaskName
	description := Cron.Description
	spec := Cron.CronSpec
	command := Cron.Command
	uuid := Cron.MachineUUID
	status := Cron.Status
	UpdateCron := model.CrontabList{
		TaskName:    TaskName,
		Description: description,
		CronSpec:    spec[:len(spec)-2],
		Command:     command,
		Status:      &status,
	}
	// 数据库内容修改
	dao.UpdateTask(id, UpdateCron)
	if !status {
		response.Fail(c, nil, "定时任务已保存,未执行")
		return
	}

	// 更新agent任务
	_, err := service.StopAndDel(uuid, id)
	if err != nil {
		msg := fmt.Sprintf("任务已保存,重启失败：%s", err)
		response.Fail(c, nil, msg)
		return
	}
	cron_start, err := service.CronStart(uuid, id, spec[:len(spec)-2], command)
	if err != nil {
		msg := fmt.Sprintf("任务已保存,重启失败：%s", err)
		response.Fail(c, nil, msg)
		return
	}

	response.JSON(c, http.StatusOK, http.StatusOK, gin.H{"cron": cron_start}, "任务更新成功,已开始执行")
}

func CronTaskStatus(c *gin.Context) {
	var Cron model.CrontabUpdate
	c.Bind(&Cron)
	id := Cron.ID
	uuid := Cron.MachineUUID
	status := Cron.Status

	if dao.IsTaskStatus(id, status) {
		response.Fail(c, nil, "请重新确认任务状态")
		return
	}

	dao.CronTaskStatus(id, status)

	if status {
		cron_stop, err := service.StopAndDel(uuid, id)
		if err != nil {
			response.Fail(c, nil, "任务暂停失败")
			return
		}
		response.JSON(c, http.StatusOK, http.StatusOK, gin.H{"cron": cron_stop}, "任务已暂停")
		return
	}

	cronSpec, Command := dao.Id2CronInfo(id)
	cron_start, err := service.CronStart(uuid, id, cronSpec, Command)
	if err != nil {
		response.Fail(c, gin.H{"error": err}, "任务执行失败!")
		return
	}

	response.JSON(c, http.StatusOK, http.StatusOK, gin.H{"cron": cron_start}, "任务已开启")
}

func CronTaskList(c *gin.Context) {
	uuid := c.Query("uuid")

	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.JSON(c, http.StatusOK, http.StatusBadRequest, gin.H{"status": false}, err.Error())
		return
	}

	cronlist := model.CrontabList{}
	list, tx := cronlist.CronList(uuid, query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	total, err := controller.CrudAll(query, tx, list)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}
	// 返回数据开始拼装分页的json
	controller.JsonPagination(c, list, total, query)
}
