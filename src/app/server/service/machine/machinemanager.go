package machine

import (
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
)

type MachineNode = dao.MachineNode
type Depart struct {
	ID int `form:"DepartId"`
}

type DeleteUUID struct {
	Deluuid []string `json:"deluuid"`
}

func MachineInfo(depart *Depart, offset, size int) (int64, []dao.Res, error) {

	var TheDeptAndSubDeptIds []int
	common.ReturnSpecifiedDepart(depart.ID, &TheDeptAndSubDeptIds)
	TheDeptAndSubDeptIds = append(TheDeptAndSubDeptIds, depart.ID)
	total, data, err := dao.GetMachinePaged(TheDeptAndSubDeptIds, offset, size)
	return total, data, err
}

func ReturnMachinePaged(departid, offset, size int) (int64, []dao.Res, error) {
	return dao.ReturnMachinePaged(departid, offset, size)
}

func MachineAllData() ([]map[string]string, error) {
	AllData, err := dao.MachineAllData()
	if err != nil {
		return nil, err
	}
	datas := make([]map[string]string, 0)
	for _, data := range AllData {
		datas = append(datas, map[string]string{"uuid": data.UUID, "ip_dept": data.IP + "-" + data.Departname, "ip": data.IP})
	}
	return datas, nil
}

func Machines() ([]dao.Res, error) {
	AllData, err := dao.MachineAllData()
	if err != nil {
		return nil, err
	}
	return AllData, nil
}

func DeleteMachine(Deluuid []string) map[string]string {
	machinelist := make(map[string]string)
	for _, machinedeluuid := range Deluuid {
		if err := dao.DeleteMachine(machinedeluuid); err != nil {
			machinelist[machinedeluuid] = err.Error()
		}
	}
	return machinelist
}

func MachineBasic(uuid string) (ip string, state int, dept string, err error) {
	return dao.MachineBasic(uuid)
}

func UpdateMachineState(uuid string, state int) error {
	return dao.UpdateMachineState(uuid, state)
}

func IsUUIDExist(uuid string) (bool, error) {
	return dao.IsUUIDExist(uuid)
}

// 根据uuid获取部门id
func UUIDForDepartId(uuid string) (int, error) {
	return dao.UUIDForDepartId(uuid)
}

// 更新机器IP及状态
func UpdateMachineIPState(uuid, ip string, state int) error {
	return dao.UpdateMachineIPState(uuid, ip, state)
}

// 新增agent机器
func AddNewMachine(Machine MachineNode) error {
	return dao.AddNewMachine(Machine)
}
