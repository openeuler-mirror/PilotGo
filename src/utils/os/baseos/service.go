package baseos

import (
	"bufio"
	"fmt"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
	"gitee.com/openeuler/PilotGo/utils/os/common"
)

// 获取服务列表
func (b *BaseOS) GetServiceList() ([]common.ListService, error) {
	list := make([]common.ListService, 0)
	exitc, result1, stde, err := utils.RunCommand("systemctl list-units --all|grep 'loaded[ ]*' | awk 'NR>2{print $1\" \" $2\" \" $3\" \" $4}'")
	if exitc == 0 && result1 != "" && stde == "" && err == nil {
	} else {
		logger.Error("failed to execute the command to get the list of services: %d, %s, %s, %v", exitc, result1, stde, err)
		return nil, fmt.Errorf("failed to execute the command to get the list of services: %d, %s, %s, %v", exitc, result1, stde, err)
	}

	exitc, result2, stde, err := utils.RunCommand("systemctl list-units --all|grep 'not-found' | awk 'NR>2{print $2\" \" $3\" \" $4\" \" $5}'")
	if exitc == 0 && result2 != "" && stde == "" && err == nil {
	} else {
		logger.Error("the command to get the list of services has failed to run: %d, %s, %s, %v", exitc, result1, stde, err)
		return nil, fmt.Errorf("the command to get the list of services has failed to run: %d, %s, %s, %v", exitc, result1, stde, err)
	}

	result := result1 + result2
	reader := strings.NewReader(result)
	scanner := bufio.NewScanner(reader)
	for {

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		str := strings.Fields(line)
		tmp := common.ListService{}
		tmp.Name = str[0]
		tmp.LOAD = str[1]
		tmp.Active = str[2]
		tmp.SUB = str[3]
		list = append(list, tmp)
	}
	return list, nil
}

func (b *BaseOS) GetService(service string) (*common.ServiceInfo, error) {
	var build strings.Builder
	build.WriteString("systemctl status ")
	build.WriteString(service)
	command := build.String()
	exitc, tmp, stde, err := utils.RunCommand(command)
	if err != nil {
		logger.Error("failed to get service status: %d, %s, %s, %v", exitc, tmp, stde, err)
		return nil, err
	}
	if stde != "" {
		logger.Error(stde)
	}
	serviceInfo := parseServiceInfo(tmp)
	return serviceInfo, nil
}

func parseServiceInfo(tmp string) *common.ServiceInfo {
	service := &common.ServiceInfo{}
	serviceResult := strings.Split(tmp, "\n")

	if len(serviceResult) > 2 {
		service.UnitName = strings.Trim(strings.Split(serviceResult[0], " - ")[0], "● ")
		service.ServiceName = strings.Split(service.UnitName, ".")[0]
		switch strings.Split(service.UnitName, ".")[1] {
		case "service":
			service.UnitType = common.ServiceUnit
			break
		case "socket":
			service.UnitType = common.SocketUnit
			break
		case "target":
			service.UnitType = common.TargetUnit
			break
		case "mount":
			service.UnitType = common.MountUnit
			break
		case "automount":
			service.UnitType = common.AutomountUnit
			break
		case "path":
			service.UnitType = common.PathUnit
			break
		case "time":
			service.UnitType = common.TimeUnit
			break
		default:
			service.UnitType = common.UNKnown
			break
		}
		for i, value := range strings.Split(strings.Trim(serviceResult[1], " "), " ") {
			if i == 2 {
				service.ServicePath = strings.Trim(value, "(;")
			}

			if i == 3 {
				switch strings.Trim(value, ";") {
				case "enabled":
					service.ServiceLoadedStatus = common.ServiceLoadedStatusEnabled
					break
				case "disabled":
					service.ServiceLoadedStatus = common.ServiceLoadedStatusDisabled
					break
				case "static":
					service.ServiceLoadedStatus = common.ServiceLoadedStatusStatic
					break
				case "mask":
					service.ServiceLoadedStatus = common.ServiceLoadedStatusMask
					break
				default:
					service.ServiceLoadedStatus = common.ServiceLoadedStatusUnknown
					break
				}
			}
		}

		for i, value := range strings.Split(strings.Trim(serviceResult[2], " "), " ") {
			if i == 1 {
				switch value {
				case "inactive":
					service.ServiceActiveStatus = common.ServiceActiveStatusInactive
					break
				case "active":
					service.ServiceActiveStatus = "active"
					break
				default:
					service.ServiceActiveStatus = common.ServiceActiveStatusUnknown
					break
				}
			}
			if i == 2 && service.ServiceActiveStatus == "active" {
				switch strings.Trim(value, "()") {
				case "running":
					service.ServiceActiveStatus = common.ServiceActiveStatusRunning
					break
				case "exited":
					service.ServiceActiveStatus = common.ServiceActiveStatusExited
					break
				case "waiting":
					service.ServiceActiveStatus = common.ServiceActiveStatusWaiting
					break
				}
			}
			if i == 5 {
				service.ServiceTime = value
			}
			if i == 6 {
				service.ServiceTime = service.ServiceTime + " " + value
			}
		}

		if len(serviceResult) > 3 && strings.Split(strings.Trim(serviceResult[3], " "), " ")[0] == "Process:" {
			service.ServiceExectStart = strings.Split(strings.Split(strings.Trim(serviceResult[3], " "), " ")[2], "=")[1]
		}
		return service
	}
	return nil
}

