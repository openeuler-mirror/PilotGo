package common

import (
	"bytes"
	"fmt"
	"strconv"
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

type NIC struct {
	Device string
	Type   string
	State  string
}

type NetworkConfig struct {
	NetworkType        string `json:"TYPE"` //以太网、无线网
	ProxyMethod        string `json:"PROXY_METHOD"`
	BrowserOnly        string `json:"BROWSER_ONLY"`
	DefRoute           string `json:"DEFROUTE"`
	IPV4_Failure_Fatal string `json:"IPV4_FAILURE_FATAL"`
	Name               string `json:"NAME"`   //接口名称
	UUID               string `json:"UUID"`   //唯一识别码
	Device             string `json:"DEVICE"` //网卡设备名字
	OnBoot             string `json:"ONBOOT"` //是否随网络服务启动当前网卡生效

	IPV6Init           string `json:"IPV6INIT"` //ipv6是否启用
	IPV6_Autoconf      string `json:"IPV6_AUTOCONF"`
	IPV6_DefRoute      string `json:"IPV6_DEFROUTE"`
	IPV6_Failure_Fatal string `json:"IPV6_FAILURE_FATAL"`
	IPv6_Addr_Gen_Mode string `json:"IPV6_ADDR_GEN_MODE"`

	MachineUUID string `json:"macUUID"`
	BootProto   string `json:"BOOTPROTO"` //dhcp或者static
	IPAddr      string `json:"IPADDR"`
	NetMask     string `json:"NETMASK"`
	GateWay     string `json:"GATEWAY"`
	DNS1        string `json:"DNS1"`
	DNS2        string `json:"DNS2"`
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
