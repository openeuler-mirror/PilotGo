// common包定义agent与server之前传输的结构体数据
package common

import (
	ocommon "gitee.com/openeuler/PilotGo/utils/os/common"
)

type AgentOverview struct {
	IP          string `mapstructure:"IP"`
	SysInfo     *ocommon.SystemInfo
	DiskUsage   []ocommon.DiskUsageINfo
	MemoryInfo  *ocommon.MemoryConfig
	CpuInfo     *ocommon.CPUInfo
	IsImmutable bool
}
