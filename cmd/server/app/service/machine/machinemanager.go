/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package machine

import (
	"errors"
	"time"

	eventSDK "gitee.com/openeuler/PilotGo-plugins/event/sdk"
	departservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/depart"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	commonSDK "gitee.com/openeuler/PilotGo/sdk/common"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
)

// 机器运行状态
const (
	OnlineStatus  = "online"
	OfflineStatus = "offline"
)

// 机器维护状态
const (
	NormalStatus      = "normal"
	MaintenanceStatus = "maintenance"
)

type MachineNode = dao.MachineNode
type Res = dao.Res
type Depart struct {
	ID int `form:"DepartId"`
}

type DeleteUUID struct {
	Deluuid []string `json:"deluuid"`
}

// 新增agent机器
func AddMachine(Machine *MachineNode) error {
	return Machine.Add()
}

func MachineInfo(depart *Depart, offset, size int, search string) (int64, []dao.Res, error) {
	var TheDeptAndSubDeptIds []int
	departservice.ReturnSpecifiedDepart(depart.ID, &TheDeptAndSubDeptIds)
	TheDeptAndSubDeptIds = append([]int{depart.ID}, TheDeptAndSubDeptIds...)
	total, data, err := dao.GetMachinePaged(TheDeptAndSubDeptIds, offset, size, search)
	return total, data, err
}

// 插件调用
func MachineAllData() ([]map[string]string, error) {
	AllData, err := dao.MachineAll()
	if err != nil {
		return nil, err
	}
	datas := make([]map[string]string, 0)
	for _, data := range AllData {
		datas = append(datas, map[string]string{"uuid": data.UUID, "ip_dept": data.IP + "-" + data.Departname, "ip": data.IP})
	}
	return datas, nil
}

func MachineAll() ([]dao.Res, error) {
	AllData, err := dao.MachineAll()
	return AllData, err
}

// 删除机器，删除之前先校验uuid是否存在
func DeleteMachine(Deluuid []string) map[string]string {
	machinelist := make(map[string]string)
	for _, machinedeluuid := range Deluuid {
		node, err := MachineInfoByUUID(machinedeluuid)
		if err != nil {
			machinelist[machinedeluuid] = err.Error()
		}
		if node.ID != 0 {
			//删除机器批次关系表数据
			if err := dao.DeleteMachineBatch(node.ID); err != nil {
				machinelist[machinedeluuid] = err.Error()
				continue
			}
			if err := dao.DeleteMachine(machinedeluuid); err != nil {
				machinelist[machinedeluuid] = err.Error()
			}
			// 发布“平台移除主机”事件
			msgData := commonSDK.MessageData{
				MsgType:     eventSDK.MsgHostRemove,
				MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgHostRemove),
				TimeStamp:   time.Now(),
				Data: eventSDK.MDHostChange{
					IP:     node.IP,
					OS:     node.Systeminfo,
					CPU:    node.CPU,
					Status: OfflineStatus,
				},
			}
			msgDataString, _ := msgData.ToMessageDataString()
			ms := commonSDK.EventMessage{
				MessageType: eventSDK.MsgHostRemove,
				MessageData: msgDataString,
			}
			plugin.PublishEvent(ms)
		} else {
			machinelist[machinedeluuid] = errors.New("该机器不存在").Error()
		}
	}
	return machinelist
}

// 根据uuid查找机器信息
func MachineInfoByUUID(uuid string) (*Res, error) {
	node, err := dao.MachineInfoByUUID(uuid)
	if err != nil {
		return nil, err
	}

	depart, err := dao.GetDepartById(node.DepartId)
	r := &Res{
		ID:          node.ID,
		Departid:    node.DepartId,
		Departname:  depart.Depart,
		IP:          node.IP,
		UUID:        node.MachineUUID,
		CPU:         node.CPU,
		Runstatus:   node.RunStatus,
		Maintstatus: node.MaintStatus,
		Systeminfo:  node.Systeminfo,
	}
	return r, err
}

func UpdateMachine(uuid string, ma *MachineNode) error {
	if ma.MaintStatus != "" {
		err := IsStatus(ma.MaintStatus)
		if err != nil {
			return err
		}
	}
	return dao.UpdateMachine(uuid, ma)
}

func ModifyMachineDepart(MachineID string, DepartID int) error {
	//查询部门节点是否存在
	depart, err := dao.GetDepartById(DepartID)
	if err != nil {
		return err
	}
	if depart.ID == 0 {
		return errors.New("此部门不存在")
	}
	ResIds := utils.String2Int(MachineID)
	for _, id := range ResIds {
		machine, err := dao.MachineInfo(id)
		if err != nil {
			return err
		}
		Ma := &dao.MachineNode{DepartId: DepartID}
		if DepartID != global.UncateloguedDepartId {
			Ma.MaintStatus = NormalStatus
		} else {
			Ma.MaintStatus = MaintenanceStatus
		}
		err = dao.UpdateMachine(machine.MachineUUID, Ma)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsStatus(maintstatus string) error {
	if maintstatus != MaintenanceStatus && maintstatus != NormalStatus {
		return errors.New("维护状态字段不存在")
	}
	return nil
}

func UpdateMaintStatus(MachineIDs string, maintstatus string) ([]int, error) {
	//判断状态名是否正确
	err := IsStatus(maintstatus)
	if err != nil {
		return nil, err
	}
	ids := utils.String2Int(MachineIDs)
	var result []int
	for _, v := range ids {
		uuid, err := dao.MachineIdToUUID(v)
		if err != nil {
			result = append(result, v)
			continue
		}
		err = dao.UpdateMachine(uuid, &MachineNode{MaintStatus: maintstatus})
		if err != nil {
			result = append(result, v)
		}
	}
	return result, nil
}