func (b *BaseOS) StartService(service string) error {
	var build strings.Builder
	build.WriteString("systemctl start ")
	build.WriteString(service)
	command := build.String()

	exitc, result, stde, err := utils.RunCommand(command)
	if exitc == 0 && result == "" && stde == "" && err == nil {
	} else {
		logger.Error("failed to execute the command to start the service: %d, %s, %s, %v", exitc, result, stde, err)
		return fmt.Errorf("failed to execute the command to start the service: %d, %s, %s, %v", exitc, result, stde, err)
	}

	serviceInfo, err := b.GetService(service)
	if err != nil {
		logger.Error("failed to retrieve the status of the service")
		return fmt.Errorf("failed to start the %s service: %s", service, err)
	}
	if serviceInfo.ServiceActiveStatus != common.ServiceActiveStatusRunning {
		logger.Error("the command to start the service has produced an invalid result!")
		return fmt.Errorf("failed to start the %s service: %s", service, err)
	}
	return nil
}
func (b *BaseOS) StopService(service string) error {
	var build strings.Builder
	build.WriteString("systemctl stop ")
	build.WriteString(service)
	command := build.String()
	exitc, result, stde, err := utils.RunCommand(command)
	if exitc == 0 && result == "" && stde == "" && err == nil {
	} else {
		logger.Error("failed to execute the command to stop the service: %d, %s, %s, %v", exitc, result, stde, err)
		return fmt.Errorf("failed to execute the command to stop the service: %d, %s, %s, %v", exitc, result, stde, err)
	}

	serviceInfo, err := b.GetService(service)
	if err != nil {
		logger.Error("failed to get the status of the service")
		return fmt.Errorf("failed to stop the %s service: %s", service, err)
	}
	if serviceInfo.ServiceActiveStatus != common.ServiceActiveStatusInactive {
		logger.Error("the command to stop the service has produced an invalid result!")
		return fmt.Errorf("failed to stop the %s service: %s", service, err)
	}
	return nil
}
func (b *BaseOS) RestartService(service string) error {
	var build strings.Builder
	build.WriteString("systemctl restart ")
	build.WriteString(service)
	command := build.String()
	exitc, result, stde, err := utils.RunCommand(command)
	if exitc == 0 && result == "" && stde == "" && err == nil {
	} else {
		logger.Error("failed to execute the command to restart the service: %d, %s, %s, %v", exitc, result, stde, err)
		return fmt.Errorf("failed to execute the command to restart the service: %d, %s, %s, %v", exitc, result, stde, err)
	}

	serviceInfo, err := b.GetService(service)
	if err != nil {
		logger.Error("failed to get the status of the service")
		return fmt.Errorf("failed to restart the %s service: %s", service, err)
	}
	if serviceInfo.ServiceActiveStatus != common.ServiceActiveStatusRunning {
		logger.Error("failed to execute the command to restart the service!")
		return fmt.Errorf("failed to restart the %s service: %s", service, err)
	}
	return nil
}