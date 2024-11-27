/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Fri Apr 7 10:02:19 2023 +0800
 */
package baseos

import (
	"fmt"
	"testing"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"github.com/stretchr/testify/assert"
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

	t.Run("test GetNICS", func(t *testing.T) {
		tmp, err := osobj.GetNICS()
		fmt.Println(tmp)
		assert.Nil(t, err)
		fmt.Println(err)
	})

	t.Run("test GetHostIp", func(t *testing.T) {
		_, err := osobj.GetHostIp()
		assert.Nil(t, err)
	})

	t.Run("test ConfigNetwork", func(t *testing.T) {
		init_config, err := osobj.GetNetworkConnInfo()
		assert.Nil(t, err)
		init_ip_assignment := init_config.BootProto
		init_ipv4_addr := init_config.IPAddr
		init_ipv4_netmask := init_config.NetMask
		init_ipv4_gateway := init_config.GateWay
		init_ipv4_dns1 := init_config.DNS1
		init_ipv4_dns2 := init_config.DNS2
		// http请求地址中的网卡配置参数
		ip_assignment := "static"
		ipv4_addr := "192.168.75.200"
		ipv4_netmask := "255.255.255.0"
		ipv4_gateway := "192.168.75.2"
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
			_, err := utils.UpdateFile(NetWorkPath, nic_name, text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name)
			assert.Nil(t, err)
		case "dhcp":
			text := NetworkDHCP(oldnet)
			_, err := utils.UpdateFile(NetWorkPath, nic_name, text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name)
			assert.Nil(t, err)
		}

		tmp, err := osobj.GetNetworkConnInfo()
		assert.Nil(t, err)
		assert.Equal(t, ip_assignment, tmp.BootProto)
		assert.Equal(t, ipv4_addr, tmp.IPAddr)
		assert.Equal(t, ipv4_netmask, tmp.NetMask)
		assert.Equal(t, ipv4_gateway, tmp.GateWay)
		assert.Equal(t, ipv4_dns1, tmp.DNS1)
		assert.Equal(t, ipv4_dns2, tmp.DNS2)

		// 测试完成将更改参数改回初始值
		switch init_ip_assignment {
		case "static":
			text := NetworkStatic(oldnet, init_ipv4_addr, init_ipv4_netmask, init_ipv4_gateway, init_ipv4_dns1, init_ipv4_dns2)
			_, err := utils.UpdateFile(NetWorkPath, nic_name, text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name)
			assert.Nil(t, err)
		case "dhcp":
			text := NetworkDHCP(oldnet)
			_, err := utils.UpdateFile(NetWorkPath, nic_name, text)
			assert.Nil(t, err)
			err = osobj.RestartNetwork(nic_name)
			assert.Nil(t, err)
		}
	})
}
