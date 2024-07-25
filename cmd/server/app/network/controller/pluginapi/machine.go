package pluginapi

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"

	machineservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/machine"
	"gitee.com/openeuler/PilotGo/sdk/response"
)

func MachineList(c *gin.Context) {
	data, err := machineservice.MachineAll()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	resp := []*common.MachineNode{}
	for _, item := range data {
		d := &common.MachineNode{
			UUID:        item.UUID,
			IP:          item.IP,
			Department:  item.Departname,
			CPUArch:     item.CPU,
			OS:          item.Systeminfo,
			RunStatus:   item.Runstatus,
			MaintStatus: item.Maintstatus,
		}

		resp = append(resp, d)
	}

	response.Success(c, resp, "获取所有的机器数据")
}

func MachineInfoByUUID(c *gin.Context) {
	machine_uuid := c.Query("machine_uuid")
	data, err := machineservice.MachineInfoByUUID(machine_uuid)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, data, "获取所有的机器数据")
}
