/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-05-18 09:08:08
 * LastEditTime: 2023-09-15 14:58:17
 * Description: 批次管理业务逻辑
 ******************************************************************************/
package batch

import (
	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/depart"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	scommon "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
)

type Batch = dao.Batch
type BatchMachines = dao.BatchMachines
type CreateBatchParam struct {
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Manager     string   `json:"Manager"`
	DepartName  []string `json:"DepartName"`
	DepartID    []int    `json:"DepartID"`
	Machines    []int    `json:"Machines"`
}

func CreateBatch(batchinfo *CreateBatchParam) error {
	ExistNameBool, err := dao.IsExistName(batchinfo.Name)
	if err != nil {
		return err
	}
	if ExistNameBool {
		return errors.New("已存在该名称批次")
	}

	if len(batchinfo.Machines) == 0 && len(batchinfo.DepartID) == 0 {
		return errors.New("请先选择机器IP或部门")
	}
	// 机器id列表
	Departids := make([]int, 0)

	if len(batchinfo.Machines) == 0 && len(batchinfo.DepartID) != 0 {
		// 点选部门创建批次
		for _, ids := range batchinfo.DepartID {
			Departids = append(Departids, ids)
			depart.ReturnSpecifiedDepart(ids, &Departids)
		}

		// 机器所属部门ids
		var departIdlist string
		for _, id := range Departids {
			departIdlist = departIdlist + "," + strconv.Itoa(id)
			departIdlist = strings.Trim(departIdlist, ",")
		}

		// 机器所属部门
		var departNamelist string
		list := dao.DepartIdsToGetDepartNames(Departids)
		departNamelist = strings.Join(list, ",")

		Batch := &dao.Batch{
			Name:        batchinfo.Name,
			Description: batchinfo.Description,
			Manager:     batchinfo.Manager,
			Depart:      departIdlist,
			DepartName:  departNamelist,
		}
		err = dao.CreateBatchMessage(Batch)
		if err != nil {
			return err
		}
	} else if len(batchinfo.DepartID) == 0 && len(batchinfo.Machines) != 0 {
		//按机器创建批次
		Batch := &dao.Batch{
			Name:        batchinfo.Name,
			Description: batchinfo.Description,
			Manager:     batchinfo.Manager,
		}
		err = dao.CreateBatchMessage(Batch)
		if err != nil {
			return err
		}

		// 添加数据到Batch2Machine
		for _, id := range batchinfo.Machines {
			err := dao.AddBatch2Machine(dao.BatchMachines{
				BatchID:       Batch.ID,
				MachineNodeID: uint(id),
			})
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("输入参数有误")
	}
	return nil
}

// 查询所有批次
func SelectBatch() ([]dao.Batch, error) {
	return dao.GetBatch()
}

// 分页查询所有批次
func GetBatchPaged(offset, size int) (int64, []Batch, error) {
	return dao.GetBatchrPaged(offset, size)
}

// 删除批次
func DeleteBatch(ids []int) error {
	for _, value := range ids {
		err := dao.DeleteBatch(value)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}

// 更改批次
func UpdateBatch(batchid int, name, description string) error {
	temp, err := dao.IsExistID(batchid)
	if err != nil {
		return err
	}
	if !temp {
		return errors.New("不存在该批次")
	}

	err = dao.UpdateBatch(batchid, name, description)
	return err
}

// 分页获取某批次的机器信息
func GetBatchMachines(offset, size, batchid int) (int64, []dao.MachineNode, error) {
	//检查批次id是否存在
	isExist, err := dao.IsExistID(batchid)
	if !isExist {
		return 0, nil, errors.New("批次不存在")
	}
	if err != nil {
		return 0, nil, err
	}

	count, machineIdlist, err := dao.GetMachineIDPaged(offset, size, batchid)
	if err != nil {
		return 0, nil, err
	}
	// 获取机器的所有信息
	machinesInfo := make([]dao.MachineNode, 0)
	for _, macId := range machineIdlist {
		MacInfo, err := dao.MachineInfo(int(macId.MachineNodeID))
		if err != nil {
			logger.Error(err.Error())
		}
		machinesInfo = append(machinesInfo, MacInfo)
	}

	return count, machinesInfo, nil
}

// 根据批次id获取所属的所有uuids
func BatchIds2UUIDs(batchIds []int) (uuids []string) {
	for _, batchid := range batchIds {
		uuids = append(uuids, GetMachineUUIDS(batchid)...)
	}
	return uuids
}

// GetMachineUUIDS from batch get all machines
func GetMachineUUIDS(batchid int) []string {
	if batchid == 0 {
		return nil
	}

	machineIdlist, err := dao.GetMachineID(batchid)
	if err != nil {
		logger.Error("Failed to get machine IDs for batch ID %d: %v", batchid, err)
		return nil
	}

	uuids := make([]string, 0, len(machineIdlist)) // 预先分配足够空间
	for _, macID := range machineIdlist {
		macuuid, err := dao.MachineIdToUUID(int(macID))
		if err != nil {
			logger.Error("Failed to get UUID for machine ID %d in batch ID %d: %v", macID, batchid, err)
			continue // 跳过当前循环，不将错误的UUID添加到列表中
		}
		uuids = append(uuids, macuuid)
	}
	return uuids
}

// get common.CmdStruct batch uuids
func GetBatchMachineUUIDS(b *scommon.Batch) []string {
	var machine_uuids []string
	if b.MachineUUIDs != nil {
		machine_uuids = append(machine_uuids, b.MachineUUIDs...)
	}
	if b.BatchIds != nil {
		for _, v := range b.BatchIds {
			machine_uuids = append(machine_uuids, GetMachineUUIDS(v)...)
		}
	}
	if b.DepartmentIDs != nil {
		for _, v := range b.DepartmentIDs {
			r, err := depart.MachineList(v)
			if err != nil {
				logger.Error("failed to get machine uuid from departid:%s", err.Error())
			}
			for _, k := range r {
				machine_uuids = append(machine_uuids, k.UUID)
			}
		}
	}
	//给uuid除重
	machine_uuids = utils.RemoveRepByMap(machine_uuids)
	return machine_uuids
}

type R interface{}

func BatchProcess(b *scommon.Batch, f func(uuid string) R, it ...interface{}) []R {
	uuids := GetBatchMachineUUIDS(b)
	result := []R{}
	for _, uuid := range uuids {
		r := f(uuid)
		result = append(result, r)
	}
	return result
}
