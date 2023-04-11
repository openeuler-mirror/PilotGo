package baseos

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServiceList(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetServiceList()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestGetServiceStatus(t *testing.T) {
	var osobj BaseOS
	service := "named.service"
	tmp, err := osobj.GetServiceStatus(service)
	switch tmp {
	case "active\n", "inactive\n":
		break
	default:
		assert.Nil(t, err)
		assert.NotNil(t, tmp)
	}
}

func TestConfigService(t *testing.T) {
	var osobj BaseOS
	service := "named.service"
	tmp, err := osobj.GetServiceStatus(service)
	switch tmp {
	case "active\n", "inactive\n":
		break
	default:
		assert.Nil(t, err)
	}
	if strings.Replace(tmp, "\n", "", -1) == "inactive" {
		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active\n", status)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetServiceStatus(service)
			switch status {
			case "active\n", "inactive\n":
				assert.Equal(t, "inactive\n", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive\n", status)
			}
		})

		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active\n", status)
		})
	} else {
		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active\n", status)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetServiceStatus(service)
			switch status {
			case "active\n", "inactive\n":
				assert.Equal(t, "inactive\n", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive\n", status)
			}
		})

		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active\n", status)
		})
	}

}
