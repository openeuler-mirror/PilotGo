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
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

type CasbinRule struct {
	PType    string `json:"ptype"`
	RoleType string `json:"role"`
	Url      string `json:"url"`
	Method   string `json:"method"`
}

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

func Casbin(conf *sconfig.MysqlDBInfo) *casbin.Enforcer {

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

	once.Do(func() {
		// m := casbinmodel.Model{}
		m, err := casbinmodel.NewModelFromString(text)
		if err != nil {
			logger.Fatal("casbin model create failed: %s", err)
		}

		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			conf.UserName,
			conf.Password,
			conf.HostName,
			conf.Port)
		a, err := gormadapter.NewAdapter("mysql", url, conf.DataBase)
		if err != nil {
			logger.Fatal("casbin adapter create failed: %s", err)
		}
		enforcer, err = casbin.NewEnforcer(m, a)
		if err != nil {
			logger.Fatal("casbin enforcer create failed: %s", err)
		}
		enforcer.LoadPolicy()
	})
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

func PolicyRemove(rule CasbinRule) bool {
	ok, err := global.PILOTGO_E.RemovePolicy(rule.RoleType, rule.Url, rule.Method)
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
	ok, err := global.PILOTGO_E.AddPolicy(rule.RoleType, rule.Url, rule.Method)
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
