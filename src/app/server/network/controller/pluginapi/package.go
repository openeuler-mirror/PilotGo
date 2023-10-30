package pluginapi

import (
	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Install Package",
	}
	auditlog.Add(log)

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)

		if agent != nil {
			log_s := &auditlog.AuditLog{
				LogUUID:    uuidservice.New().String(),
				ParentUUID: log.LogUUID,
				Module:     auditlog.ModulePlugin,
				Status:     auditlog.StatusOK,
				UserID:     u.ID,
				Action:     "Install Package",
				Message:    "agentuuid:" + uuid,
			}
			auditlog.Add(log_s)

			data, resp_message_err, err := agent.InstallRpm(param.Package)
			if resp_message_err != "" {
				auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+resp_message_err)
				auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
				logger.Error(resp_message_err)
			}
			if err != nil {
				auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
				auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
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

	u, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "Uninstall Package",
	}
	auditlog.Add(log)

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)

		if agent != nil {
			log_s := &auditlog.AuditLog{
				LogUUID:    uuidservice.New().String(),
				ParentUUID: log.LogUUID,
				Module:     auditlog.ModulePlugin,
				Status:     auditlog.StatusOK,
				UserID:     u.ID,
				Action:     "Uninstall Package",
				Message:    "agentuuid:" + uuid,
			}
			auditlog.Add(log_s)
			data, resp_message_err, err := agent.RemoveRpm(param.Package)
			if resp_message_err != "" {
				auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+resp_message_err)
				auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
				logger.Error(resp_message_err)
			}
			if err != nil {
				auditlog.UpdateMessage(log_s, "agentuuid:"+uuid+err.Error())
				auditlog.UpdateStatus(log_s, auditlog.StatusFailed)
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
