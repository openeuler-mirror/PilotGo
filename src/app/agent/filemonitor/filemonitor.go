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
 * Date: 2022-07-05 13:03:16
 * LastEditTime: 2023-02-21 19:02:55
 * Description: file monitor init
 ******************************************************************************/
package filemonitor

import (
	"fmt"
	"path"
	"strings"

	"gitee.com/openeuler/PilotGo/app/agent/localstorage"
	"gitee.com/openeuler/PilotGo/app/agent/network"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils/message/protocol"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

// repo路径
const RepoPath = "/etc/yum.repos.d"

// 网络配置
const NetWorkPath = "/etc/sysconfig/network-scripts"

var RESP_MSG = make(chan interface{})

func FileMonitorInit() error {
	// get agent uuid
	uuid := localstorage.AgentUUID()

	// 1、NewWatcher 初始化一个 watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 2、使用 watcher 的 Add 方法增加需要监听的文件或目录到监听队列中
	err = watcher.Add(RepoPath)
	if err != nil {
		logger.Debug("failed to monitor repo")
	}
	logger.Info("start to monitor repo")

	err = watcher.Add(NetWorkPath)
	if err != nil {
		logger.Debug("failed to monitor network")
	}
	logger.Info("start to monitor network")

	//3、创建新的 goroutine，等待管道中的事件或错误
	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}
				fileExt := path.Ext(e.Name)
				if strings.Contains(fileExt, ".sw") || strings.Contains(fileExt, "~") || strings.Contains(e.Name, "~") {
					continue
				}

				if e.Op&fsnotify.Write == fsnotify.Write {
					RESP_MSG <- fmt.Sprintf("机器 %s 上的文件已被修改 : %s", uuid, e.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("error: %s", err)
			}
		}
	}()
	<-done
	return nil
}

func FileMonitor(client *network.SocketClient) {
	for data := range RESP_MSG {
		if data == nil {
			continue
		}

		msg := &protocol.Message{
			UUID:   uuid.New().String(),
			Type:   protocol.FileMonitor,
			Status: 0,
			Data:   data,
		}

		if err := client.Send(msg); err != nil {
			logger.Debug("send message failed, error: %s", err)
		}

	}
}
