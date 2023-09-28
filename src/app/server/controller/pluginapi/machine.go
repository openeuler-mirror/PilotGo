package pluginapi

import (
	"gitee.com/PilotGo/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"

	machineservice "gitee.com/PilotGo/PilotGo/app/server/service/machine"
	"gitee.com/PilotGo/PilotGo/sdk/response"
)

func MachineList(c *gin.Context) {
	data, err := machineservice.Machines()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	resp := []*common.MachineNode{}
	for _, item := range data {
		d := &common.MachineNode{
			UUID:       item.UUID,
			IP:         item.IP,
			Department: item.Departname,
			CPUArch:    item.CPU,
			OS:         item.Systeminfo,
			State:      item.State,
		}

		resp = append(resp, d)
	}

	response.Success(c, resp, "获取所有的机器数据")
}
