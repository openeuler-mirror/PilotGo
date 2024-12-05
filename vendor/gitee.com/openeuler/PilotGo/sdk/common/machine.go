/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
package common

type MachineNode struct {
	UUID        string `json:"uuid"`
	Department  string `json:"department"`
	IP          string `json:"ip"`
	CPUArch     string `json:"cpu_arch"`
	OS          string `json:"os"`
	RunStatus   string `json:"runstatus"`
	MaintStatus string `json:"maintatatus"`
}

type BatchList struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Manager     string `json:"manager"`
}

type Batch struct {
	BatchIds      []int    `json:"batch_ids"`
	DepartmentIDs []int    `json:"department_ids"`
	MachineUUIDs  []string `json:"machine_uuids"`
}

type File struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type NodeResult struct {
	UUID  string
	Error string
	Data  interface{}
}
