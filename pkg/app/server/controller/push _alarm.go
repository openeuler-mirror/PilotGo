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
 * LastEditTime: 2022-06-24 16:48:55
 * Description: 通过web socket方式推送告警
 ******************************************************************************/

package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/agentmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

type ConnClient struct {
	Conn *websocket.Conn
}

var clients = make(map[int]*ConnClient)
var i int = 0
var lock sync.Mutex

func PushAlarmHandler(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	lock.Lock()
	i++
	key := i
	lock.Unlock()
	client := &ConnClient{Conn: conn}
	clients[key] = client
	go func() {
		for {
			data := <-agentmanager.WARN_MSG
			for _, cli := range clients {
				cli.Conn.WriteMessage(websocket.TextMessage, []byte(data.(string)))
			}
		}
	}()
}
