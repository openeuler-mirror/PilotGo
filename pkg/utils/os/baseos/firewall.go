package baseos

import (
	"fmt"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func (b *BaseOS) Config() (common.FireWalldConfig, error) {
	nic_interface, err := b.GetNICName()
	firewalldConfig := common.FireWalldConfig{}
	firewalldConfig.Set()
	if err != nil {
		return firewalldConfig, fmt.Errorf("failed to get network card name")
	}

	exitc, firewall_state, stde, err := utils.RunCommandnew("firewall-cmd --state")
	if exitc != 0 && firewall_state == "" && strings.Replace(stde, "\n", "", -1) == "not running" && err == nil {
		firewalldConfig.Status = "not running"
		firewalldConfig.Nic = strings.Split(nic_interface, "-")[1]
		return firewalldConfig, nil
	}

	exitc, zone_default, stde, err := utils.RunCommandnew("firewall-cmd --get-default-zone")
	if exitc == 0 && zone_default != "" && stde == "" && err == nil {
	} else {
		firewalldConfig.Status = strings.Replace(firewall_state, "\n", "", -1)
		firewalldConfig.Nic = strings.Split(nic_interface, "-")[1]
		return firewalldConfig, nil
	}

	exitc, zones, stde, err := utils.RunCommandnew("firewall-cmd --get-zones")
	if exitc == 0 && zones != "" && stde == "" && err == nil {
	} else {
		firewalldConfig.Status = strings.Replace(firewall_state, "\n", "", -1)
		firewalldConfig.Nic = strings.Split(nic_interface, "-")[1]
		firewalldConfig.DefaultZone = strings.Replace(zone_default, "\n", "", -1)
		return firewalldConfig, nil
	}
	Zones := strings.Split(strings.Replace(zones, "\n", "", -1), " ")

	exitc, services, stde, err := utils.RunCommandnew("firewall-cmd --get-services")
	if exitc == 0 && services != "" && stde == "" && err == nil {
	} else {
		firewalldConfig.Status = strings.Replace(firewall_state, "\n", "", -1)
		firewalldConfig.Nic = strings.Split(nic_interface, "-")[1]
		firewalldConfig.DefaultZone = strings.Replace(zone_default, "\n", "", -1)
		firewalldConfig.Zones = Zones
		return firewalldConfig, nil
	}
	Services := strings.Split(strings.Replace(services, "\n", "", -1), " ")
	firewalldConfig.Status = strings.Replace(firewall_state, "\n", "", -1)
	firewalldConfig.Nic = strings.Split(nic_interface, "-")[1]
	firewalldConfig.DefaultZone = strings.Replace(zone_default, "\n", "", -1)
	firewalldConfig.Zones = Zones
	firewalldConfig.Services = Services

	return firewalldConfig, nil
}

func (b *BaseOS) FirewalldSetDefaultZone(zone string) (string, error) {
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --set-default-zone=%v", zone))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		return strings.Replace(stdo, "\n", "", -1), nil
	}
	return "", fmt.Errorf("failed to change default zone of firewall: %d, %s, %s, %v", exitc, stdo, stde, err)
}

func (b *BaseOS) FirewalldZoneConfig(zone string) (*common.FirewalldCMDList, error) {
	exitc, conf, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --zone=%v --list-all", zone))
	if exitc == 0 && conf != "" && stde == "" && err == nil {
		var firewall = &common.FirewalldCMDList{}
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
	return nil, fmt.Errorf("failed to get zone config: %d, %s, %s, %v", exitc, conf, stde, err)
}

func (b *BaseOS) FirewalldSourceAdd(zone, source string) error {
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-source=%v", zone, source))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			return nil
		}
		return fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return fmt.Errorf("failed to add source %s: %d, %s, %s, %v", source, exitc, stdo, stde, err)

}

func (b *BaseOS) FirewalldSourceRemove(zone, source string) error {
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-source=%v", zone, source))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			return nil
		}
		return fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return fmt.Errorf("failed to remove source %s: %d, %s, %s, %v", source, exitc, stdo, stde, err)
}

func (b *BaseOS) FirewalldServiceAdd(zone, service string) error {
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-service=%v", zone, service))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			return nil
		}
		return fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return fmt.Errorf("failed to add service %s: %d, %s, %s, %v", service, exitc, stdo, stde, err)
}

func (b *BaseOS) FirewalldServiceRemove(zone, service string) error {
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-service=%v", zone, service))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			return nil
		}
		return fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return fmt.Errorf("failed to remove service %s: %d, %s, %s, %v", service, exitc, stdo, stde, err)
}

func (b *BaseOS) Restart() bool {
	exitc, stdo, stde, err := utils.RunCommandnew("systemctl restart firewalld.service")
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
		return true
	}
	return false
}

func (b *BaseOS) Stop() bool {
	exitc, stdo, stde, err := utils.RunCommandnew("systemctl stop firewalld.service")
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
		return true
	}
	return false
}

func (b *BaseOS) DelZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-port=%v/%v", zone, port, protocol))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			stdo2 = strings.Replace(stdo2, "\n", "", -1)
			return stdo2, nil
		}
		return "", fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return "", fmt.Errorf("failed to remove zone port %s/%s: %d, %s, %s, %v", port, protocol, exitc, stdo, stde, err)
}

func (b *BaseOS) AddZonePort(zone, port, protocol string) (string, error) { //zone = block dmz drop external home internal public trusted work
	exitc, stdo, stde, err := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-port=%v/%v", zone, port, protocol))
	if exitc == 0 && strings.Replace(stdo, "\n", "", -1) == "success" && stde == "" && err == nil {
		exitc2, stdo2, stde2, err2 := utils.RunCommandnew("firewall-cmd --reload")
		if exitc2 == 0 && strings.Replace(stdo2, "\n", "", -1) == "success" && stde2 == "" && err2 == nil {
			stdo2 = strings.Replace(stdo2, "\n", "", -1)
			return stdo2, nil
		}
		return "", fmt.Errorf("failed to reload firewall: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}
	return "", fmt.Errorf("failed to add zone port %s/%s: %d, %s, %s, %v", port, protocol, exitc, stdo, stde, err)
}

// TODO: firewall完善zone interface的添加和删除接口，完善firewall接口重复add会报错的逻辑
func (b *BaseOS) AddZoneInterface(zone, NIC string) (string, error) {
	_, tmp, _, _ := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-interface=%v", zone, NIC))

	return strings.Replace(tmp, "\n", "", -1), nil
}

func (b *BaseOS) DelZoneInterface(zone, NIC string) (string, error) {
	_, tmp, _, _ := utils.RunCommandnew(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-interface=%v", zone, NIC))

	return strings.Replace(tmp, "\n", "", -1), nil
}
