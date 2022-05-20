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
 * LastEditTime: 2022-04-13 01:51:51
 * Description: provide agent network manager functions.
 ******************************************************************************/
package agentcontroller

import (
	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
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
