// 提供公共数据结构定义
package common

const (
	// 插件正在运行
	StatusRunning = "running"
	// 插件已加载，但未运行
	StatusLoaded = "loaded"
	// 插件离线，无法访问
	StatusOffline = "offline"
)
