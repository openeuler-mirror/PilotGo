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
 * LastEditTime: 2022-07-04 09:25:41
 * Description: casbin服务
 ******************************************************************************/
package service

import (
	"sync"

	"github.com/casbin/casbin"
	casbinmodel "github.com/casbin/casbin/model"
	gormadapter "github.com/casbin/gorm-adapter"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/dbmanager/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

func Casbin() *casbin.Enforcer {

	text := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`

	m := casbinmodel.Model{}
	m.LoadModelFromText(text)

	once.Do(func() {
		a := gormadapter.NewAdapter("mysql", mysqlmanager.Url, true)
		enforcer = casbin.NewEnforcer(m, a)
	})
	enforcer.LoadPolicy()
	return enforcer
}

func AllPolicy() (interface{}, int) {
	casbin := make([]map[string]interface{}, 0)
	list := global.PILOTGO_E.GetPolicy()
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

func PolicyRemove(rule model.CasbinRule) bool {
	if ok := global.PILOTGO_E.RemovePolicy(rule.RoleType, rule.Url, rule.Method); !ok {
		return false
	} else {
		return true
	}
}

func PolicyAdd(rule model.CasbinRule) bool {
	if ok := global.PILOTGO_E.AddPolicy(rule.RoleType, rule.Url, rule.Method); !ok {
		return false
	} else {
		return true
	}
}
