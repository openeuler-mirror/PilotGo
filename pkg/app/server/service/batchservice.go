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
 * LastEditTime: 2022-05-18 16:25:41
 * Description: 批次管理业务逻辑
 ******************************************************************************/
package service

import (
	"errors"
	"strconv"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func CreateBatch(batchinfo *model.CreateBatch) error {
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
			ReturnSpecifiedDepart(ids, &Departids)
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

	Batch := model.Batch{
		Name:        batchinfo.Name,
		Description: batchinfo.Description,
		Manager:     batchinfo.Manager,
		Depart:      departIdlist,
		DepartName:  departNamelist,
		Machinelist: machinelist,
	}
	dao.CreateBatch(Batch)

	return nil
}

// TODO: *[]model.Batch 应该定义为指针数组
func GetBatches(query *model.PaginationQ) (*[]model.Batch, int64, error) {

	batch := model.Batch{}
	list, tx := batch.ReturnBatch(query)

	total, err := CrudAll(query, tx, list)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func DeleteBatch(ids []int) error {
	for _, value := range ids {
		dao.DeleteBatch(value)
	}
	return nil
}

func UpdateBatch(batchid int, name, description string) error {
	if !dao.IsExistID(batchid) {
		return errors.New("不存在该批次")
	}

	dao.UpdateBatch(batchid, name, description)
	return nil
}

func GetBatchMachines(batchid int) ([]model.MachineNode, error) {
	machinelist := dao.GetMachineID(batchid)
	machineIdlist := utils.String2Int(machinelist) // 获取批次里所有机器的id

	// 获取机器的所有信息
	machinesInfo := make([]model.MachineNode, 0)
	for _, macId := range machineIdlist {
		MacInfo, err := dao.MachineData(macId)
		if err != nil {
			return machinesInfo, err
		}
		machinesInfo = append(machinesInfo, MacInfo)
	}

	return machinesInfo, nil
}

func SelectBatch() ([]model.Batch, error) {
	return dao.GetBatch(), nil
}
