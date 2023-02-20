package baseos

import (
	"bufio"
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type listService struct {
	Name   string
	LOAD   string
	Active string
	SUB    string
}

const (
	ServiceStart   = 0
	ServiceStop    = 1
	ServiceRestart = 2
)

// 获取服务列表
func (b *BaseOS) GetServiceList() ([]listService, error) {
	list := make([]listService, 0)
	result1, err := utils.RunCommand("systemctl list-units --all|grep 'loaded[ ]*active' | awk 'NR>2{print $1\" \" $2\" \" $3\" \" $4}'")
	if err != nil {
		logger.Error("获取服务列表命令运行失败: ", err)
		return []listService{}, err
	}
	result2, err := utils.RunCommand("systemctl list-units --all|grep 'not-found' | awk 'NR>2{print $2\" \" $3\" \" $4\" \" $5}'")
	if err != nil {
		logger.Error("获取服务列表命令运行失败: ", err)
		return []listService{}, err
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
		tmp := listService{}
		tmp.Name = str[0]
		tmp.LOAD = str[1]
		tmp.Active = str[2]
		tmp.SUB = str[3]
		list = append(list, tmp)
	}
	return list, nil
}

//operate 0,1,2分别代表开启，关闭，重启

func (b *BaseOS) GetServiceStatus(service string) (string, error) {
	var build strings.Builder
	build.WriteString("systemctl is-active ")
	build.WriteString(service)
	command := build.String()
	output, err := utils.RunCommand(command)
	return output, err
}
func verifyStatus(output string, operate int) bool {
	var judge bool
	if strings.Contains(string(output), "inactive") {
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
	result, err := utils.RunCommand(command)
	if err != nil || len(result) != 0 {
		logger.Error("开启服务命令运行失败: ", err)
		return fmt.Errorf(" %s 服务启动失败%s", service, err)
	}
	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("获取服务状态失败")
		return fmt.Errorf(" %s 服务重新启动失败%s", service, err)
	}
	if !verifyStatus(output, ServiceStart) {
		logger.Error("开启服务命令运行结果无效!")
		return fmt.Errorf(" %s 服务启动失败%s", service, err)
	}
	return nil
}
func (b *BaseOS) StopService(service string) error {
	var build strings.Builder
	build.WriteString("systemctl stop ")
	build.WriteString(service)
	command := build.String()
	result, err := utils.RunCommand(command)
	if err != nil || len(result) != 0 {
		logger.Error("关闭服务命令运行失败: ", err)
		return fmt.Errorf(" %s 服务关闭失败%s", service, err)
	}
	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("获取服务状态失败")
		return fmt.Errorf(" %s 服务重新启动失败%s", service, err)
	}
	if !verifyStatus(output, ServiceStop) {
		logger.Error("关闭服务命令运行结果无效!")
		return fmt.Errorf(" %s 服务关闭失败%s", service, err)
	}
	return nil
}
func (b *BaseOS) RestartService(service string) error {
	var build strings.Builder
	build.WriteString("systemctl restart ")
	build.WriteString(service)
	command := build.String()
	result, err := utils.RunCommand(command)
	if err != nil || len(result) != 0 {
		logger.Error("重启服务命令运行失败: ", err)
		return fmt.Errorf(" %s 服务重新启动失败%s", service, err)
	}
	output, err := b.GetServiceStatus(service)
	if err != nil {
		logger.Error("获取服务状态失败")
		return fmt.Errorf(" %s 服务重新启动失败%s", service, err)
	}
	if !verifyStatus(output, ServiceRestart) {
		logger.Error("重启服务命令运行结果无效!")
		return fmt.Errorf(" %s 服务重新启动失败%s", service, err)
	}
	return nil
}
