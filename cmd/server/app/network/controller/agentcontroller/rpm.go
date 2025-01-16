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
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/agentmanager"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type RPMS struct {
	UUIDs        []string `json:"uuid"`
	RPM          string   `json:"rpm"`
	UserName     string   `json:"userName"`
	UserDeptName string   `json:"userDept"`
}

func AllRpmHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_all, err := agent.AllRpm()
	if err != nil {
		response.Fail(c, nil, "获取已安装rpm包列表失败!")
		return
	}
	response.Success(c, gin.H{"rpm_all": rpm_all}, "Success")
}
func RpmSourceHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_source, err := agent.RpmSource(rpmname)
	if err != nil {
		response.Fail(c, nil, "获取源软件包名以及源失败!")
		return
	}
	response.Success(c, gin.H{"rpm_source": rpm_source}, "Success")
}
func RpmInfoHandler(c *gin.Context) {
	uuid := c.Query("uuid")
	rpmname := c.Query("rpm")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	rpm_info, Err, err := agent.RpmInfo(rpmname)
	if len(Err) != 0 || err != nil {
		response.Fail(c, gin.H{"error": Err}, "获取源软件包信息失败!")
		return
	} else {
		response.Success(c, gin.H{"rpm_info": rpm_info}, "Success")
	}

}

func InstallRpmHandler(c *gin.Context) {
	var rpm RPMS
	if err := c.Bind(&rpm); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}
	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	if !utils.CheckString(rpm.RPM) {
		response.Fail(c, nil, "软件包名有除_+-.以外的特殊字符")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "软件包安装-" + rpm.RPM,
		Module:     auditlog.RPMInstall,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	var codeMap = make(map[string]struct{})
	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)

		if agent == nil {
			continue
		}
		subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
			LogId:        logId,
			ActionObject: "软件包安装：" + agent.IP,
			UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
		})

		info, err := agent.AgentOverview()
		if err != nil {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "获取agent主机基础信息失败")
			codeMap[uuid] = struct{}{}
			continue
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "Install rpm is not supported on NestOS For Container")
			codeMap[uuid] = struct{}{}
			continue
		}

		_, Err, err := agent.InstallRpm(rpm.RPM)
		if err != nil || len(Err) != 0 {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, Err)
			codeMap[uuid] = struct{}{}
			continue
		} else {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, rpm.RPM+"安装成功")
		}
	}

	global.SendRemindMsg(
		global.MachineSendMsg,
		fmt.Sprintf("用户 %s 执行 %s 软件包安装, machines: %v", u.Username, rpm.RPM, rpm.UUIDs),
	)

	if len(codeMap) == 0 {
		auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	} else {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
	}
	response.Success(c, nil, "软件包安装指令下发成功")
}

func RemoveRpmHandler(c *gin.Context) {
	var rpm RPMS
	if err := c.Bind(&rpm); err != nil {
		response.Fail(c, nil, "parameter error")
		return
	}

	if len(rpm.UUIDs) == 0 {
		response.Fail(c, nil, "机器uuid不能为空")
		return
	}

	if !utils.CheckString(rpm.RPM) {
		response.Fail(c, nil, "软件包名有除_+-.以外的特殊字符")
		return
	}

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	logId, _ := auditlog.Add(&auditlog.AuditLog{
		Action:     "软件包卸载-" + rpm.RPM,
		Module:     auditlog.RPMRemove,
		User:       u.Username,
		Batches:    "",
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	})

	var codeMap = make(map[string]struct{})
	for _, uuid := range rpm.UUIDs {
		agent := agentmanager.GetAgent(uuid)

		if agent == nil {
			continue
		}
		subLogId, _ := auditlog.AddSubLog(&auditlog.SubLog{
			LogId:        logId,
			ActionObject: "软件包卸载：" + agent.IP,
			UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
		})

		info, err := agent.AgentOverview()
		if err != nil {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "获取agent主机基础信息失败")
			codeMap[uuid] = struct{}{}
			continue
		}
		if info.SysInfo.Platform == "NestOS For Container" {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, "Install rpm is not supported on NestOS For Container")
			codeMap[uuid] = struct{}{}
			continue
		}

		_, Err, err := agent.RemoveRpm(rpm.RPM)
		if len(Err) != 0 || err != nil {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusFail, Err)
			codeMap[uuid] = struct{}{}
			continue
		} else {
			auditlog.UpdateSubLog(subLogId, auditlog.StatusSuccess, rpm.RPM+"卸载成功")
		}
	}

	if len(codeMap) == 0 {
		auditlog.UpdateLog(logId, auditlog.StatusSuccess)
	} else {
		auditlog.UpdateLog(logId, auditlog.StatusFail)
	}
	response.Success(c, nil, "软件包卸载指令下发成功")
}
