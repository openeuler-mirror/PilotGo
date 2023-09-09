package nestos

import (
	"fmt"
	"net"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/baseos"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

type NestOS struct {
	baseos.BaseOS
}

func (b *NestOS) GetNetworkConnInfo() (*common.NetworkConfig, error) {
	var network = &common.NetworkConfig{}

	conn, err := net.Dial("udp", "openeuler.org:80")
	if err != nil {
		return nil, fmt.Errorf("failed to get IP: %s", err)
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]

	exitc, gateway, stde, err := utils.RunCommand("route -n |awk '{print $2}' | sed -n '3p'")
	if exitc == 0 && gateway != "" && stde == "" && err == nil {
	} else {
		return nil, fmt.Errorf("failed to get gateway: %d, %s, %s, %v", exitc, gateway, stde, err)
	}

	exitc, DNS, stde, err := utils.RunCommand("cat /etc/resolv.conf | egrep 'nameserver' | awk '{print $2}'")
	if exitc == 0 && DNS != "" && stde == "" && err == nil {
	} else {
		return nil, fmt.Errorf("failed to get dns: %d, %s, %s, %v", exitc, DNS, stde, err)
	}

	exitc, NETMASK, stde, err := utils.RunCommand("ifconfig |grep netmask|awk '{print $4}'|awk 'NR==1'")
	if exitc == 0 && NETMASK != "" && stde == "" && err == nil {
	} else {
		return nil, fmt.Errorf("failed to get dns: %d, %s, %s, %v", exitc, DNS, stde, err)
	}
	network.NetMask = strings.Replace(NETMASK, "\n", "", -1)
	network.DNS1 = strings.Replace(DNS, "\n", "", -1)
	network.BootProto = "dhcp"
	network.IPAddr = ip
	network.GateWay = strings.Replace(gateway, "\n", "", -1)
	return network, nil
}
