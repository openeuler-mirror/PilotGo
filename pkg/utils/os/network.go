package os

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/net"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
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
		reg2 := regexp.MustCompile(`([0-9]{1,3}.){3}[0-9]`)
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
