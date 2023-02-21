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
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/net"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

// 获取当前TCP网络连接信息
func (b *BaseOS) GetTCP() ([]common.NetConnect, error) {
	info, err := net.Connections("tcp")
	if err != nil {
		logger.Error("tcp信息获取失败: ", err)
		return []common.NetConnect{}, fmt.Errorf("tcp信息获取失败")
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
	info, err := net.Connections("udp")
	if err != nil {
		logger.Error("udp信息获取失败: ", err)
		return []common.NetConnect{}, fmt.Errorf("udp信息获取失败")
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
	info, err := net.IOCounters(true)
	if err != nil {
		logger.Error("网络读写字节／包的个数获取失败: ", err)
		return []common.IOCnt{}, fmt.Errorf("网络读写字节／包的个数获取失败")
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
	result, err := utils.RunCommand("cat /proc/net/arp")
	if err != nil {
		logger.Error("网卡信息获取失败: ", err)
		return []common.NetInterfaceCard{}, fmt.Errorf("网卡信息获取失败")
	}
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

// 配置网络连接
func (b *BaseOS) ConfigNetworkConnect() (interface{}, error) {
	filePath := "/home"
	network, err := utils.GetFiles(filePath)
	if err != nil {
		return "", fmt.Errorf("获取网络配置文件失败:%s", err)
	}
	var filename string
	for _, n := range network {
		if ok := strings.Contains(n, "ifcfg-e"); !ok {
			continue
		}
		filename = n
	}

	text, err := utils.RunCommand("cat " + filePath + "/" + filename)
	if err != nil {
		return "", fmt.Errorf("读取网络配置数据失败:%s", err)
	}

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

func (b *BaseOS) GetNetworkConnInfo() (interface{}, error) {
	netPath, err := utils.GetFiles(global.NetWorkPath)
	if err != nil {
		return nil, fmt.Errorf("获取网络配置源文件失败:%s", err)
	}
	var filename string
	for _, n := range netPath {
		if ok := strings.Contains(n, "ifcfg-e"); !ok {
			continue
		}
		filename = n
	}

	result, _ := utils.RunCommand("cat " + global.NetWorkPath + "/" + filename + " | egrep 'BOOTPROTO=.*'")
	ip_assignment_method := strings.Split(result, "=")[1]

	var network = &common.NetworkConfig{}
	switch strings.Replace(ip_assignment_method, "\n", "", -1) {
	case "static":
		tmp, err := utils.RunCommand("cat " + global.NetWorkPath + "/" + filename)
		if err != nil {
			return nil, fmt.Errorf("读取网络配置源数据失败:%s", err)
		}
		lines := strings.Split(tmp, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			strSlice := strings.Split(line, "=")
			if len(strSlice) == 1 {
				continue
			}
			ModuleMatch(strSlice[0], strSlice[1], network)
		}
	case "dhcp":
		IP, err := utils.RunCommand("hostname -I")
		if err != nil {
			return nil, fmt.Errorf("获取IP失败:%s", err)
		}
		str := strings.Split(IP, " ")
		ip := str[0]

		gateway, err := utils.RunCommand("route -n |awk '{print $2}' | sed -n '3p'")
		if err != nil {
			return nil, fmt.Errorf("获取网关失败:%s", err)
		}

		DNS, err := utils.RunCommand("cat /etc/resolv.conf | egrep 'nameserver' | awk '{print $2}'")
		if err != nil {
			return nil, fmt.Errorf("获取DNS失败:%s", err)
		}
		network.DNS1 = strings.Replace(DNS, "\n", "", -1)
		network.BootProto = "dhcp"
		network.IPAddr = ip
		network.GateWay = strings.Replace(gateway, "\n", "", -1)

	default:
		tmp, err := utils.RunCommand("cat " + global.NetWorkPath + "/" + filename)
		if err != nil {
			return nil, fmt.Errorf("读取网络配置源数据失败:%s", err)
		}
		lines := strings.Split(tmp, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			strSlice := strings.Split(line, "=")
			if len(strSlice) == 1 {
				continue
			}
			ModuleMatch(strSlice[0], strSlice[1], network)
		}
	}
	return network, nil
}

func (b *BaseOS) GetNICName() (interface{}, error) {
	network, err := utils.GetFiles(global.NetWorkPath)
	if err != nil {
		return nil, fmt.Errorf("获取网络配置文件失败:%s", err)
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
	_, err := utils.RunCommand("nmcli c reload")
	if err != nil {
		return fmt.Errorf("网络配置文件重载失败:%s", err)
	}

	str := "nmcli c up " + strings.Split(nic, "-")[1]
	_, err = utils.RunCommand(str)
	if err != nil {
		return fmt.Errorf("网络配置文件未生效:%s", err)
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
	IP, err := utils.RunCommand("hostname -I")
	if err != nil {
		return "", err
	}
	str := strings.Split(IP, " ")
	IP = str[0]
	return IP, nil
}
