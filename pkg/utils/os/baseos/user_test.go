/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Thu Apr 13 14:50:44 2023 +0800
 */
package baseos

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUserInfo(t *testing.T) {
	var osobj BaseOS
	tmp := osobj.GetCurrentUserInfo()
	assert.NotNil(t, tmp.GroupName)
	assert.NotNil(t, tmp.Groupid)
	assert.NotNil(t, tmp.HomeDir)
	assert.NotNil(t, tmp.Userid)
	assert.NotNil(t, tmp.Username)
}

func TestGetAllUserInfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetAllUserInfo()
	assert.Nil(t, err)
	for _, v := range tmp {
		assert.NotNil(t, v.Description)
		assert.NotNil(t, v.GroupId)
		assert.NotNil(t, v.HomeDir)
		assert.NotNil(t, v.ShellType)
		assert.NotNil(t, v.UserId)
		assert.NotNil(t, v.Username)
	}
}

func TestConfigUser(t *testing.T) {
	var osobj BaseOS

	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	username := base64.URLEncoding.EncodeToString(randomBytes)
	password := "china666*"
	permission := "444"
	file := "testfile"
	fileabs := "/home/" + username + "/" + file

	t.Run("test AddLinuxUser", func(t *testing.T) {
		err := osobj.AddLinuxUser(username, password)
		assert.Nil(t, err)
	})

	t.Run("test ChangePermission", func(t *testing.T) {
		exitc, stdo, stde, err := utils.RunCommand("touch /home/" + username + "/" + file)
		assert.Equal(t, 0, exitc)
		assert.Equal(t, "", strings.Replace(stdo, "\n", "", -1))
		assert.Equal(t, "", strings.Replace(stde, "\n", "", -1))
		assert.Nil(t, err)

		_, err = osobj.ChangePermission(permission, fileabs)
		assert.Nil(t, err)

		exitc2, stdo2, stde2, err2 := utils.RunCommand("ls -l " + fileabs)
		assert.Equal(t, 0, exitc2)
		assert.NotNil(t, stdo2)
		assert.Equal(t, "", strings.Replace(stde2, "\n", "", -1))
		assert.Nil(t, err2)
		assert.Equal(t, "-r--r--r--", strings.Replace(strings.Split(stdo2, " ")[0], ".", "", -1))
	})

	t.Run("test ChangeFileOwner", func(t *testing.T) {
		_, err := osobj.ChangeFileOwner(username, fileabs)
		assert.Nil(t, err)

		exitc, stdo, stde, err2 := utils.RunCommand("ls -l " + fileabs)
		assert.Equal(t, 0, exitc)
		assert.NotNil(t, stdo)
		assert.Equal(t, "", strings.Replace(stde, "\n", "", -1))
		assert.Nil(t, err2)
		assert.Equal(t, username, strings.Split(stdo, " ")[2])
	})

	t.Run("test DelUser", func(t *testing.T) {
		_, err := osobj.DelUser(username)
		assert.Nil(t, err)

		exitc, stdo, stde, err2 := utils.RunCommand("cat /etc/passwd | cut -d : -f 1 | grep \"" + username + "\"")
		assert.Equal(t, 1, exitc)
		assert.Equal(t, "", strings.Replace(stde, "\n", "", -1))
		assert.Nil(t, err2)
		assert.Equal(t, "", strings.Replace(stdo, "\n", "", -1))
	})
}
