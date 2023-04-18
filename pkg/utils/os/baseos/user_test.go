package baseos

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
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
	tmp := osobj.GetAllUserInfo()
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
	password := "123456"
	permission := "444"
	file := "testfile"
	fileabs := "/home/" + username + "/" + file

	t.Run("test AddLinuxUser", func(t *testing.T) {
		err := osobj.AddLinuxUser(username, password)
		assert.Nil(t, err)
	})

	t.Run("test ChangePermission", func(t *testing.T) {
		_, err := utils.RunCommand("touch /home/" + username + "/" + file)
		assert.Nil(t, err)

		_, err = osobj.ChangePermission(permission, fileabs)
		assert.Nil(t, err)

		output, err := utils.RunCommand("ls -l " + fileabs)
		assert.Nil(t, err)
		assert.Equal(t, "-r--r--r--", strings.Replace(strings.Split(output, " ")[0], ".", "", -1))
	})

	t.Run("test ChangeFileOwner", func(t *testing.T) {
		_, err := osobj.ChangeFileOwner(username, fileabs)
		assert.Nil(t, err)

		output, err := utils.RunCommand("ls -l " + fileabs)
		assert.Nil(t, err)
		assert.Equal(t, username, strings.Split(output, " ")[2])
	})

	t.Run("test DelUser", func(t *testing.T) {
		_, err := osobj.DelUser(username)
		assert.Nil(t, err)

		output, _ := utils.RunCommand("cat /etc/passwd | cut -d : -f 1 | grep \"" + username + "\"")
		assert.Nil(t, err)
		assert.Equal(t, "", strings.Replace(output, "\n", "", -1))
	})
}
