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
 * LastEditTime: 2022-04-13 01:51:51
 * Description: get agent network information.
 ******************************************************************************/
package os

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/net"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type NetConnect struct {
	Localaddr  string
	Remoteaddr string
	Status     string
	Uids       []int32
	Pid        int32
}
type IOCnt struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}
type NetInterfaceCard struct {
	Name    string
	IPAddr  string
	MacAddr string
}

type NetworkConfig struct {
	NetworkType        string `json:"type"` //以太网、无线网
	ProxyMethod        string `json:"proxy_method"`
	BrowserOnly        string `json:"browser_only"`
	DefRoute           string `json:"defroute"`
	IPV4_Failure_Fatal string `json:"ipv4_failure_fatal"`
	Name               string `json:"name"`   //接口名称
	UUID               string `json:"uuid"`   //唯一识别码
	Device             string `json:"device"` //网卡设备名字
	OnBoot             string `json:"onboot"` //是否随网络服务启动当前网卡生效

	IPV6Init           string `json:"ipv6_init"` //ipv6是否启用
	IPV6_Autoconf      string `json:"ipv6_autoconf"`
	IPV6_DefRoute      string `json:"ipv6_defroute"`
	IPV6_Failure_Fatal string `json:"ipv6_failure_fatal"`
	IPv6_Addr_Gen_Mode string `json:"ipv6_addr_gen_mode"`

	MachineUUID string `json:"macUUID"`
	BootProto   string `json:"BOOTPROTO"` //dhcp或者static
	IPAddr      string `json:"IPADDR"`
	NetMask     string `json:"NETMASK"`
	GateWay     string `json:"GATEWAY"`
	DNS1        string `json:"DNS1"`
	DNS2        string `json:"DNS2"`
}

//获取当前TCP网络连接信息
func GetTCP() ([]NetConnect, error) {
	info, err := net.Connections("tcp")
	if err != nil {
		logger.Error("tcp信息获取失败: ", err)
		return []NetConnect{}, fmt.Errorf("tcp信息获取失败")
	}
	tcpConf := make([]NetConnect, 0)
	for _, value := range info {
		tmp := NetConnect{}
		tmp.Localaddr = value.Laddr.IP + ":" + fmt.Sprint(value.Laddr.Port)
		tmp.Remoteaddr = value.Raddr.IP + ":" + fmt.Sprint(value.Raddr.Port)
		tmp.Status = value.Status
		tmp.Uids = value.Uids
		tmp.Pid = value.Pid
		tcpConf = append(tcpConf, tmp)
	}
	return tcpConf, nil
}

//获取当前UDP网络连接信息
func GetUDP() ([]NetConnect, error) {
	info, err := net.Connections("udp")
	if err != nil {
		logger.Error("udp信息获取失败: ", err)
		return []NetConnect{}, fmt.Errorf("udp信息获取失败")
	}
	tcpConf := make([]NetConnect, 0)
	for _, value := range info {
		tmp := NetConnect{}
		tmp.Localaddr = value.Laddr.IP + ":" + fmt.Sprint(value.Laddr.Port)
		tmp.Remoteaddr = value.Raddr.IP + ":" + fmt.Sprint(value.Raddr.Port)
		tmp.Status = value.Status
		tmp.Uids = value.Uids
		tmp.Pid = value.Pid
		tcpConf = append(tcpConf, tmp)
	}
	return tcpConf, nil
}

//获取网络读写字节／包的个数
func GetIOCounter() ([]IOCnt, error) {
	info, err := net.IOCounters(true)
	if err != nil {
		logger.Error("网络读写字节／包的个数获取失败: ", err)
		return []IOCnt{}, fmt.Errorf("网络读写字节／包的个数获取失败")
	}
	IOConf := make([]IOCnt, 0)
	for _, value := range info {
		tmp := IOCnt{}
		tmp.Name = value.Name
		tmp.BytesSent = value.BytesSent
		tmp.BytesRecv = value.BytesRecv
		tmp.PacketsSent = value.PacketsSent
		tmp.PacketsRecv = value.PacketsRecv
		IOConf = append(IOConf, tmp)
	}
	return IOConf, nil
}

func GetNICConfig() ([]NetInterfaceCard, error) {
	NICConfig := make([]NetInterfaceCard, 0)
	result, err := utils.RunCommand("cat /proc/net/arp")
	if err != nil {
		logger.Error("网卡信息获取失败: ", err)
		return []NetInterfaceCard{}, fmt.Errorf("网卡信息获取失败")
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
		tmp := NetInterfaceCard{}
		tmp.IPAddr = y[0]
		tmp.MacAddr = x[0]
		tmp.Name = z[0]
		NICConfig = append(NICConfig, tmp)
	}
	return NICConfig, nil
}

// 配置网络连接
func ConfigNetworkConnect() (interface{}, error) {
	filePath := "/home"
	network, err := GetFiles(filePath)
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

func GetNetworkConnInfo() (interface{}, error) {
	netPath, err := GetFiles(global.NetWorkPath)
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

	var network = &NetworkConfig{}
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

func GetNICName() (interface{}, error) {
	network, err := GetFiles(global.NetWorkPath)
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

func RestartNetwork(nic string) error {
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

func ModuleMatch(key string, value string, network *NetworkConfig) {
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
		network.NetMask = LenToSubNetMask(prefix)
	}
}

// 网络长度转换成子网掩码
func LenToSubNetMask(subnet int) string {
	var buff bytes.Buffer
	for i := 0; i < subnet; i++ {
		buff.WriteString("1")
	}
	for i := subnet; i < 32; i++ {
		buff.WriteString("0")
	}
	masker := buff.String()
	a, _ := strconv.ParseUint(masker[:8], 2, 64)
	b, _ := strconv.ParseUint(masker[8:16], 2, 64)
	c, _ := strconv.ParseUint(masker[16:24], 2, 64)
	d, _ := strconv.ParseUint(masker[24:32], 2, 64)
	resultMask := fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
	return resultMask

}
