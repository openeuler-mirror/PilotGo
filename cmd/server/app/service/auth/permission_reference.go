/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Tue Aug 12 16:18:53 2025 +0800
 */
package auth

// permission reference
var PermissionListMap *PermissionMap

type PermissionMap struct {
	data map[string]map[string][]string
}

func NewPermissionMap() {
	PermissionListMap = &PermissionMap{
		data: make(map[string]map[string][]string),
	}
}

func (pm *PermissionMap) Set(domain, action, resource string) {
	if pm.data[domain] == nil {
		pm.data[domain] = make(map[string][]string)
	}
	resources := pm.data[domain][action]
	for _, a := range resources {
		if a == resource {
			continue
		}
	}
	pm.data[domain][action] = append(resources, resource)
}

func (pm *PermissionMap) Get(domain, action string) ([]string, bool) {
	if innerMap, ok := pm.data[domain]; ok {
		val, ok := innerMap[action]
		return val, ok
	}
	return []string{}, false
}

func (pm *PermissionMap) Delete(domain, resource, action string) {
	if innerMap, ok := pm.data[domain]; ok {
		if resources, ok := innerMap[action]; ok {
			// 删除指定 resource
			newResources := []string{}
			for _, r := range resources {
				if r != resource {
					newResources = append(newResources, r)
				}
			}
			if len(newResources) > 0 {
				innerMap[action] = newResources
			} else {
				delete(innerMap, action)
			}
		}
		if len(innerMap) == 0 {
			delete(pm.data, domain)
		}
	}
}

func (pm *PermissionMap) DeleteDomain(domain string) {
	delete(pm.data, domain)
}

func (pm *PermissionMap) FindByInnerKey(resource string) (string, string) {
	for domain, innerMap := range pm.data {
		for action, resources := range innerMap {
			for _, r := range resources {
				if r == resource {
					return domain, action
				}
			}
		}
	}
	return "", ""
}

func (pm *PermissionMap) GetAll() map[string]map[string][]string {
	return pm.data
}

const (
	DomainPilotGo = "PilotGo-server"
)

var (
	PermissionList = []string{
		"rpm_install",
		"rpm_uninstall",
		"batch_create",
		"batch_update",
		"batch_delete",
		"user_add",
		"user_import",
		"user_edit",
		"user_reset",
		"user_del",
		"role_add",
		"role_update",
		"role_delete",
		"role_modify",
		"dept_change",
		"dept_add",
		"dept_delete",
		"dept_update",
		"machine_delete",
		"run_script",
		"update_script_blacklist",
		"plugin_operate",
	}

	MenuList = []string{
		"overview",
		"cluster",
		"usermanager",
		"rolemanager",
		"audit",
		"plugin",
		"terminal",
		"script",
	}
)
