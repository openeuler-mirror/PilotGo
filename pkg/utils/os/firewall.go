package os

import (
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func Config() (interface{}, error) {
	nic_interface, err := GetNICName()
	if err != nil {
		return nil, fmt.Errorf("未获取到网卡名称")
	}

	firewall_state, err := utils.RunCommand("firewall-cmd --state")
	if err != nil {
		firewalldConfig := map[string]interface{}{
			"status":      "not running",
			"nic":         strings.Split(nic_interface.(string), "-")[1],
			"defaultZone": nil,
			"zones":       nil,
			"services":    nil,
		}
		return firewalldConfig, nil
	}

	zone_default, err := utils.RunCommand("firewall-cmd --get-default-zone")
	if err != nil {
		firewalldConfig := map[string]interface{}{
			"status":      strings.Replace(firewall_state, "\n", "", -1),
			"nic":         strings.Split(nic_interface.(string), "-")[1],
			"defaultZone": nil,
			"zones":       nil,
			"services":    nil,
		}
		return firewalldConfig, nil
	}

	zones, err := utils.RunCommand("firewall-cmd --get-zones")
	if err != nil {
		firewalldConfig := map[string]interface{}{
			"status":      strings.Replace(firewall_state, "\n", "", -1),
			"nic":         strings.Split(nic_interface.(string), "-")[1],
			"defaultZone": strings.Replace(zone_default, "\n", "", -1),
			"zones":       nil,
			"services":    nil,
		}
		return firewalldConfig, nil
	}
	Zones := strings.Split(strings.Replace(zones, "\n", "", -1), " ")

	services, err := utils.RunCommand("firewall-cmd --get-services")
	if err != nil {
		firewalldConfig := map[string]interface{}{
			"status":      strings.Replace(firewall_state, "\n", "", -1),
			"nic":         strings.Split(nic_interface.(string), "-")[1],
			"defaultZone": strings.Replace(zone_default, "\n", "", -1),
			"zones":       Zones,
			"services":    nil,
		}
		return firewalldConfig, nil
	}
	Services := strings.Split(strings.Replace(services, "\n", "", -1), " ")

	firewalldConfig := map[string]interface{}{
		"status":      strings.Replace(firewall_state, "\n", "", -1),
		"nic":         strings.Split(nic_interface.(string), "-")[1],
		"defaultZone": strings.Replace(zone_default, "\n", "", -1),
		"zones":       Zones,
		"services":    Services,
	}
	return firewalldConfig, nil
}

func FirewalldSetDefaultZone(zone string) (interface{}, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --set-default-zone=%v", zone))
	if err != nil {
		return nil, fmt.Errorf("防火墙默认区域更改失败")
	}
	return tmp, nil
}

func FirewalldZoneConfig(zone string) (interface{}, error) {
	conf, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --list-all", zone))
	if err != nil {
		return nil, fmt.Errorf("防火墙未运行")
	}

	var firewall = &FirewalldCMDList{}
	lines := strings.Split(conf, "\n")
	for _, line := range lines {
		if ok := strings.Contains(line, "sources"); ok {
			firewall.Sources = strings.Split(strings.Split(line, ": ")[1], " ")
		} else if ok := strings.Contains(line, "services"); ok {
			firewall.Service = strings.Split(strings.Split(line, ": ")[1], " ")
		} else if ok := strings.Contains(line, " ports"); ok {
			ports := strings.Split(strings.Split(line, ": ")[1], " ")
			datas := make([]map[string]string, 0)
			for _, port := range ports {
				strSlice := strings.Split(port, "/")
				if len(strSlice) == 1 {
					continue
				}
				data := map[string]string{
					"port":     strSlice[0],
					"protocol": strSlice[1],
				}
				datas = append(datas, data)
			}
			firewall.Ports = datas
		} else {
			continue
		}
	}

	return firewall, nil
}

func FirewalldSourceAdd(zone, source string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --add-source=%v", zone, source))
	if err != nil {
		return fmt.Errorf("INVALID_ADDR:%s", source)
	}
	return nil
}

func FirewalldSourceRemove(zone, source string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --remove-source=%v", zone, source))
	if err != nil {
		return fmt.Errorf("UNKNOWN_SOURCE: '%s' is not in any zone", source)
	}
	return nil
}

func FirewalldServiceAdd(zone, service string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --add-service=%v", zone, service))
	if err != nil {
		return err
	}
	return nil
}

func FirewalldServiceRemove(zone, service string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --remove-service=%v", zone, service))
	if err != nil {
		return err
	}
	return nil
}

func Restart() bool {
	tmp, _ := utils.RunCommand("service firewalld restart")
	return len(tmp) == 0
}

func Stop() bool {
	tmp, _ := utils.RunCommand("service firewalld stop")
	return len(tmp) == 0
}

func DelZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-port=%v/%v", zone, port, protocol))
	if err != nil {
		return tmp, fmt.Errorf("防火墙未运行")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("重新加载防火墙失败")
	}
	return tmpp, nil
}

func AddZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-port=%v/%v", zone, port, protocol))
	if err != nil {
		return tmp, fmt.Errorf("防火墙未运行")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("重新加载防火墙失败")
	}
	return tmpp, nil
}

type FirewalldCMDList struct {
	Service []string    `json:"services"` // 列出允许通过这个防火墙的服务
	Ports   interface{} `json:"ports"`    //列出允许通过这个防火墙的目标端口。（即 需要对外开放的端口）
	Sources []string    `json:"sources"`  // 允许通过的IP或mac
}
