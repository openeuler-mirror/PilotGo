package pluginapi

import (
	"github.com/gin-gonic/gin"

	"gitee.com/openeuler/PilotGo/app/server/agentmanager"
	"gitee.com/openeuler/PilotGo/app/server/network/jwt"
	"gitee.com/openeuler/PilotGo/app/server/service/auditlog"
	"gitee.com/openeuler/PilotGo/app/server/service/batch"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
)

func Service(c *gin.Context) {
	// TODO: support batch
	d := &common.ServiceStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			log_s := auditlog.New_sub(auditlog.LogTypePlugin, "GetService", uuid, "", auditlog.StatusSuccess, user.ID)
			auditlog.Add(log_s)
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
				auditlog.UpdateMessage(log_s, err.Error())
				auditlog.UpdateStatus(log_s, auditlog.StatusFail)
				logger.Error("获取服务状态失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("获取服务状态结果:%v", serviceInfo)

			serviceSample := common.ServiceInfo{
				ServiceName:         serviceInfo.ServiceName,
				UnitName:            serviceInfo.UnitName,
				UnitType:            serviceInfo.UnitType,
				ServicePath:         serviceInfo.ServicePath,
				ServiceExectStart:   serviceInfo.ServiceExectStart,
				ServiceActiveStatus: serviceInfo.ServiceActiveStatus,
				ServiceLoadedStatus: serviceInfo.ServiceLoadedStatus,
				StartTime:           serviceInfo.ServiceTime,
			}
			re := common.ServiceResult{
				MachineUUID:   uuid,
				MachineIP:     agent.IP,
				ServiceSample: serviceSample,
			}
			return re
		}
		return common.ServiceResult{}
	}

	result := batch.BatchProcess(d.Batch, f, d.ServiceName)
	response.Success(c, gin.H{"service_status": result}, "Success")
}

func StartService(c *gin.Context) {
	// TODO: support batch
	d := &common.ServiceStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			log_s := auditlog.New_sub(auditlog.LogTypePlugin, "StartService", uuid, "", auditlog.StatusSuccess, user.ID)
			auditlog.Add(log_s)
			service_status, Err, err := agent.ServiceStart(d.ServiceName)
			if len(Err) != 0 || err != nil {
				auditlog.UpdateMessage(log_s, Err)
				auditlog.UpdateStatus(log_s, auditlog.StatusFail)
				logger.Error("开启服务失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("开启服务结果:%v", service_status)

			log_s = auditlog.New_sub(auditlog.LogTypePlugin, "GetService", uuid, "", auditlog.StatusSuccess, user.ID)
			auditlog.Add(log_s)
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
				auditlog.UpdateMessage(log_s, err.Error())
				auditlog.UpdateStatus(log_s, auditlog.StatusFail)
				logger.Error("获取服务状态失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("获取服务状态结果:%v", serviceInfo)
			serviceSample := common.ServiceInfo{
				ServiceName:         serviceInfo.ServiceName,
				UnitName:            serviceInfo.UnitName,
				UnitType:            serviceInfo.UnitType,
				ServicePath:         serviceInfo.ServicePath,
				ServiceExectStart:   serviceInfo.ServiceExectStart,
				ServiceActiveStatus: serviceInfo.ServiceActiveStatus,
				ServiceLoadedStatus: serviceInfo.ServiceLoadedStatus,
				StartTime:           serviceInfo.ServiceTime,
			}
			re := common.ServiceResult{
				MachineUUID:   uuid,
				MachineIP:     agent.IP,
				ServiceSample: serviceSample,
			}
			return re
		}
		return common.ServiceResult{}
	}

	result := batch.BatchProcess(d.Batch, f, d.ServiceName)
	response.Success(c, gin.H{"service_start": result}, "Success")
}

func StopService(c *gin.Context) {
	// TODO: support batch
	d := &common.ServiceStruct{}
	err := c.ShouldBind(d)
	if err != nil {
		logger.Debug("bind batch param error:%s", err)
		response.Fail(c, nil, "parameter error")
		return
	}

	user, err := jwt.ParseUser(c)
	if err != nil {
		response.Fail(c, nil, "user token error:"+err.Error())
		return
	}

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			log_s := auditlog.New_sub(auditlog.LogTypePlugin, "StopService", uuid, "", auditlog.StatusSuccess, user.ID)
			auditlog.Add(log_s)
			service_status, Err, err := agent.ServiceStop(d.ServiceName)
			if len(Err) != 0 || err != nil {
				auditlog.UpdateMessage(log_s, Err)
				auditlog.UpdateStatus(log_s, auditlog.StatusFail)
				logger.Error("停止服务失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("停止服务结果:%v", service_status)

			log_s = auditlog.New_sub(auditlog.LogTypePlugin, "GetService", uuid, "", auditlog.StatusSuccess, user.ID)
			auditlog.Add(log_s)
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
				auditlog.UpdateMessage(log_s, err.Error())
				auditlog.UpdateStatus(log_s, auditlog.StatusFail)
				logger.Error("获取服务状态失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("获取服务状态结果:%v", serviceInfo)

			serviceSample := common.ServiceInfo{
				ServiceName:         serviceInfo.ServiceName,
				UnitName:            serviceInfo.UnitName,
				UnitType:            serviceInfo.UnitType,
				ServicePath:         serviceInfo.ServicePath,
				ServiceExectStart:   serviceInfo.ServiceExectStart,
				ServiceActiveStatus: serviceInfo.ServiceActiveStatus,
				ServiceLoadedStatus: serviceInfo.ServiceLoadedStatus,
				StartTime:           serviceInfo.ServiceTime,
			}
			re := common.ServiceResult{
				MachineUUID:   uuid,
				MachineIP:     agent.IP,
				ServiceSample: serviceSample,
			}
			return re
		}
		return common.ServiceResult{}
	}

	result := batch.BatchProcess(d.Batch, f, d.ServiceName)
	response.Success(c, gin.H{"service_stop": result}, "Success")

}
