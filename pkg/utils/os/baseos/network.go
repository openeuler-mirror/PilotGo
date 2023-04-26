/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: wanghao
 * Date: 2022-02-17 02:43:29
 * LastEditTime: 2023-02-21 16:00:35
 * Description: get agent network information.
 ******************************************************************************/
package baseos

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	gnet "github.com/shirou/gopsutil/net"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取当前TCP网络连接信息
func (b *BaseOS) GetTCP() ([]common.NetConnect, error) {
	info, err := gnet.Connections("tcp")
	if err != nil {
		logger.Error("failed to get tcp message: %s", err)
		return []common.NetConnect{}, fmt.Errorf("failed to get tcp message")
	}
	tcpConf := make([]common.NetConnect, 0)
	for _, value := range info {
		tmp := common.NetConnect{}
		tmp.Localaddr = value.Laddr.IP + ":" + fmt.Sprint(value.Laddr.Port)
		tmp.Remoteaddr = value.Raddr.IP + ":" + fmt.Sprint(value.Raddr.Port)
		tmp.Status = value.Status
		tmp.Uids = value.Uids
		tmp.Pid = value.Pid
		tcpConf = append(tcpConf, tmp)
	}
	return tcpConf, nil
}

// 获取当前UDP网络连接信息
func (b *BaseOS) GetUDP() ([]common.NetConnect, error) {
	info, err := gnet.Connections("udp")
	if err != nil {
		logger.Error("failed to get udp message: %s", err)
		return []common.NetConnect{}, fmt.Errorf("failed to get udp message")
	}
	tcpConf := make([]common.NetConnect, 0)
	for _, value := range info {
		tmp := common.NetConnect{}
		tmp.Localaddr = value.Laddr.IP + ":" + fmt.Sprint(value.Laddr.Port)
		tmp.Remoteaddr = value.Raddr.IP + ":" + fmt.Sprint(value.Raddr.Port)
		tmp.Status = value.Status
		tmp.Uids = value.Uids
		tmp.Pid = value.Pid
		tcpConf = append(tcpConf, tmp)
	}
	return tcpConf, nil
}

// 获取网络读写字节／包的个数
func (b *BaseOS) GetIOCounter() ([]common.IOCnt, error) {
	info, err := gnet.IOCounters(true)
	if err != nil {
		logger.Error("failed to get number of bytes/packets for network read/write: %s", err)
		return []common.IOCnt{}, fmt.Errorf("failed to get number of bytes/packets for network read/write")
	}
	IOConf := make([]common.IOCnt, 0)
	for _, value := range info {
		tmp := common.IOCnt{}
		tmp.Name = value.Name
		tmp.BytesSent = value.BytesSent
		tmp.BytesRecv = value.BytesRecv
		tmp.PacketsSent = value.PacketsSent
		tmp.PacketsRecv = value.PacketsRecv
		IOConf = append(IOConf, tmp)
	}
	return IOConf, nil
}

func (b *BaseOS) GetNICConfig() ([]common.NetInterfaceCard, error) {
	NICConfig := make([]common.NetInterfaceCard, 0)
	exitc, result, stde, err := utils.RunCommandnew("cat /proc/net/arp")
	if exitc == 0 && result != "" && stde == "" && err == nil {
		reader := strings.NewReader(result)
		scanner := bufio.NewScanner(reader)
		for {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)

			reg1 := regexp.MustCompile(`[A-Fa-f0-9]{2}(:[A-Fa-f0-9]{2}){5}`)
			reg2 := regexp.MustCompile(`([0-9]{1,3}.){3}[0-9]{1,3}`)
			reg3 := regexp.MustCompile(`[a-zA-Z0-9]+$`)
			x := reg1.FindAllString(line, -1)
			y := reg2.FindAllString(line, -1)
			z := reg3.FindAllString(line, -1)
			if x == nil || y == nil || z == nil {
				continue
			}
			tmp := common.NetInterfaceCard{}
			tmp.IPAddr = y[0]
			tmp.MacAddr = x[0]
			tmp.Name = z[0]
			NICConfig = append(NICConfig, tmp)
		}
		return NICConfig, nil
	}
	logger.Error("faile to get network card message: %d, %s, %s, %v", exitc, result, stde, err)
	return nil, fmt.Errorf("faile to get network card message: %d, %s, %s, %v", exitc, result, stde, err)

}

