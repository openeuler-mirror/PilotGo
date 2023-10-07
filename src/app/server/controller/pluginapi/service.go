package pluginapi

import (
	"github.com/gin-gonic/gin"

	"gitee.com/PilotGo/PilotGo/app/server/agentmanager"
	"gitee.com/PilotGo/PilotGo/app/server/service/batch"
	"gitee.com/PilotGo/PilotGo/sdk/common"
	"gitee.com/PilotGo/PilotGo/sdk/logger"
	"gitee.com/PilotGo/PilotGo/sdk/response"
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

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
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

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			service_status, Err, err := agent.ServiceStart(d.ServiceName)
			if len(Err) != 0 || err != nil {
				logger.Error("开启服务失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("开启服务结果:%v", service_status)
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
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

	f := func(uuid string) batch.R {
		agent := agentmanager.GetAgent(uuid)
		if agent != nil {
			service_status, Err, err := agent.ServiceStop(d.ServiceName)
			if len(Err) != 0 || err != nil {
				logger.Error("停止服务失败!, agent:%s, command:%s", uuid, d.ServiceName)
			}
			logger.Debug("停止服务结果:%v", service_status)
			serviceInfo, err := agent.GetService(d.ServiceName)
			if err != nil {
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
