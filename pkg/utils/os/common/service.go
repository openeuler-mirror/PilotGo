package common

type ListService struct {
	Name   string
	LOAD   string
	Active string
	SUB    string
}

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
	UNKnown       = "unknown"
)

type ServiceInfo struct {
	ServiceName         string
	UnitName            string
	UnitType            string
	ServicePath         string //配置文件放置的目录
	ServiceExectStart   string //实际执行daemon的指令或脚本程序
	ServiceTime         string //开启时间
	ServiceActiveStatus string
	ServiceLoadedStatus string
}
