/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package auth

import (
	"errors"
	"fmt"

	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	suser "gitee.com/openeuler/PilotGo/cmd/server/app/service/user"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var G_Enfocer *casbin.Enforcer

type CasbinRule struct {
	PType    string `json:"ptype"`
	RoleType string `json:"role"`
	Url      string `json:"url"`
	Method   string `json:"method"`
}

const (
	DomainPilotGo = "PilotGo-server"
)

func Casbin(conf *options.MysqlDBInfo) {
	text := `
	[request_definition]
	r = sub, obj, act, domain

	[policy_definition]
	p = sub, obj, act, domain

	[role_definition]
	g = _, _

	[matchers]
	m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)&& regexMatch(r.domain, p.domain)

	[policy_effect]
	e = some(where (p.eft == allow))
	`

	// m := casbinmodel.Model{}
	m, err := casbinmodel.NewModelFromString(text)
	if err != nil {
		logger.Error("casbin model create failed: %s", err)
		return
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		conf.UserName,
		conf.Password,
		conf.HostName,
		conf.Port)
	a, err := gormadapter.NewAdapter("mysql", url, conf.DataBase)
	if err != nil {
		logger.Error("casbin adapter create failed: %s", err)
		return
	}
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		logger.Error("casbin enforcer create failed: %s", err)
		return
	}
	if err := enforcer.LoadPolicy(); err != nil {
		logger.Error("casbin load Policy failed: %s", err.Error())
	}

	G_Enfocer = enforcer

	// TODO:
	initAdminPolicy()
}

func addPolicy(role, resource, action, domain string) (bool, error) {
	return G_Enfocer.AddPolicy(role, resource, action, domain)
}

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
	}

	MenuList = []string{
		"overview",
		"cluster",
		"batch",
		"usermanager",
		"rolemanager",
		"audit",
		"plugin",
		"terminal",
		"script",
	}
)

// 添加插件权限到列表中
func AddPluginPermission(role string, permissions []common.Permission, uuid string) error {
	//TODO；先添加到列表中可以展示
	for _, v := range permissions {
		ok, err := addPolicy(role, v.Resource, v.Operate, uuid)
		if err != nil {
			logger.Error("init plugin-admin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("plugin-admin %s permission already exists: %s", v.Operate, v.Resource)
		}
	}
	return nil
}

// 删除插件权限
func DeletePluginPermission(permissions []common.Permission, uuid string) error {
	for _, v := range permissions {
		ok, err := G_Enfocer.RemoveFilteredPolicy(1, v.Resource, v.Operate, uuid)
		if err != nil {
			logger.Error("delete plugin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("delete plugin %s permission failed: %s", v.Operate, v.Resource)
		}
	}
	return nil
}

func initAdminPolicy() {
	G_Enfocer.AddRoleForUser("admin", "admin")

	for _, p := range PermissionList {
		ok, err := addPolicy("admin", p, "button", DomainPilotGo)
		if err != nil {
			logger.Error("init admin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("admin button permission already exists: %s", p)
		}
	}

	for _, m := range MenuList {
		ok, err := addPolicy("admin", m, "menu", DomainPilotGo)
		if err != nil {
			logger.Error("init admin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("admin menu access permission already exists: %s", m)
		}
	}

	// test
	{
		ok, err := addPolicy("admin", "plugins", "get", DomainPilotGo)
		if err != nil {
			logger.Error("init admin policy failed:%s", err)
		}
		if !ok {
			logger.Debug("admin test permission already exists")
		}
	}
}

func CheckAuth(user, resource, action, domain string) (bool, error) {
	roles, err := suser.GetUserRoles(user)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		ok, err := G_Enfocer.Enforce(role, resource, action, domain)
		logger.Debug("check %s auth: %s %s %s %s, result: %t", user, role, resource, action, domain, ok)

		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func GetAllRoles() []string {
	return G_Enfocer.GetAllRoles()
}

func GetUserRoles(user string) ([]string, error) {
	// TODO:
	// return G_Enfocer.GetImplicitRolesForUser(user)
	return G_Enfocer.GetRolesForUser(user)
}

func AddRole(role string) error {
	// TODO: 为了兼容历史版本创建空role，创建一个无用的权限
	addPolicy(role, "empty", "emply", DomainPilotGo)
	return nil
}

type Policy struct {
	Role     string
	Resource string
	Action   string
	Domain   string
}

func GetAllPolicies() []Policy {
	policies := G_Enfocer.GetPolicy()

	result := []Policy{}
	for _, p := range policies {
		result = append(result, Policy{
			Role:     p[0],
			Resource: p[1],
			Action:   p[2],
		})
	}
	return result
}

func GetFilteredPolicy(role, object, action, domain string) [][]string {
	return G_Enfocer.GetFilteredPolicy(0, role, object, action, domain)
}

func UpdateRolePermissions(role string, buttons, menus []string) error {
	for _, p := range buttons {
		_, err := addPolicy(role, p, "button", DomainPilotGo)
		if err != nil {
			logger.Error("add role:%s buttion policy failed:%s", role, err)
			return err
		}
	}

	for _, m := range menus {
		_, err := addPolicy(role, m, "menu", DomainPilotGo)
		if err != nil {
			logger.Error("add role:%s menu policy failed:%s", role, err)
			return err
		}
	}
	return nil
}

func DeleteRole(role string) error {
	if role == "admin" {
		return errors.New("admin角色不可删除")
	}
	if _, err := G_Enfocer.DeleteRole(role); err != nil {
		return err
	}
	return nil
}
