/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package plugin

import (
	"fmt"
	"time"

	eventSDK "gitee.com/openeuler/PilotGo-plugins/event/sdk"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"

	commonSDK "gitee.com/openeuler/PilotGo/sdk/common"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CheckPluginHeartbeats(stopCh <-chan struct{}) {
	go wait.Until(func() {
		checkAndRebind()
	}, client.HeartbeatInterval, stopCh)
}

func checkAndRebind() {
	plugins, err := GetPlugins()
	if err != nil {
		logger.Error("get plugins failed:%v", err.Error())
	}
	for _, p := range plugins {
		key := client.HeartbeatKey + p.Url
		plugin_status, err := redismanager.Get(key, &client.PluginStatus{})
		if err != nil {
			logger.Error("Error getting %v last heartbeat: %v", p.Url, err)
			continue
		}

		if !plugin_status.(*client.PluginStatus).Connected || time.Since(plugin_status.(*client.PluginStatus).LastConnect) > client.HeartbeatInterval+1*time.Second {
			err := Handshake(p.Url, p)
			if err != nil {
				if time.Since(plugin_status.(*client.PluginStatus).LastConnect) <= client.HeartbeatInterval*2+1*time.Second {
					_addr := p.Url
					if p.Url == "localhost" || p.Url == "127.0.0.1" {
						node, err := dao.MachineInfoByUUID(p.UUID)
						if err == nil {
							_addr = node.IP
						} else {
							logger.Error("fail to get machineinfo by uuid: %s", err.Error())
						}
					}
					global.SendRemindMsg(
						global.PluginSendMsg,
						fmt.Sprintf("%s 插件离线 %s", p.Name, _addr),
					)
				}

				logger.Error("rebind plugin and pilotgo server failed:%v", err.Error())
				value := client.PluginStatus{
					Connected:   false,
					LastConnect: plugin_status.(*client.PluginStatus).LastConnect,
				}
				redismanager.Set(key, value)

				// 缓存，发布“插件离线”事件
				offlineKey := "offline:" + p.UUID
				offlineValue := struct {
					OfflineTime time.Time
				}{
					OfflineTime: time.Now(),
				}
				ok, err := redismanager.SetNX(offlineKey, offlineValue)
				if ok && err == nil {
					msgData := commonSDK.MessageData{
						MsgType:     eventSDK.MsgPluginOffline,
						MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgPluginOffline),
						TimeStamp:   time.Now(),
						Data: eventSDK.MDPluginChange{
							PluginName:  p.Name,
							Version:     p.Version,
							Url:         p.Url,
							Description: p.Description,
							Status:      false,
						},
					}
					msgDataString, err := msgData.ToMessageDataString()
					if err != nil {
						logger.Error("event message data marshal failed:%v", err.Error())
					}
					ms := commonSDK.EventMessage{
						MessageType: eventSDK.MsgPluginOffline,
						MessageData: msgDataString,
					}
					PublishEvent(ms)
				}
			} else {
				if time.Since(plugin_status.(*client.PluginStatus).LastConnect) > client.HeartbeatInterval*2+1*time.Second {
					_addr := p.Url
					if p.Url == "localhost" || p.Url == "127.0.0.1" {
						node, err := dao.MachineInfoByUUID(p.UUID)
						if err == nil {
							_addr = node.IP
						} else {
							logger.Error("fail to get machineinfo by uuid: %s", err.Error())
						}
					}
					global.SendRemindMsg(
						global.PluginSendMsg,
						fmt.Sprintf("%s 插件上线 %s", p.Name, _addr),
					)
				}

				value := client.PluginStatus{
					Connected:   true,
					LastConnect: time.Now(),
				}
				redismanager.Set(key, value)

				//删除缓存，发布“插件上线”事件
				redismanager.Delete("offline:" + p.UUID)
				msgData := commonSDK.MessageData{
					MsgType:     eventSDK.MsgPluginOnline,
					MessageType: eventSDK.GetMessageTypeString(eventSDK.MsgPluginOnline),
					TimeStamp:   time.Now(),
					Data: eventSDK.MDPluginChange{
						PluginName:  p.Name,
						Version:     p.Version,
						Url:         p.Url,
						Description: p.Description,
						Status:      true,
					},
				}
				msgDataString, err := msgData.ToMessageDataString()
				if err != nil {
					logger.Error("event message data marshal failed:%v", err.Error())
				}
				ms := commonSDK.EventMessage{
					MessageType: eventSDK.MsgPluginOnline,
					MessageData: msgDataString,
				}
				PublishEvent(ms)
			}
		}
	}
}
