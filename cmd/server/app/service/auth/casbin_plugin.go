/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Apr 10 16:18:53 2025 +0800
 */
package auth

import (
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/mitchellh/mapstructure"
)

// 添加插件服务权限到列表中
func AddPluginServicePermission(role string, permissions []common.Permission, serviceName string) error {
	//TODO；先添加到列表中可以展示
	for _, v := range permissions {
		PermissionListMap.Set(serviceName, v.Operate, v.Resource)
		ok, err := addPolicy(role, v.Resource, v.Operate, serviceName)
		if err != nil {
			logger.Error("init %s policy failed:%s", serviceName, err)
		}
		if !ok {
			logger.Debug("%s %s permission already exists: %s", serviceName, v.Operate, v.Resource)
		}
	}
	return nil
}

// 删除插件服务权限
func DeletePluginServicePermission(permissions []common.Permission, serviceName string) error {
	PermissionListMap.DeleteDomain(serviceName)
	for _, v := range permissions {
		ok, err := G_Enfocer.RemoveFilteredPolicy(1, v.Resource, v.Operate, serviceName)
		if err != nil {
			logger.Error("delete %s policy failed:%s", serviceName, err)
		}
		if !ok {
			logger.Debug("delete %s %s permission failed: %s", serviceName, v.Operate, v.Resource)
		}
	}
	return nil
}

func initPluginServicePermission() {
	services := global.GW.GetAllServices()
	for _, service := range services {
		var permissions []common.Permission
		if perms, ok := service["permissions"]; ok {
			if err := mapstructure.Decode(perms, &permissions); err != nil {
				logger.Error("decode permission policy failed:%s", err)
			}

			if err := AddPluginServicePermission("admin", permissions, service["serviceName"].(string)); err != nil {
				logger.Error("Failed to add permissions for service %s: %v", service["serviceName"].(string), err)
			}
		}
	}
}
