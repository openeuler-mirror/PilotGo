package baseos

import (
	"strings"
	"testing"

	"gitee.com/openeuler/PilotGo/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAllRpm(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetAllRpm()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
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
	tmp, err := osobj.GetRpmInfo(rpm)
	assert.Nil(t, err)
	assert.Equal(t, rpm, tmp.Name)
}

func TestInstallAndRemoveRpm(t *testing.T) {
	var osobj BaseOS
	rpm := "bind"

	exitc, stdo, stde, err := utils.RunCommand("rpm -qi " + rpm)
	if exitc == 0 && len(stdo) > 0 && stde == "" && err == nil {
		t.Run("test remove rpm", func(t *testing.T) {
			assert.Nil(t, osobj.RemoveRpm(rpm))
			_, err := osobj.GetRpmInfo(rpm)
			assert.NotNil(t, err)
		})

		t.Run("test install rpm", func(t *testing.T) {
			assert.Nil(t, osobj.InstallRpm(rpm))
			tmp, err := osobj.GetRpmInfo(rpm)
			assert.Nil(t, err)
			assert.Equal(t, rpm, tmp.Name)
		})
	} else if exitc == 1 && strings.Replace(stdo, "\n", "", -1) == "package bind is not installed" && stde == "" && err == nil {
		t.Run("test install rpm", func(t *testing.T) {
			assert.Nil(t, osobj.InstallRpm(rpm))
			tmp, err := osobj.GetRpmInfo(rpm)
			assert.Nil(t, err)
			assert.Equal(t, rpm, tmp.Name)
		})

		t.Run("test remove rpm", func(t *testing.T) {
			assert.Nil(t, osobj.RemoveRpm(rpm))
			_, err := osobj.GetRpmInfo(rpm)
			assert.NotNil(t, err)
		})
	} else if exitc == 127 && stdo == "" && strings.Contains(stde, "command not found") == true && err == nil {
		t.Errorf("[TestInstallAndRemoveRpm]rpm command not found: %d, %s, %s, %v\n", exitc, stdo, stde, err)
	} else {
		t.Errorf("[TestInstallAndRemoveRpm]other error: %d, %s, %s, %v\n", exitc, stdo, stde, err)
	}
}
