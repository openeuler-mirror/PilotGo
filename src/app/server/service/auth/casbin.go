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
 * Date: 2021-07-04 09:08:08
 * LastEditTime: 2023-09-04 11:12:43
 * Description: casbin服务
 ******************************************************************************/
package auth

import (
	"errors"
	"fmt"
	"sync"

	sconfig "gitee.com/openeuler/PilotGo/app/server/config"
	suser "gitee.com/openeuler/PilotGo/app/server/service/user"
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

var (
	once sync.Once
)

const (
	DomainPilotGo = "PilotGo-server"
)

func Casbin(conf *sconfig.MysqlDBInfo) {
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

	once.Do(func() {
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
	})
}

// deprecated
func AllPolicy() (interface{}, int) {
	casbin := make([]map[string]interface{}, 0)
	list := G_Enfocer.GetPolicy()
	for _, vlist := range list {
		policy := make(map[string]interface{})
		policy["role"] = vlist[0]
		policy["url"] = vlist[1]
		policy["method"] = vlist[2]
		casbin = append(casbin, policy)
	}
	total := len(casbin)
	return casbin, total
}

// deprecated
func PolicyRemove(rule CasbinRule) bool {
	ok, err := G_Enfocer.RemovePolicy(rule.RoleType, rule.Url, rule.Method)
	if err == nil {
		if !ok {
			return false
		} else {
			return true
		}
	}
	logger.Error("remove policy error:%s", err)
	return false
}

// deprecated
func PolicyAdd(rule CasbinRule) bool {
	ok, err := G_Enfocer.AddPolicy(rule.RoleType, rule.Url, rule.Method)
	if err == nil {
		if !ok {
			return false
		} else {
			return true
		}
	}
	logger.Error("add policy error:%s", err)
	return false
}

func addPolicy(role, resource, action, domain string) (bool, error) {
	return G_Enfocer.AddPolicy(role, resource, action, domain)
}

var (
	PermissionList = []string{
		"rpm_install",
		"rpm_uninstall",
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
		"config_install",
		"dept_change",
	}

	MenuList = []string{
		"overview",
		"cluster",
		"batch",
		"usermanager",
		"rolemanager",
		"config",
		"log",
		"plugin",
	}
)

// 添加信息到列表中
func AddPluginPermission(permissions []common.Permission, uuid string) error {
	//TODO；先添加到列表中可以展示，再通过修改权限进行调整
	// 添加admin的插件权限
	for _, v := range permissions {
		ok, err := addPolicy("admin", v.Resource, v.Operate, uuid)
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
		fmt.Println(v)
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

// 获取指定用户的buttion权限和menu权限
func GetPermissionsOfUser(user string) ([]string, []string, error) {
	ps, err := G_Enfocer.GetImplicitPermissionsForUser(user)
	if err != nil {
		return nil, nil, err
	}
	// fmt.Printf("user permissions: %v\n", ps)

	pbutton := []string{}
	pmenu := []string{}
	for _, ss := range ps {
		// user := ss[0]
		resource := ss[1]
		ptype := ss[2]

		switch ptype {
		case "button":
			pbutton = append(pbutton, resource)
		case "menu":
			pmenu = append(pmenu, resource)
		case "empty":
			// TODO: 历史兼容性保留的空权限
			continue
		default:
			logger.Warn("unknown permission type: %s", ptype)
		}
	}
	return pbutton, pmenu, nil
}

func UpdateRolePermissions(role string, buttons, menus []string) error {
	if role == "admin" {
		return errors.New("admin角色权限不可修改")
	}

	if _, err := G_Enfocer.DeleteRole(role); err != nil {
		return err
	}

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
