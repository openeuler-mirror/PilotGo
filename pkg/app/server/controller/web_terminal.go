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
 * Date: 2022-04-20 16:48:55
 * LastEditTime: 2022-04-20 17:48:55
 * Description: web socket连接控制
 ******************************************************************************/

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/webSocket"
)

func ShellWs(c *gin.Context) {
	msg := c.DefaultQuery("msg", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := dao.Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	// 后端获取到前端传来的主机信息,以此建立ssh客户端
	sshClient, err := webSocket.DecodedMsgToSSHClient(msg)
	if err != nil {
		c.Error(err)
		return
	}
	if sshClient.IpAddress == "" || sshClient.Password == "" {
		c.Error(&dao.ApiError{Message: "IP地址或密码不能为空", Code: 400})
		return
	}
	// 升级协议并获得socket连接
	conn, err := webSocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	// 生成ssh socket客户端，建立session、client、channel
	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(1, []byte("验证失败，请检查用户名或密码..."))
		conn.Close()
		return
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
}
