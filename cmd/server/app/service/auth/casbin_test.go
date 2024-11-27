/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

// func TestMain(m *testing.M) {
// 	err := sconfig.Init("D:\\tmp\\PilotGo-projects\\PilotGo\\config_server.yaml")
// 	if err != nil {
// 		fmt.Println("failed to load configure, exit..", err)
// 		os.Exit(-1)
// 	}

// 	// 鉴权模块初始化
// 	Casbin(&sconfig.Config().MysqlDBinfo)

// 	m.Run()
// }
