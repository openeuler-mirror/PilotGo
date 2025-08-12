/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package plugin

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type PluginPermission struct {
	ServiceName string
	Permissions []common.Permission
}

// 从db获取某个角色的所有插件权限
func GetRolePluginPermission(role string) map[string][]string {
	pluginsPermission := make(map[string][]string)

	pluginServices := global.GW.GetAllServices()
	for _, ps := range pluginServices {
		menu := []string{}
		buttons := []string{}
		policys := auth.GetFilteredPolicy(role, "", "button", ps["serviceName"].(string))
		for _, v := range policys {
			buttons = append(buttons, v[1])
		}

		policys = auth.GetFilteredPolicy(role, "", "menu", ps["serviceName"].(string))
		for _, v := range policys {
			menu = append(menu, v[1])
		}

		pluginsPermission["menu"] = append(pluginsPermission["menu"], menu...)
		pluginsPermission["button"] = append(pluginsPermission["button"], buttons...)
	}
	return pluginsPermission
}

// 更新插件角色权限
func UpdatePluginPermissions(role string, PluginPermissions []PluginPermission) error {
	for _, p := range PluginPermissions {
		err := auth.AddPluginServicePermission(role, p.Permissions, p.ServiceName)
		if err != nil {
			logger.Error("add role:%s buttion policy failed:%s", role, err)
			return err
		}
	}
	return nil
}
