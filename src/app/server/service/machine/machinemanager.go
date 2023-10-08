package machine

import (
	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
)

type MachineNode = dao.MachineNode
type Depart struct {
	ID int `form:"DepartId"`
}

type DeleteUUID struct {
	Deluuid []string `json:"deluuid"`
}

func MachineInfo(depart *Depart, query *common.PaginationQ) (interface{}, int, error) {

	var TheDeptAndSubDeptIds []int
	common.ReturnSpecifiedDepart(depart.ID, &TheDeptAndSubDeptIds)
	TheDeptAndSubDeptIds = append(TheDeptAndSubDeptIds, depart.ID)
	machinelist, err := dao.MachineList(TheDeptAndSubDeptIds)
	if err != nil {
		return nil, 0, err
	}
	lens := len(machinelist)
	data, err := common.DataPaging(query, machinelist, lens)
	if err != nil {
		return nil, 0, err
	}
	return data, lens, nil
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
