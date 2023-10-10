package common

import "encoding/json"

// event消息类型定义
const (
	// 主机安装软件包
	MsgPackageInstall = 0
	// 主机升级软件包
	MsgPackageUpdate = 1
	// 主机卸载软件包
	MsgPackageUninstall = 2
	// 主机ip变更
	MsgIPChange = 3

	// 平台新增主机
	MsgHostAdd = 10
	// 平台移除主机
	MsgHostRemove = 11

	// 插件添加
	MsgPluginAdd = 20
	// 插件卸载
	MsgPluginRemove = 21
)

type EventMessage struct {
	MessageType int
	MessageData string
}

// 将 MessageData json字符串转换成指定结构体的message消息数据
func ToMessage(d string, s interface{}) error {
	return json.Unmarshal([]byte(d), s)
}

type MDPackageInstall struct {
	HostUUID string
	Name     string
	Version  string
	Time     string
}

type MDPackageUpdate struct {
	HostUUID string
	Name     string
	Version  string
	Time     string
}

type MDPackageUninstall struct {
	HostUUID string
	Name     string
	Version  string
	Time     string
}

type MDIPChange struct {
	HostUUID string
	NewIP    string
}
