/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	scriptservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/script"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 存储脚本文件
func AddScriptHandler(c *gin.Context) {
	script := &scriptservice.Script{}
	if err := c.ShouldBindJSON(script); err != nil {
		logger.Error("fail to create script(bind): %s", err.Error())
		response.Fail(c, nil, fmt.Sprintf("脚本文件添加失败: %s", err.Error()))
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	cmds, err := scriptservice.GetDangerousCommandsInBlackList()
	if err != nil {
		logger.Error("创建高危命令脚本失败: %v", err.Error())
		response.Fail(c, nil, "创建失败")
		return
	}
	positions, matchedCommands := global.FindDangerousCommandsPos(script.Content, cmds)
	if len(positions) > 0 {
		logger.Error("脚本中检测到高危命令: %v", strings.Join(matchedCommands, ","))
		response.Fail(c, nil, "脚本中检测到高危命令: "+strings.Join(matchedCommands, "\n"))
		return
	}

	if err := scriptservice.AddScript(script); err != nil {
		logger.Error("创建脚本失败: %v", err.Error())
		response.Fail(c, nil, fmt.Sprintf("脚本文件添加失败: %s", err.Error()))
		return
	}

	global.SendRemindMsg(
		global.ServerSendMsg,
		fmt.Sprintf("用户 %s 创建脚本 %s", u.Username, script.Name),
	)

	response.Success(c, nil, "成功")
}

func UpdateScriptHandler(c *gin.Context) {
	script := &scriptservice.Script{}
	if err := c.ShouldBindJSON(script); err != nil {
		logger.Error("fail to edit script(bind): %s", err.Error())
		response.Fail(c, nil, fmt.Sprintf("脚本文件添加失败: %s", err.Error()))
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "编辑脚本",
		Module:     auditlog.ScriptEdit,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "编辑脚本：" + script.Name,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	cmds, err := scriptservice.GetDangerousCommandsInBlackList()
	if err != nil {
		logger.Error("更新脚本检测到高危命令: %v", err.Error())
		response.Fail(c, nil, "更新失败")
		return
	}
	positions, matchedCommands := global.FindDangerousCommandsPos(script.Content, cmds)
	if len(positions) > 0 {
		logger.Error("脚本中检测到高危命令: %v", strings.Join(matchedCommands, ","))
		response.Fail(c, nil, "脚本中检测到高危命令: "+strings.Join(matchedCommands, "\n"))
		return
	}

	if err := scriptservice.UpdateScript(script); err != nil {
		logger.Error("更新脚本失败: %v", err.Error())
		response.Fail(c, nil, fmt.Sprintf("脚本文件更新失败: %s", err.Error()))
		return
	}

	global.SendRemindMsg(
		global.ServerSendMsg,
		fmt.Sprintf("用户 %s 更新脚本 %s", u.Username, script.Name),
	)
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")

	response.Success(c, nil, "成功")
}

func DeleteScriptHandler(c *gin.Context) {
	req_body := struct {
		ScriptID uint   `json:"script_id"`
		Version  string `json:"version"`
	}{}
	if err := c.ShouldBindJSON(&req_body); err != nil {
		logger.Error("fail to delete script(bind): %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	var script_name string
	_script, err := scriptservice.GetScriptByID(req_body.ScriptID)
	if err != nil {
		logger.Error("fail to get script by id: %s", err.Error())
		script_name = ""
	} else {
		script_name = _script.Name
	}

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "删除脚本",
		Module:     auditlog.ScriptDelete,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})
	subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
		LogId:        logId,
		ActionObject: "删除脚本：" + script_name,
		UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	})

	if err := scriptservice.DeleteScript(req_body.ScriptID, req_body.Version); err != nil {
		logger.Error("fail to delete script: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	global.SendRemindMsg(
		global.MachineSendMsg,
		fmt.Sprintf("用户 %s 删除脚本 %s %s", u.Username, script_name, req_body.Version),
	)
	auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, "操作成功")

	response.Success(c, nil, "成功")
}

func RunScriptHandler(c *gin.Context) {
	body := &scriptservice.RunScriptMeta{}
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Error("fail to run script(bind): %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	var script_name string
	script, err := scriptservice.GetScriptByID(body.ScriptID)
	if err != nil {
		logger.Error("fail to get script by id: %s", err.Error())
		script_name = ""
	} else {
		script_name = script.Name
	}

	batch := &common.Batch{}
	if body.BatchID < 1 && len(body.MachineUUIDs) == 0 {
		logger.Error("fail to run script, batchid and machine_uuids are both empty")
		response.Fail(c, nil, "目标类型错误")
		return
	}
	if body.BatchID >= 1 {
		batch.BatchIds = append(batch.BatchIds, int(body.BatchID))
	} else {
		batch.MachineUUIDs = append(batch.MachineUUIDs, body.MachineUUIDs...)
	}

	result, err := scriptservice.RunScript(u.Username, body, batch)
	if err != nil {
		logger.Error("fail to run script: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	result_str := []string{}
	for _, v := range result {
		cmdR, ok := v.(common.CmdResult)
		if !ok {
			continue
		}
		result_str = append(result_str, fmt.Sprintf("\n机器IP:%s\nstdout:\n%s\nstderr:\n%s", cmdR.MachineIP, cmdR.Stdout, cmdR.Stderr))
	}

	global.SendRemindMsg(
		global.MachineSendMsg,
		fmt.Sprintf("用户 %s 执行脚本 %s %s, batch: %v, machines: %v", u.Username, script_name, body.Version, body.BatchID, body.MachineUUIDs),
	)

	response.Success(c, result, "成功")
}

func GetScriptListHandler(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		logger.Error("fail to get script list(bind): %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	scripts, total, err := scriptservice.ScriptList(query)
	if err != nil {
		logger.Error("fail to get script list: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.DataPagination(c, scripts, total, query)
}

func GetScriptHistoryVersionHandler(c *gin.Context) {
	scriptid := c.Query("script_id")

	id, err := strconv.Atoi(scriptid)
	if err != nil {
		logger.Error("fail to get script history version: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	scripts, err := scriptservice.ScriptHistoryVersion(uint(id))
	if err != nil {
		logger.Error("fail to get script history version: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, scripts, "成功")
}

func UpdateCommandsBlackListHandler(c *gin.Context) {
	body := &struct {
		WhiteList []uint `json:"white_list"`
	}{}
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Error("fail to update script blacklist(bind): %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}

	if err := scriptservice.UpdateCommandsBlackList(body.WhiteList); err != nil {
		logger.Error("fail to update script blacklist: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")
}

func GetDangerousCommandsList(c *gin.Context) {
	commands, err := scriptservice.GetDangerousCommandsList()
	if err != nil {
		logger.Error("fail to get dangerous commands list: %s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, commands, "成功")
}