// 配置网络连接
func (b *BaseOS) ConfigNetworkConnect() ([]map[string]string, error) {
	network, err := utils.GetFiles(global.NetWorkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get network configuration file: %s", err)
	}
	var filename string
	for _, n := range network {
		if ok := strings.Contains(n, "ifcfg-e"); !ok {
			continue
		}
		filename = n
	}

	exitc, text, stde, err := utils.RunCommandnew("cat " + global.NetWorkPath + "/" + filename)
	if exitc == 0 && text != "" && stde == "" && err == nil {
		var oldnet []map[string]string
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			strSlice := strings.Split(line, "=")
			if len(strSlice) == 1 {
				continue
			}
			if strings.Contains(strSlice[0], "#") {
				continue
			}
			net := map[string]string{
				strSlice[0]: strSlice[1],
			}
			oldnet = append(oldnet, net)
		}
		return oldnet, nil
	}
	return nil, fmt.Errorf("failed to read network configuration data: %d, %s, %s, %v", exitc, text, stde, err)

}

func (b *BaseOS) GetNetworkConnInfo() (*common.NetworkConfig, error) {
	netPath, err := utils.GetFiles(global.NetWorkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get network configuration source file: %s", err)
	}
	var filename string
	for _, n := range netPath {
		if ok := strings.Contains(n, "ifcfg-e"); !ok {
			continue
		}
		filename = n
	}

	exitc, result, stde, err := utils.RunCommandnew("cat " + global.NetWorkPath + "/" + filename + " | egrep 'BOOTPROTO=.*'")
	if exitc == 0 && result != "" && stde == "" && err == nil {
		ip_assignment_method := strings.Split(result, "=")[1]

		var network = &common.NetworkConfig{}
		switch strings.Replace(ip_assignment_method, "\n", "", -1) {
		case "static":
			exitc, tmp, stde, err := utils.RunCommandnew("cat " + global.NetWorkPath + "/" + filename)
			if exitc == 0 && tmp != "" && stde == "" && err == nil {
				lines := strings.Split(tmp, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					strSlice := strings.Split(line, "=")
					if len(strSlice) == 1 {
						continue
					}
					ModuleMatch(strSlice[0], strSlice[1], network)
				}
			} else {
				return nil, fmt.Errorf("failed to read network configuration source file: %d, %s, %s, %v", exitc, tmp, stde, err)
			}
		case "dhcp":
			conn, err := net.Dial("udp", "openeuler.org:80")
			if err != nil {
				return nil, fmt.Errorf("failed to get IP: %s", err)
			}
			defer conn.Close()
			ip := strings.Split(conn.LocalAddr().String(), ":")[0]

			exitc, gateway, stde, err := utils.RunCommandnew("route -n |awk '{print $2}' | sed -n '3p'")
			if exitc == 0 && gateway != "" && stde == "" && err == nil {
			} else {
				return nil, fmt.Errorf("failed to get gateway: %d, %s, %s, %v", exitc, gateway, stde, err)
			}

			exitc, DNS, stde, err := utils.RunCommandnew("cat /etc/resolv.conf | egrep 'nameserver' | awk '{print $2}'")
			if exitc == 0 && DNS != "" && stde == "" && err == nil {
			} else {
				return nil, fmt.Errorf("failed to get dns: %d, %s, %s, %v", exitc, DNS, stde, err)
			}

			network.DNS1 = strings.Replace(DNS, "\n", "", -1)
			network.BootProto = "dhcp"
			network.IPAddr = ip
			network.GateWay = strings.Replace(gateway, "\n", "", -1)
		default:
			exitc, tmp, stde, err := utils.RunCommandnew("cat " + global.NetWorkPath + "/" + filename)
			if exitc == 0 && tmp != "" && stde == "" && err == nil {
				lines := strings.Split(tmp, "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					strSlice := strings.Split(line, "=")
					if len(strSlice) == 1 {
						continue
					}
					ModuleMatch(strSlice[0], strSlice[1], network)
				}
			} else {
				return nil, fmt.Errorf("failed to read network configuration source data: %d, %s, %s, %v", exitc, tmp, stde, err)
			}
		}
		return network, nil
	}
	return nil, fmt.Errorf("failed to get BOOTPROTO: %d, %s, %s, %v", exitc, result, stde, err)
}

