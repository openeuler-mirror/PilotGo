/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
package common

const (
	ServiceActiveStatusRunning  = "running"
	ServiceActiveStatusExited   = "exited"
	ServiceActiveStatusWaiting  = "waiting"
	ServiceActiveStatusInactive = "inactive"
	ServiceActiveStatusUnknown  = "unknown"

	ServiceLoadedStatusEnabled  = "enabled"
	ServiceLoadedStatusDisabled = "disabled"
	ServiceLoadedStatusStatic   = "static"
	ServiceLoadedStatusMask     = "mask"
	ServiceLoadedStatusUnknown  = "unknown"
)

const (
	ServiceUnit   = "service"
	SocketUnit    = "socket"
	TargetUnit    = "target"
	MountUnit     = "mount"
	AutomountUnit = "automount"
	PathUnit      = "path"
	TimeUnit      = "time"
)

type ServiceResult struct {
	MachineUUID   string
	MachineIP     string
	ServiceSample ServiceInfo
}

type ServiceInfo struct {
	ServiceName            string
	UnitName               string
	UnitType               string
	ServicePath            string //配置文件放置的目录
	ServiceAfter           string //在什么服务启动后启动
	ServiceBefore          string //在什么服务启动前启动
	ServiceRequires        string //需要的daemon
	ServiceWants           string //与requires相反
	ServiceEnvironmentFile string //启动脚本的环境配置文件
	ServiceExectStart      string //实际执行daemon的指令或脚本程序
	ServiceActiveStatus    string
	ServiceLoadedStatus    string
	StartTime              string
}

type ServiceStruct struct {
	Batch       *Batch `json:"batch"`
	ServiceName string `json:"service"`
}

type Result struct {
	Code    int              `json:"code"`
	Mseeage string           `json:"msg"`
	Data    []*ServiceResult `json:"data"`
}
