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
 * LastEditTime: 2023-09-01 16:22:14
 * Description: casbin服务
 ******************************************************************************/
package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	sconfig "gitee.com/openeuler/PilotGo/app/server/config"
)

func TestGetRoles(t *testing.T) {
	roles := GetAllRoles()
	assert.NotNil(t, roles)

	fmt.Printf("roles: %v\n", roles)
}

func TestGetUserRoles(t *testing.T) {
	roles, err := GetUserRoles("admin")
	assert.NoError(t, err)

	fmt.Printf("user roles: %v\n", roles)
}

func TestGetAllPolicy(t *testing.T) {
	policies := GetAllPolicies()

	fmt.Printf("policies: %v\n", policies)
}

func TestMain(m *testing.M) {
	err := sconfig.Init("D:\\tmp\\PilotGo-projects\\PilotGo\\config_server.yaml")
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	// 鉴权模块初始化
	Casbin(&sconfig.Config().MysqlDBinfo)

	m.Run()
}

func TestGetPermissionsOfUser(t *testing.T) {
	pbutton, pmenu, err := GetPermissionsOfUser("admin")
	assert.NoError(t, err)
	fmt.Printf("pbutton: %v, pmenu: %v\n", pbutton, pmenu)
}
