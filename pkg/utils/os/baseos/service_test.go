package baseos

import (
	"fmt"
	"testing"

	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceList(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetServiceList()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestGetService(t *testing.T) {
	var osobj BaseOS
	service := "mysqld"
	tmp, err := osobj.GetService(service)
	fmt.Println(tmp)
	if tmp == nil {
		fmt.Println("service is nil")
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, "mysqld", tmp.ServiceName)
	assert.Equal(t, "mysqld.service", tmp.UnitName)
	assert.Equal(t, common.ServiceUnit, tmp.UnitType)
	assert.Equal(t, "/usr/lib/systemd/system/mysqld.service", tmp.ServicePath)
	assert.NotEqual(t, "", tmp.ServiceActiveStatus)
	assert.NotEqual(t, "", tmp.ServiceLoadedStatus)
}

func TestConfigService(t *testing.T) {
	var osobj BaseOS
	service := "named.service"
	tmp, err := osobj.GetService(service)
	switch tmp.ServiceActiveStatus {
	case "active", "inactive":
		break
	default:
		assert.Nil(t, err)
	}
	if tmp.ServiceActiveStatus == "inactive" {
		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetService(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status.ServiceActiveStatus)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetService(service)
			switch status.ServiceActiveStatus {
			case "active", "inactive":
				assert.Equal(t, "inactive", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive", status)
			}
		})

		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetService(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status.ServiceActiveStatus)
		})
	} else {
		t.Run("test RestartService", func(t *testing.T) {
			assert.Nil(t, osobj.RestartService(service))
			status, err := osobj.GetService(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status.ServiceActiveStatus)
		})

		t.Run("test StopService", func(t *testing.T) {
			assert.Nil(t, osobj.StopService(service))
			status, err := osobj.GetService(service)
			switch status.ServiceActiveStatus {
			case "active", "inactive":
				assert.Equal(t, "inactive", status)
			default:
				assert.Nil(t, err)
				assert.Equal(t, "inactive", status)
			}
		})

		t.Run("test StartService", func(t *testing.T) {
			assert.Nil(t, osobj.StartService(service))
			status, err := osobj.GetService(service)
			assert.Nil(t, err)
			assert.Equal(t, "active", status.ServiceActiveStatus)
		})
	}

}
