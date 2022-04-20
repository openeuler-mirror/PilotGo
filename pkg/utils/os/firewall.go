package os

import (
	"fmt"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

func Config() ([]string, error) {
	tmp, err := utils.RunCommand("firewall-cmd --list-all")
	if err != nil {
		return nil, fmt.Errorf("FirewallD is not running")
	}

	tmp = strings.TrimSpace(tmp)
	t := strings.Split(tmp, "\n")

	return t, nil
}

func Restart() bool {
	tmp, _ := utils.RunCommand("service firewalld restart")
	return len(tmp) == 0
}

func Stop() bool {
	tmp, _ := utils.RunCommand("service firewalld stop")
	return len(tmp) == 0
}

func DelZonePort(zone, port string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --remove-port=%v/tcp", zone, port))
	if err != nil {
		return tmp, fmt.Errorf("FirewallD is not running")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("重新加载防火墙失败")
	}
	return tmpp, nil
}

func AddZonePort(zone, port string) (string, error) { //zone = block dmz drop external home internal public trusted work
	tmp, err := utils.RunCommand(fmt.Sprintf("firewall-cmd --permanent --zone=%v --add-port=%v/tcp", zone, port))
	if err != nil {
		return tmp, fmt.Errorf("FirewallD is not running")
	}
	tmpp, err := utils.RunCommand("firewall-cmd --reload")
	tmpp = strings.Replace(tmpp, "\n", "", -1)
	if err != nil {
		return "", fmt.Errorf("重新加载防火墙失败")
	}
	return tmpp, nil
}
