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
 * LastEditTime: 2023-09-01 15:36:32
 * Description: casbin服务
 ******************************************************************************/
package auth

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
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

func Casbin(conf *sconfig.MysqlDBInfo) {
	text := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[matchers]
	m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act

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

func addPolicy(role, resource, action string) (bool, error) {
	return G_Enfocer.AddPolicy(role, resource, action)
}

func initAdminPolicy() {
	G_Enfocer.AddRoleForUser("admin", "admin")

	ok, err := addPolicy("admin", "/*", "get")
	if err != nil {
		logger.Error("init admin policy failed:%s", err)
	}
	if !ok {
		logger.Info("admin policy already exists")
	}
}

func CheckAuth(user, resource, action string) (bool, error) {
	return G_Enfocer.Enforce(user, resource, action)
}
