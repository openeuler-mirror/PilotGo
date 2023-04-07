package baseos

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllRpm(t *testing.T) {
	var osobj BaseOS
	assert.NotNil(t, osobj.GetAllRpm())
}

func TestGetRpmSource(t *testing.T) {
	var osobj BaseOS
	rpm := "time"
	tmp, err := osobj.GetRpmSource(rpm)
	assert.Nil(t, err)
	assert.Equal(t, rpm, strings.Split(tmp[0].Name, "-")[0])
}

func TestGetRpmInfo(t *testing.T) {
	var osobj BaseOS
	rpm := "time"
	tmp, err1, err2 := osobj.GetRpmInfo(rpm)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, rpm, tmp.Name)
}

func TestInstallAndRemoveRpm(t *testing.T) {
	var osobj BaseOS
	rpm := "bind"
	t.Run("test install rpm", func(t *testing.T) {
		assert.Nil(t, osobj.InstallRpm(rpm))
		tmp, err1, err2 := osobj.GetRpmInfo(rpm)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		assert.Equal(t, rpm, tmp.Name)
	})

	t.Run("test remove rpm", func(t *testing.T) {
		assert.Nil(t, osobj.RemoveRpm(rpm))
		_, err1, err2 := osobj.GetRpmInfo(rpm)
		assert.NotNil(t, err1)
		assert.NotNil(t, err2)
	})
}
