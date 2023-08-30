package baseos

import (
	"bufio"
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

const (
	ServiceStart   = 0
	ServiceStop    = 1
	ServiceRestart = 2
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

// operate 0,1,2分别代表开启，关闭，重启
// TODO: 软件包在未安装情况下，'systemctl is-active'返回的结果和软件包已安装且服务未运行时的结果一致
func (b *BaseOS) GetServiceStatus(service string) (string, error) {
	var build strings.Builder
	build.WriteString("systemctl is-active ")
	build.WriteString(service)
	command := build.String()
	exitc, tmp, stde, err := utils.RunCommand(command)
	output := strings.Trim(tmp, "\n")
	switch output {
	case "active", "inactive":
		return output, nil
	default:
		return "", fmt.Errorf("failed to get service status: %d, %s, %s, %v", exitc, tmp, stde, err)
	}
}

func (b *BaseOS) GetService(service string) *common.ServiceInfo {
	var build strings.Builder
	build.WriteString("systemctl status ")
	build.WriteString(service)
	command := build.String()
	exitc, tmp, stde, err := utils.RunCommand(command)
	if err != nil {
		logger.Error("failed to get service status: %d, %s, %s, %v", exitc, tmp, stde, err)
	}
	if stde != "" {
		logger.Error(stde)
		return nil
	}
	serviceInfo := parseServiceInfo(tmp)
	return serviceInfo
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

func verifyStatus(output string, operate int) bool {
	var judge bool
	if output == "inactive" {
		switch operate {
		case 0:
			judge = false
		case 1:
			judge = true
		case 2:
			judge = false
		}
	} else {
		switch operate {
		case 0:
			judge = true
		case 1:
			judge = false
		case 2:
			judge = true
		}
	}
	return judge
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

	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("failed to retrieve the status of the service")
		return fmt.Errorf("failed to start the %s service: %s", service, err)
	}
	if !verifyStatus(output, ServiceStart) {
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

	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("failed to get the status of the service")
		return fmt.Errorf("failed to stop the %s service: %s", service, err)
	}
	if !verifyStatus(output, ServiceStop) {
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

	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("failed to get the status of the service")
		return fmt.Errorf("failed to restart the %s service: %s", service, err)
	}
	if !verifyStatus(output, ServiceRestart) {
		logger.Error("failed to execute the command to restart the service!")
		return fmt.Errorf("failed to restart the %s service: %s", service, err)
	}
	return nil
}
