package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

type PackageStruct struct {
	Batch   *common.Batch `json:"batch"`
	Package string
}

func InstallPackage(c *gin.Context) {
	param := PackageStruct{}
	if err := c.Bind(&param); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)

		if agent != nil {
			data, resp_message_err, err := agent.InstallRpm(param.Package)
			if resp_message_err != "" {
				logger.Error(resp_message_err)
			}
			if err != nil {
				logger.Error("agent %s install package %s failed: %s", uuid, param.Package, err)
			}
			logger.Debug("install package on agent result:%v", data)
			return data
		}
		return ""
	}

	result := batch.BatchProcess(param.Batch, f, param.Package)
	response.Success(c, result, "软件包安装完成!")
}

func UninstallPackage(c *gin.Context) {
	param := PackageStruct{}
	if err := c.Bind(&param); err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)

		if agent != nil {
			data, resp_message_err, err := agent.RemoveRpm(param.Package)
			if resp_message_err != "" {
				logger.Error(resp_message_err)
			}
			if err != nil {
				logger.Error("agent %s uninstall package %s failed: %s", uuid, param.Package, err)
			}
			logger.Debug("uninstall package on agent result:%v", data)
			return data
		}
		return ""
	}

	result := batch.BatchProcess(param.Batch, f, param.Package)
	response.Success(c, result, "软件包卸载完成!")
}
