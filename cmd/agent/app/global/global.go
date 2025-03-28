/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package global

import (
	"os"

	"gitee.com/openeuler/PilotGo/cmd/agent/app/network"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/message/protocol"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

var AgentVersion = "v0.0.1"

type ConfigMessage struct {
	ConfigName    string
	ConfigContent string
	ConfigChange  string
	Machine_uuid  string
}

// 配置文件的监听器
func Configfsnotify(ConMess ConfigMessage, client *network.SocketClient) error {
	//创建一个监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("NewWatcher failed: %s", err)
		return err
	}
	defer watcher.Close()
	done := make(chan bool)
	err = watcher.Add(ConMess.ConfigName)
	if err != nil {
		logger.Error("Add failed:%s", err)
	}
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Debug("在文件 %s 上进行 : %s", event.Name, event.Op)
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Write == fsnotify.Write {
					ConMess.ConfigChange = event.Op.String()
					ConMess.ConfigContent, err = utils.FileReadString(ConMess.ConfigName)
					if err != nil {
						logger.Debug("error: %s", err)
					}
					msg := &protocol.Message{
						UUID:   uuid.New().String(),
						Type:   protocol.ConfigFileMonitor,
						Status: 0,
						Data:   ConMess,
					}
					if err := client.Send(msg); err != nil {
						logger.Debug("send message failed, error: %s", err)
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					_, err := os.Stat(ConMess.ConfigName)
					if err == nil {
						err = watcher.Add(ConMess.ConfigName)
					}
				}
			case err, ok := <-watcher.Errors:
				logger.Debug("error: %s", err)
				if !ok {
					return
				}
			}
		}
	}()
	<-done
	return nil
}
