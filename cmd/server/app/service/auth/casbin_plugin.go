/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Thu Apr 10 16:18:53 2025 +0800
 */
package auth

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// 添加插件服务权限到列表中
func AddPluginServicePermission(role string, permissions []common.Permission, serviceName string) error {
	//TODO；先添加到列表中可以展示
	for _, v := range permissions {
		ok, err := addPolicy(role, v.Resource, v.Operate, serviceName)
		if err != nil {
			logger.Error("init plugin-admin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("plugin-admin %s permission already exists: %s", v.Operate, v.Resource)
		}
	}
	return nil
}

// 删除插件服务权限
func DeletePluginServicePermission(permissions []common.Permission, serviceName string) error {
	for _, v := range permissions {
		ok, err := G_Enfocer.RemoveFilteredPolicy(1, v.Resource, v.Operate, serviceName)
		if err != nil {
			logger.Error("delete plugin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("delete plugin %s permission failed: %s", v.Operate, v.Resource)
		}
	}
	return nil
}