func (b *BaseOS) GetNICName() (interface{}, error) {
	network, err := utils.GetFiles(global.NetWorkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get network configuration file: %s", err)
	}
	var filename string
	for _, n := range network {
		if ok := strings.Contains(n, "ifcfg-e"); !ok {
			continue
		}
		filename = n
	}

	return filename, nil
}

func (b *BaseOS) RestartNetwork(nic string) error {
	exitc, stdo, stde, err := utils.RunCommandnew("nmcli c reload")
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
	} else {
		return fmt.Errorf("failed to reload network configuration file: %d, %s, %s, %v", exitc, stdo, stde, err)
	}

	exitc2, stdo2, stde2, err2 := utils.RunCommandnew("nmcli c up " + strings.Split(nic, "-")[1])
	if exitc2 == 0 && stdo2 != "" && stde2 == "" && err2 == nil {
	} else {
		return fmt.Errorf("network configuration file not effective: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}

	return nil
}

func ModuleMatch(key string, value string, network *common.NetworkConfig) {
	if key == "IPADDR" {
		network.IPAddr = value
	} else if key == "NETMASK" {
		network.NetMask = value
	} else if key == "GATEWAY" {
		network.GateWay = value
	} else if key == "DNS1" {
		network.DNS1 = value
	} else if key == "DNS2" {
		network.DNS2 = value
	} else if key == "BOOTPROTO" {
		network.BootProto = value
	} else if key == "PREFIX" {
		prefix, _ := strconv.Atoi(value)
		network.NetMask = common.LenToSubNetMask(prefix)
	}
}

func (b *BaseOS) GetHostIp() (string, error) {
	conn, err := net.Dial("udp", "openeuler.org:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

// dhcp方式配置网络
func NetworkDHCP(net []map[string]string) (text string) {
	for _, n := range net {
		for key, value := range n {
			if key == "BOOTPROTO" {
				text += key + "=" + "dhcp" + "\n"
			} else if key == "IPADDR" {
				break
			} else if key == "NETMASK" {
				break
			} else if key == "GATEWAY" {
				break
			} else if key == "DNS1" {
				break
			} else if key == "DNS2" {
				break
			} else {
				text += key + "=" + value + "\n"
			}
		}
	}
	return
}

// static方式配置网络
func NetworkStatic(net []map[string]string, ip string, netmask string, gateway string, dns1 string, dns2 string) (text string) {
	for _, n := range net {
		for key, value := range n {
			if key == "BOOTPROTO" {
				text += key + "=" + "static" + "\n"
			} else if key == "IPADDR" {
				text += key + "=" + ip + "\n"
			} else if key == "NETMASK" {
				text += key + "=" + netmask + "\n"
			} else if key == "GATEWAY" {
				text += key + "=" + gateway + "\n"
			} else if key == "DNS1" {
				text += key + "=" + dns1 + "\n"
			} else if key == "DNS2" && len(dns2) != 0 {
				text += key + "=" + dns2 + "\n"
			} else {
				text += key + "=" + value + "\n"
			}
		}
	}
	if ok := strings.Contains(text, "IPADDR"); !ok {
		t := "IPADDR" + "=" + ip + "\n"
		text += t
	}
	if ok := strings.Contains(text, "NETMASK"); !ok {
		t := "NETMASK" + "=" + netmask + "\n"
		text += t
	}
	if ok := strings.Contains(text, "GATEWAY"); !ok {
		t := "GATEWAY" + "=" + gateway + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS1"); !ok {
		t := "DNS1" + "=" + dns1 + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS2"); !ok {
		if len(dns2) != 0 {
			t := "DNS2" + "=" + dns2 + "\n"
			text += t
		}
	}
	return
}
