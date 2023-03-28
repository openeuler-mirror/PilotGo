package baseos

import (
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

func (b *BaseOS) Config() (interface{}, error) {
	nic_interface, err := b.GetNICName()
	if err != nil {
		return nil, fmt.Errorf("failed to get network card name")
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

func (b *BaseOS) FirewalldSetDefaultZone(zone string) (interface{}, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --set-default-zone=%v", zone))
	if err != nil {
		return nil, fmt.Errorf("failed to change default zone of firewall")
	}
	return tmp, nil
}

func (b *BaseOS) FirewalldZoneConfig(zone string) (interface{}, error) {
	conf, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --list-all", zone))
	if err != nil {
		return nil, fmt.Errorf("firewall not running")
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

func (b *BaseOS) FirewalldSourceAdd(zone, source string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --add-source=%v", zone, source))
	if err != nil {
		return fmt.Errorf("INVALID_ADDR:%s", source)
	}
	return nil
}

func (b *BaseOS) FirewalldSourceRemove(zone, source string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --remove-source=%v", zone, source))
	if err != nil {
		return fmt.Errorf("UNKNOWN_SOURCE: '%s' is not in any zone", source)
	}
	return nil
}

func (b *BaseOS) FirewalldServiceAdd(zone, service string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --add-service=%v", zone, service))
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseOS) FirewalldServiceRemove(zone, service string) error {
	_, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --zone=%v --remove-service=%v", zone, service))
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseOS) Restart() bool {
	tmp, _ := utils.RunCommand("service firewalld restart")
	return len(tmp) == 0
}

func (b *BaseOS) Stop() bool {
	tmp, _ := utils.RunCommand("service firewalld stop")
	return len(tmp) == 0
}

func (b *BaseOS) DelZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-port=%v/%v", zone, port, protocol))
	if err != nil {
		return tmp, fmt.Errorf("firewall not running")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("failed to reload firewall")
	}
	return tmpp, nil
}

func (b *BaseOS) AddZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-port=%v/%v", zone, port, protocol))
	if err != nil {
		return tmp, fmt.Errorf("firewall not running")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("failed to reload firewall")
	}
	return tmpp, nil
}

type FirewalldCMDList struct {
	Service []string    `json:"services"` // 列出允许通过这个防火墙的服务
	Ports   interface{} `json:"ports"`    //列出允许通过这个防火墙的目标端口。（即 需要对外开放的端口）
	Sources []string    `json:"sources"`  // 允许通过的IP或mac
}
