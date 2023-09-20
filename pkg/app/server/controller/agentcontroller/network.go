/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-02-17 02:43:29
 * LastEditTime: 2023-02-21 16:01:22
 * Description: provide agent network manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/baseos"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func NetTCPHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net_tcp, err := agent.NetTCP()
	if err != nil {
		response.Fail(c, nil, "获取当前TCP网络连接信息失败!")
		return
	}
	response.Success(c, gin.H{"net_tcp": net_tcp}, "Success")
}
func NetUDPHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net_udp, err := agent.NetUDP()
	if err != nil {
		response.Fail(c, nil, "获取当前UDP网络连接信息失败!")
		return
	}
	response.Success(c, gin.H{"net_udp": net_udp}, "Success")
}
func NetIOCounterHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net_io, err := agent.NetIOCounter()
	if err != nil {
		response.Fail(c, nil, "获取网络读写字节/包的个数失败!")
		return
	}
	response.Success(c, gin.H{"net_io": net_io}, "Success")
}
func NetNICConfigHandler(c *gin.Context) {
	uuid := c.Query("uuid")

	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net_nic, err := agent.NetNICConfig()
	if err != nil {
		response.Fail(c, nil, "获取网卡配置失败!")
		return
	}
	response.Success(c, gin.H{"net_nic": net_nic}, "Success")
}

func GetAgentNetworkConnect(c *gin.Context) {
	uuid := c.Query("uuid")
	agent := agentmanager.GetAgent(uuid)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	net, Err, err := agent.GetNetWorkConnInfo()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	response.Success(c, net, "获取到网络连接信息")
}

func ConfigNetworkConnect(c *gin.Context) {
	var network common.NetworkConfig
	c.Bind(&network)

	ip_assignment := network.BootProto
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ip分配方式不能为空")
		return
	}
	ipv4_addr := network.IPAddr
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4地址不能为空")
		return
	}
	ipv4_netmask := network.NetMask
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4子网掩码不能为空")
		return
	}
	if ok := strings.Contains(ipv4_netmask, "."); !ok {
		prefix, _ := strconv.Atoi(ipv4_netmask)
		ipv4_netmask = common.LenToSubNetMask(prefix)
	}
	ipv4_gateway := network.GateWay
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4网关不能为空")
		return
	}
	ipv4_dns1 := network.DNS1
	if len(ip_assignment) == 0 {
		response.Fail(c, nil, "ipv4 DNS1 不能为空")
		return
	}

	agent := agentmanager.GetAgent(network.MachineUUID)
	if agent == nil {
		response.Fail(c, nil, "获取uuid失败!")
		return
	}

	nic_name, Err, err := agent.GetNICName()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}

	oldnet, Err, err := agent.GetNetWorkConnectInfo()
	if len(Err) != 0 || err != nil {
		response.Fail(c, nil, Err)
		return
	}
	// oldnets1 := oldnet.([]interface{})
	// var oldnets2 []map[string]interface{}
	// for _, m := range *oldnet {
	// 	m1 := m.(map[string]interface{})
	// 	oldnets2 = append(oldnets2, m1)
	// }
	// var oldnets3 []map[string]string
	// for _, m := range *oldnet {
	// 	for k, v := range m {
	// 		m2 := make(map[string]string)
	// 		m2[k] = v.(string)
	// 		oldnets3 = append(oldnets3, m2)
	// 	}
	// }
	var oldnets3 = []map[string]string{
		*oldnet,
	}

	switch ip_assignment {
	case "static":
		text := baseos.NetworkStatic(oldnets3, ipv4_addr, ipv4_netmask, ipv4_gateway, ipv4_dns1, network.DNS2)
		_, Err, err := agent.UpdateConfigFile(global.NetWorkPath, nic_name, text)
		if len(Err) != 0 || err != nil {
			response.Fail(c, nil, Err)
			return
		}
		Err, err = agent.RestartNetWork(nic_name)
		if len(Err) != 0 || err != nil {
			response.Fail(c, nil, Err)
			return
		}
		response.Success(c, nil, "网络配置更新成功")

	case "dhcp":
		text := baseos.NetworkDHCP(oldnets3)
		_, Err, err := agent.UpdateConfigFile(global.NetWorkPath, nic_name, text)
		if len(Err) != 0 || err != nil {
			response.Fail(c, nil, Err)
			return
		}
		Err, err = agent.RestartNetWork(nic_name)
		if len(Err) != 0 || err != nil {
			response.Fail(c, nil, Err)
			return
		}
		response.Success(c, nil, "网络配置更新成功")

	default:
		response.Fail(c, nil, "请重新检查ip分配方式")
	}
}
