package pluginapi

import (
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/batch"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"

	"gitee.com/openeuler/PilotGo-plugins/sdk/common"
)

func InstallPackage(c *gin.Context) {
	param := struct {
		Batch   *common.Batch `json:"batch"`
		Package string        `json:"package"`
	}{}
	if err := c.Bind(&param); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	machines := batch.GetMachines(param.Batch)
	for _, uuid := range machines {
		// TODO: Improve error handling logic
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			logger.Error("cannot find agent %s", uuid)
			continue
		}

		_, _, err := agent.InstallRpm(param.Package)
		if err != nil {
			logger.Error("agent %s install package %s failed: %s", uuid, param.Package, err)
		}
	}

	response.Success(c, nil, "软件包安装完成!")
}

func UninstallPackage(c *gin.Context) {
	param := struct {
		Batch   *common.Batch `json:"batch"`
		Package string
	}{}
	if err := c.Bind(&param); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	machines := batch.GetMachines(param.Batch)
	for _, uuid := range machines {
		// TODO: Improve error handling logic
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			logger.Error("cannot find agent %s", uuid)
			continue
		}

		_, _, err := agent.RemoveRpm(param.Package)
		if err != nil {
			logger.Error("agent %s uninstall package %s failed: %s", uuid, param.Package, err)
		}
	}

	response.Success(c, nil, "软件包安装完成!")
}
