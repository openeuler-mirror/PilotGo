package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func TestNetwork(t *testing.T) {
	var osobj BaseOS

	t.Run("test GetTCP", func(t *testing.T) {
		_, err := osobj.GetTCP()
		assert.Nil(t, err)
	})

	t.Run("test GetUDP", func(t *testing.T) {
		_, err := osobj.GetUDP()
		assert.Nil(t, err)
	})

	t.Run("test GetIOCounter", func(t *testing.T) {
		_, err := osobj.GetIOCounter()
		assert.Nil(t, err)
	})
	// NICname、ip、mac
	t.Run("test GetNICConfig", func(t *testing.T) {
		_, err := osobj.GetNICConfig()
		assert.Nil(t, err)
	})
	// 获取基础网卡配置
	t.Run("test GetNetworkConnInfo", func(t *testing.T) {
		_, err := osobj.GetNetworkConnInfo()
		assert.Nil(t, err)
	})
	// 获取完整旧网卡配置
	t.Run("test ConfigNetworkConnect", func(t *testing.T) {
		_, err := osobj.ConfigNetworkConnect()
		assert.Nil(t, err)
	})

	t.Run("test GetNICName", func(t *testing.T) {
		_, err := osobj.GetNICName()
		assert.Nil(t, err)
	})

	t.Run("test GetHostIp", func(t *testing.T) {
		_, err := osobj.GetHostIp()
		assert.Nil(t, err)
	})

	t.Run("test ConfigNetwork", func(t *testing.T) {
		// http请求地址中的网卡配置参数
		ip_assignment := "static"
		ipv4_addr := "1.1.1.1"
		ipv4_netmask := "255.255.255.0"
		ipv4_gateway := "1.1.1.254"
		ipv4_dns1 := "114.114.114.114"
		ipv4_dns2 := "8.8.8.8"
		// getnicname接口获取的网卡配置参数
		nic_name, err := osobj.GetNICName()
		assert.Nil(t, err)
		// confignetworkconnect接口获取原网卡配置参数
		oldnet, err := osobj.ConfigNetworkConnect()
		assert.Nil(t, err)

		switch ip_assignment {
		case "static":
			text := NetworkStatic(oldnet, ipv4_addr, ipv4_netmask, ipv4_gateway, ipv4_dns1, ipv4_dns2)
			_, err := utils.UpdateFile(global.NetWorkPath, nic_name.(string), text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name.(string))
			assert.Nil(t, err)

		case "dhcp":
			text := NetworkDHCP(oldnet)
			_, err := utils.UpdateFile(global.NetWorkPath, nic_name.(string), text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name.(string))
			assert.Nil(t, err)

		}

		tmp, err := osobj.GetNetworkConnInfo()
		assert.Nil(t, err)
		assert.Equal(t, ip_assignment, tmp.(*common.NetworkConfig).BootProto)
		assert.Equal(t, ipv4_addr, tmp.(*common.NetworkConfig).IPAddr)
		assert.Equal(t, ipv4_netmask, tmp.(*common.NetworkConfig).NetMask)
		assert.Equal(t, ipv4_gateway, tmp.(*common.NetworkConfig).GateWay)
		assert.Equal(t, ipv4_dns1, tmp.(*common.NetworkConfig).DNS1)
		assert.Equal(t, ipv4_dns2, tmp.(*common.NetworkConfig).DNS2)

	})
}
