package baseos

import (
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
	case "active", "inactive":
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
	case "active", "inactive":
		break
	default:
		assert.Nil(t, err)
	}
	if tmp == "inactive" {
		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetServiceStatus(service)
			switch status {
			case "active", "inactive":
				assert.Equal(t, "inactive", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive", status)
			}
		})

		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status)
		})
	} else {
		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetServiceStatus(service)
			switch status {
			case "active", "inactive":
				assert.Equal(t, "inactive", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive", status)
			}
		})

		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetServiceStatus(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status)
		})
	}

}
