package service

import (
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
)

func MachineInfo(depart *model.Depart, query *model.PaginationQ) (interface{}, int, error) {

	var TheDeptAndSubDeptIds []int
	ReturnSpecifiedDepart(depart.ID, &TheDeptAndSubDeptIds)
	TheDeptAndSubDeptIds = append(TheDeptAndSubDeptIds, depart.ID)
	machinelist, err := dao.MachineList(TheDeptAndSubDeptIds)
	if err != nil {
		return nil, 0, err
	}
	lens := len(machinelist)
	data, err := DataPaging(query, machinelist, lens)
	if err != nil {
		return nil, 0, err
	}
	return data, lens, nil
}

func MachineAllData() ([]map[string]string, error) {
	AllData := dao.MachineAllData()
	datas := make([]map[string]string, 0)
	for _, data := range AllData {
		datas = append(datas, map[string]string{"uuid": data.UUID, "ip_dept": data.IP + "-" + data.Departname, "ip": data.IP})
	}
	return datas, nil
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
