/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugins licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 24 10:02:04 2024 +0800
 */
package sdk

import "time"

// event消息类型定义/msgTypeId
const (
	// 用户登录
	MsgUserLogin = 0
	// 用户退出
	MsgUserLogout = 1

	// 平台新增主机
	MsgHostAdd = 10
	// 平台移除主机
	MsgHostRemove = 11
	// 主机上线
	MsgHostOnline = 12
	// 平台离线
	MsgHostOffline = 13

	// 主机安装软件包
	MsgPackageInstall = 20
	// 主机升级软件包
	MsgPackageUpdate = 21
	// 主机卸载软件包
	MsgPackageUninstall = 22
	// 主机ip变更
	MsgIPChange = 30

	// 插件添加
	MsgPluginAdd = 40
	// 插件卸载
	MsgPluginRemove = 41
	// 插件上线
	MsgPluginOnline = 42
	// 插件离线
	MsgPluginOffline = 43
)

func GetMessageTypeString(msgType int) string {
	switch msgType {
	case MsgUserLogin:
		return "用户登录"
	case MsgUserLogout:
		return "用户退出"
	case MsgHostAdd:
		return "平台新增主机"
	case MsgHostRemove:
		return "平台移除主机"
	case MsgHostOnline:
		return "主机上线"
	case MsgHostOffline:
		return "主机离线"
	case MsgPackageInstall:
		return "主机安装软件包"
	case MsgPackageUpdate:
		return "主机升级软件包"
	case MsgPackageUninstall:
		return "主机卸载软件包"
	case MsgIPChange:
		return "主机ip变更"
	case MsgPluginAdd:
		return "插件添加"
	case MsgPluginRemove:
		return "插件卸载"
	case MsgPluginOnline:
		return "插件上线"
	case MsgPluginOffline:
		return "插件离线"
	default:
		return "未知的消息类型"
	}
}

var MessageTypes = []string{"用户登录", "用户退出", "平台新增主机", "平台移除主机", "主机上线", "主机离线", "主机安装软件包", "主机升级软件包", "主机卸载软件包", "主机ip变更", "插件添加", "插件卸载", "插件上线", "插件离线"}

type MessageData struct {
	MsgType     int         `json:"msg_type_id"`
	MessageType string      `json:"msg_type"`
	TimeStamp   time.Time   `json:"timestamp"`
	Data        interface{} `json:"data"`
}

type MDUserSystemSession struct { //平台登录、退出
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type MDHostChange struct { // 主机新增、移除、上线、离线
	IP         string `json:"ip"`
	OS         string `json:"os"`
	OSVersion  string `json:"os_version"`
	PrettyName string `json:"pretty_name"`
	CPU        string `json:"cpu"`
	Status     string `json:"status"` //在线状态
}

type MDHostPackageOpt struct { //软件包安装、升级、卸载
	HostIP  string `json:"hostIp"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type MDHostIPChange struct { //主机ip变更
	HostUUID string `json:"host_uuid"`
	NewIP    string `json:"new_ip"`
	OldIP    string `json:"old_ip"`
}

type MDPluginChange struct { // 插件新增、移除、上线、离线
	PluginName  string `json:"plugin_name"`
	Version     string `json:"version"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}
