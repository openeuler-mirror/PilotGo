/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Thu Aug 31 09:15:52 2023 +0800
 */
// common包定义agent与server之前传输的结构体数据
package common

import (
	ocommon "gitee.com/openeuler/PilotGo/pkg/utils/os/common"
)

type AgentOverview struct {
	IP          string `mapstructure:"IP"`
	SysInfo     *ocommon.SystemInfo
	DiskUsage   []ocommon.DiskUsageINfo
	MemoryInfo  *ocommon.MemoryConfig
	CpuInfo     *ocommon.CPUInfo
	IsImmutable bool
}
