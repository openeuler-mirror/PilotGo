package pluginapi

import (
	"github.com/gin-gonic/gin"
	machineservice "openeuler.org/PilotGo/PilotGo/pkg/app/server/service/machine"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func MachineList(c *gin.Context) {
	data, err := machineservice.MachineAllData()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, data, "获取所有的机器数据")
}
