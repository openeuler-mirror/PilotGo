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
 * Date: 2022-06-24 10:48:55
 * LastEditTime: 2022-06-24 16:48:55
 * Description: 通过web socket方式推送告警
 ******************************************************************************/

package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/network/websocket"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
)

func PushAlarmHandler(c *gin.Context) {
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Debug("webSocket 建立连接: %s", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := websocket.NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.Read()
	go client.Write()

	// 用户连接事件
	websocket.CliManager.Register <- client
}
