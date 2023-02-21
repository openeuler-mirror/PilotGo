package common

type SystemInfo struct {
	IP              string
	Platform        string //系统平台
	PlatformVersion string //系统版本
	KernelVersion   string //内核版本
	KernelArch      string //内核支持架构
	HostId          string //系统id
	Uptime          string //系统最新启动时间
}
