package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSysctlConfig(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetSysctlConfig()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestConfigSysctl(t *testing.T) {
	var osobj BaseOS
	param := "net.ipv4.ip_forward"

	t.Run("test GetVarNameValue", func(t *testing.T) {
		tmp, err := osobj.GetVarNameValue(param)
		switch tmp {
		case "0", "1":
			break
		default:
			assert.Nil(t, err)
			assert.NotNil(t, tmp)
		}
	})

	t.Run("test TempModifyPar", func(t *testing.T) {
		tmp, err := osobj.GetVarNameValue(param)
		assert.Nil(t, err)
		switch tmp {
		case "0":
			newparam, err := osobj.TempModifyPar(param + "=1")
			assert.Nil(t, err)
			assert.Equal(t, "net.ipv4.ip_forward = 1", newparam)
			osobj.TempModifyPar(param + "=0")
		case "1":
			newparam, err := osobj.TempModifyPar(param + "=0")
			assert.Nil(t, err)
			assert.Equal(t, "net.ipv4.ip_forward = 0", newparam)
			osobj.TempModifyPar(param + "=1")
		default:
			t.Error("failed to test TempModifyPar")
		}
	})
}
