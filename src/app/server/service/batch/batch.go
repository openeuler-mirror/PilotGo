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
	//"errors"

	"strconv"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	scommon "gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
)

type Batch = dao.Batch
type Batch2Machine = dao.Batch2Machine
type CreateBatchParam struct {
	Name        string   `json:"Name"`
	Description string   `json:"Descrip"`
	Manager     string   `json:"Manager"`
	DepartName  []string `json:"DepartName"`
	DepartID    []int    `json:"DepartID"`
	Machines    []int    `json:"Machines"`
	DepartIDs   []int    `json:"deptids"`
}

func CreateBatch(batchinfo *CreateBatchParam) error {
	if len(batchinfo.Name) == 0 {
		return errors.New("请输入批次名称")
	}
	ExistNameBool, err := dao.IsExistName(batchinfo.Name)
	if err != nil {
		return err
	}
	if ExistNameBool {
		return errors.New("已存在该名称批次")
	}
	if len(batchinfo.Manager) == 0 {
		return errors.New("创建人未输入")
	}

	if len(batchinfo.Machines) == 0 && len(batchinfo.DepartIDs) == 0 {
		return errors.New("请先选择机器IP或部门")
	}
	// 机器id列表
	var machinelist string
	Departids := make([]int, 0)
	if len(batchinfo.Machines) == 0 {
		// 点选部门创建批次
		var machineids []int
		for _, ids := range batchinfo.DepartIDs {
			Departids = append(Departids, ids)
			common.ReturnSpecifiedDepart(ids, &Departids)
		}

		machines, err := dao.MachineList(Departids)
		if err != nil {
			return err
		}
		for _, mamachine := range machines {
			machineids = append(machineids, mamachine.ID)
		}
		if len(machineids) == 0 {
			return errors.New("该部门下没有机器，请重新确认")
		}
		batchinfo.Machines = machineids
		for _, id := range machineids {
			machinelist = machinelist + "," + strconv.Itoa(id)
			machinelist = strings.Trim(machinelist, ",")
		}
	} else {
		// 点选ip创建批次
		for _, id := range batchinfo.Machines {
			machinelist = machinelist + "," + strconv.Itoa(id)
			machinelist = strings.Trim(machinelist, ",")
		}
	}

	// 机器所属部门ids
	var departIdlist string
	if len(batchinfo.DepartID) == 0 {
		for _, id := range Departids {
			departIdlist = departIdlist + "," + strconv.Itoa(id)
			departIdlist = strings.Trim(departIdlist, ",")
		}
	} else {
		for _, id := range batchinfo.DepartID {
			departIdlist = departIdlist + "," + strconv.Itoa(id)
			departIdlist = strings.Trim(departIdlist, ",")
		}
	}

	// 机器所属部门
	var departNamelist string
	if len(batchinfo.DepartID) == 0 {
		list := dao.DepartIdsToGetDepartNames(Departids)
		departNamelist = strings.Join(list, ",")
	} else {
		List := dao.DepartIdsToGetDepartNames(batchinfo.DepartID)
		departNamelist = strings.Join(List, ",")
	}

	Batch := &dao.Batch{
		Name:        batchinfo.Name,
		Description: batchinfo.Description,
		Manager:     batchinfo.Manager,
		Depart:      departIdlist,
		DepartName:  departNamelist,
		//Machinelist: machinelist,
	}
	err = dao.CreateBatchMessage(Batch)
	if err != nil {
		return err
	}

	//添加数据到Batch2Machine
	for _, id := range batchinfo.Machines {
		err := dao.AddBatch2Machine(dao.Batch2Machine{
			BatchID:       Batch.ID,
			MachineNodeID: uint(id),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func SelectBatch() ([]dao.Batch, error) {
	return dao.GetBatch()
}

// 分页查询所有批次
func GetBatchPaged(offset, size int) (int64, []Batch, error) {
	return dao.GetBatchrPaged(offset, size)
}

func DeleteBatch(ids []int) error {
	for _, value := range ids {
		err := dao.DeleteBatch(value)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}

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
	count, machineIdlist, err := dao.GetMachineIDPaged(offset, size, batchid)
	if err != nil {
		return 0, nil, err
	}
	// 获取机器的所有信息
	machinesInfo := make([]dao.MachineNode, 0)
	for _, macId := range machineIdlist {
		MacInfo, err := dao.MachineData(int(macId.MachineNodeID))
		if err != nil {
			logger.Error(err.Error())
		}
		machinesInfo = append(machinesInfo, MacInfo)
	}

	return count, machinesInfo, nil
}

// from batch get all machines
func GetMachineUUIDS(batchid int) []string {
	if batchid != 0 {
		machineIdlist, err := dao.GetMachineID(batchid)
		if err != nil {
			return nil
		}

		// 获取机器的所有信息
		uuids := make([]string, 0)
		for _, macId := range machineIdlist {
			macuuid, err := dao.MachineIdToUUID(int(macId))
			if err != nil {
				logger.Error(err.Error())
			}
			uuids = append(uuids, macuuid)
		}
		return uuids
	}
	return nil
}

type R interface{}

func BatchProcess(b *scommon.Batch, f func(uuid string) R, it ...interface{}) []R {
	var uuids []string
	if b.MachineUUIDs != nil {
		uuids = b.MachineUUIDs
	} else {
		uuids = GetMachineUUIDS(b.BatchId)
	}
	result := []R{}
	for _, uuid := range uuids {
		r := f(uuid)
		result = append(result, r)
	}

	return result
}

func BatchIds2UUIDs(batchIds []int) (uuids []string) {
	return dao.BatchIds2UUIDs(batchIds)
}
