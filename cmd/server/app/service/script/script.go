/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package script

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	batchservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/batch"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
)

type Script = dao.Script

type HistoryVersion = dao.HistoryVersion

type RunScriptMeta struct {
	BatchID      uint     `json:"batch_id"`
	MachineUUIDs []string `json:"machine_uuids"`
	ScriptID     uint     `json:"script_id"`
	Version      string   `json:"version"`
	Params       []string `json:"params"`
}

// 存储脚本文件
func AddScript(script *dao.Script) error {
	if len(script.Name) == 0 {
		return errors.New("请输入脚本文件名字")
	}
	if len(script.Content) == 0 {
		return errors.New("请输入脚本内容")
	}
	if len(script.Description) == 0 {
		return errors.New("请输入脚本描述")
	}

	if err := dao.AddScript(*script); err != nil {
		return errors.New("脚本文件添加失败")
	}
	return nil
}

func UpdateScript(script *dao.Script) error {
	if len(script.Name) == 0 {
		return errors.New("请输入脚本文件名字")
	}
	if len(script.Content) == 0 {
		return errors.New("请输入脚本内容")
	}
	if len(script.Description) == 0 {
		return errors.New("请输入脚本描述")
	}

	if err := dao.UpdateScript(*script); err != nil {
		return errors.New("脚本文件添加失败")
	}
	return nil
}

func DeleteScript(script_id uint, version string) error {
	if script_id < 1 && version == "" {
		return errors.New("script id or version abnormal")
	}

	if err := dao.DeleteScript(script_id, version); err != nil {
		return err
	}

	return nil
}

func ScriptList(query *response.PaginationQ) ([]*dao.Script, int, error) {
	if query.PageSize == 0 && query.Page == 0 {
		return nil, 0, errors.Errorf("pagesize: %d, page: %d", query.PageSize, query.Page)
	}

	scripts, total, err := dao.ScriptList(query)
	if err != nil {
		return nil, 0, err
	}
	return scripts, total, nil
}

func ScriptHistoryVersion(scriptid uint) ([]*dao.HistoryVersion, error) {
	if scriptid < 1 {
		return nil, errors.Errorf("script id abnormal, id: %d", scriptid)
	}

	scripts, err := dao.GetScriptHistoryVersion(scriptid)
	if err != nil {
		return nil, err
	}
	return scripts, nil
}

func RunScript(createName string, runscriptmeta *RunScriptMeta, batch *common.Batch) ([]batchservice.R, error) {
	var err error

	if runscriptmeta.ScriptID == 0 {
		return nil, errors.Errorf("script id abnormal, id: %d", runscriptmeta.ScriptID)
	}

	script_content := ""
	if runscriptmeta.Version == "" {
		script_content, err = dao.ShowScriptContent(runscriptmeta.ScriptID)
		if err != nil {
			return nil, err
		}
	} else {
		script_content, err = dao.ShowScriptWithVersion(runscriptmeta.ScriptID, runscriptmeta.Version)
		if err != nil {
			return nil, err
		}
	}

	var script_name string
	script, err := GetScriptByID(runscriptmeta.ScriptID)
	if err != nil {
		logger.Error("fail to get script by id: %s", err.Error())
		script_name = ""
	} else {
		script_name = script.Name
	}

	var batch_name []string
	batchName, _ := dao.GetBatchName(runscriptmeta.BatchID)
	batch_name = append(batch_name, batchName)
	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     fmt.Sprintf("脚本运行(%v)", batchName),
		Module:     auditlog.ScriptExec,
		User:       createName,
		Batches:    strings.Join(batch_name, ","),
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	cmds, err := GetDangerousCommandsInBlackList()
	if err != nil {
		return nil, errors.Errorf("run script error(dangerous commands list): %s", err.Error())
	}
	positions, matchedCommands := global.FindDangerousCommandsPos(script_content, cmds)
	if len(positions) > 0 {
		return nil, errors.New("Dangerous commands detected in script: " + strings.Join(matchedCommands, "\n"))
	}

	f := func(uuid string) batchservice.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
				LogId:        logId,
				ActionObject: "执行主机：" + agent.IP,
				UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
			})
			data, err := agent.RunScript(script_content, runscriptmeta.Params)
			if err != nil {
				logger.Error("run script error, agent:%s, command:%s", uuid, script_content)
			}
			logger.Debug("run script on agent result:%v", data)
			re := common.CmdResult{
				MachineUUID: uuid,
				MachineIP:   agent.IP,
				RetCode:     data.RetCode,
				Stdout:      data.Stdout,
				Stderr:      data.Stderr,
			}
			if len(data.Stderr) != 0 {
				auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, fmt.Sprintf("脚本 -> %s\n%s", script_name, data.Stderr))
			} else {
				auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, fmt.Sprintf("脚本 -> %s\n%s", script_name, data.Stdout))
			}
			return re
		}
		return common.CmdResult{}
	}

	result := batchservice.BatchProcess(batch, f, script_content, runscriptmeta.Params)

	if dao.GetSubLogStatus(logId) {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
	} else {
		auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	}

	return result, nil
}

func GetScriptByID(id uint) (*dao.Script, error) {
	script, err := dao.ShowScript(id)
	if err != nil {
		return nil, err
	}
	return script, nil
}
