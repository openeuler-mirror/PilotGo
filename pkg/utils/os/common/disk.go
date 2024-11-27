/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

type DiskIOInfo struct {
	PartitionName string
	Label         string
	ReadCount     uint64
	WriteCount    uint64
	ReadBytes     uint64
	WriteBytes    uint64
	IOTime        uint64
}

type DiskUsageINfo struct {
	Device      string `json:"device"`
	Path        string `json:"path"`
	Fstype      string `json:"fstype"`
	Total       string `json:"total"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedPercent"`
}
