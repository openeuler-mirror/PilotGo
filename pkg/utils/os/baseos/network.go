/*
  * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
*/
package baseos

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"

	aconfig "gitee.com/openeuler/PilotGo/cmd/agent/app/config"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	gnet "github.com/shirou/gopsutil/v3/net"
)

// 网络配置
const NetWorkPath = "/etc/sysconfig/network-scripts"

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
	exitc, result, stde, err := utils.RunCommand("cat /proc/net/arp")
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
	network, err := utils.GetFiles(NetWorkPath, false)
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

	exitc, text, stde, err := utils.RunCommand("cat " + NetWorkPath + "/" + filename)
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

func parseInterfaces(text string) []common.InterfaceInfo {
	interfaces := make([]common.InterfaceInfo, 0)
	lines := strings.Split(text, "\n")
	info := common.InterfaceInfo{}
	for _, line := range lines {

		if len(line) == 0 {
			interfaces = append(interfaces, info)
			info = common.InterfaceInfo{}
			continue
		}

		partitions := strings.Split(strings.TrimSpace(line), ":")
		switch partitions[0] {
		case "GENERAL.DEVICE":
			info.Name = strings.TrimSpace(partitions[1])
			continue
		case "GENERAL.TYPE":
			info.Type = strings.TrimSpace(partitions[1])
			continue
		case "GENERAL.HWADDR":
			info.MACAddress = strings.TrimSpace(partitions[1])
			continue
		case "GENERAL.STATE":
			info.State = strings.TrimSpace(partitions[1])
			continue
		case "IP4.ADDRESS[1]":
			IPAndNetmask := strings.Split(strings.TrimSpace(partitions[1]), "/")
			info.Inet4, info.Netmask = IPAndNetmask[0], IPAndNetmask[1]
			continue
		case "IP4.DNS":
			info.DNS1 = strings.TrimSpace(partitions[1])
			continue
		case "IP4.DNS[1]":
			info.DNS1 = strings.TrimSpace(partitions[1])
			continue
		case "IP4.DNS[2]":
			info.DNS2 = strings.TrimSpace(partitions[1])
			continue
		case "IP4.GATEWAY":
			info.GATEWAY = strings.TrimSpace(partitions[1])
			continue
		case "inet6":
			info.Inet6 = strings.TrimSpace(partitions[1])
			continue
		case "IP6.GATEWAY":
			info.IP6GATEWAY = strings.TrimSpace(partitions[1])
			continue
		default:
			continue
		}
	}
	interfaces = append(interfaces, info)
	return interfaces

}

func (b *BaseOS) GetNetworkConnInfo() (*common.NetworkConfig, error) {
	var network = &common.NetworkConfig{}
	exitc, message, stde, err := utils.RunCommand("nmcli d show")
	if exitc == 0 && message != "" && stde == "" && err == nil {
		interfaces := parseInterfaces(message)
		for i, info := range interfaces {
			name := info.Name
			exitc, message, stde, err := utils.RunCommand("ip addr show " + name + " dynamic")
			if exitc == 0 && message != "" && stde == "" && err == nil {
				interfaces[i].BootProto = "dhcp"
			} else {
				interfaces[i].BootProto = "static"
			}
		}
		network.BootProto = interfaces[0].BootProto
		network.IPAddr = interfaces[0].Inet4
		network.NetMask = interfaces[0].Netmask
		network.GateWay = interfaces[0].GATEWAY
		network.DNS1 = interfaces[0].DNS1
		network.DNS2 = interfaces[0].DNS2

	} else {
		return nil, fmt.Errorf("failed to get network message: %d, %s, %s, %v", exitc, message, stde, err)
	}
	return network, nil
}

// Deprecated
func (b *BaseOS) GetNICName() (string, error) {
	names, err := b.GetNICS()
	for _, value := range names {
		if value.Type == "ethernet" {
			return value.Device, nil
		}
	}
	return "", err
}

func (b *BaseOS) GetNICS() ([]common.NIC, error) {
	var nics []common.NIC
	exitc, stdo, stde, err := utils.RunCommand("nmcli device")
	if err != nil || stde != "" {
		return nics, fmt.Errorf("failed to reload network configuration file: %d, %s, %s, %v", exitc, stdo, stde, err)
	}

	if len(strings.Split(stdo, "\n")) > 1 {
		for _, ns := range strings.Split(stdo, "\n")[1:] {
			n := strings.Fields(ns)
			u := common.NIC{
				Device: n[0],
				Type:   n[1],
				State:  n[2],
			}
			nics = append(nics, u)
		}
	}

	return nics, nil
}

func (b *BaseOS) RestartNetwork(nic string) error {
	exitc, stdo, stde, err := utils.RunCommand("nmcli c reload")
	if exitc == 0 && stdo == "" && stde == "" && err == nil {
	} else {
		return fmt.Errorf("failed to reload network configuration file: %d, %s, %s, %v", exitc, stdo, stde, err)
	}

	exitc2, stdo2, stde2, err2 := utils.RunCommand("nmcli c up " + strings.Split(nic, "-")[1])
	if exitc2 == 0 && stdo2 != "" && stde2 == "" && err2 == nil {
	} else {
		return fmt.Errorf("network configuration file not effective: %d, %s, %s, %v", exitc2, stdo2, stde2, err2)
	}

	return nil
}

func (b *BaseOS) GetHostIp() (string, error) {
	var IP string
	//判断配置文件是否为localhost
	if strings.Split(aconfig.Config().Server.Addr, ":")[0] == "localhost" {
		//通过网卡信息方式获取ip地址
		n, err := b.GetNetworkConnInfo()
		if err != nil {
			return "", err
		}
		IP = n.IPAddr
	} else {
		//获取IP
		conn, err := net.Dial("udp", aconfig.Config().Server.Addr)
		if err != nil {
			return "", err
		}
		defer conn.Close()
		IP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return IP, nil
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
