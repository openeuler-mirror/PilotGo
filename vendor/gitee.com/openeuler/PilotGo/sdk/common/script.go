/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Oct 11 16:18:11 2023 +0800
 */
package common

type AsyncCmdResult struct {
	TaskID string       `json:"task_id"`
	Result []*CmdResult `json:"result"`
}

type CmdResult struct {
	MachineUUID string `json:"machine_uuid"`
	MachineIP   string `json:"machine_ip"`
	RetCode     int    `json:"retcode"`
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
}

type CmdStruct struct {
	Batch   *Batch `json:"batch"`
	Command string `json:"command"`
}
