package pluginapi

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"

	machineservice "gitee.com/openeuler/PilotGo/app/server/service/machine"
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
