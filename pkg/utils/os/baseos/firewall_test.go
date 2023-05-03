package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test firewall
func TestFirewall(t *testing.T) {
	var osobj BaseOS
	defaultzone := "work"
	dropzone := "drop"
	t.Run("test FirewallConfig", func(t *testing.T) {
		_, err := osobj.Config()
		assert.Nil(t, err)
	})

	t.Run("test FirewalldZoneConfig", func(t *testing.T) {
		_, err := osobj.FirewalldZoneConfig(defaultzone)
		assert.Nil(t, err)
	})

	t.Run("test FirewalldSetDefaultZone", func(t *testing.T) {
		osobj.FirewalldSetDefaultZone(defaultzone)
		tmp, err := osobj.Config()
		assert.Nil(t, err)
		new_zone := tmp.DefaultZone
		assert.Equal(t, defaultzone, new_zone)
		osobj.FirewalldSetDefaultZone("public")
	})

	t.Run("test Stop", func(t *testing.T) {
		osobj.Stop()
		tmp, err := osobj.Config()
		assert.Nil(t, err)
		new_status := tmp.Status
		assert.Equal(t, "not running", new_status)
	})

	t.Run("test Restart", func(t *testing.T) {
		osobj.Restart()
		tmp, err := osobj.Config()
		assert.Nil(t, err)
		new_status := tmp.Status
		assert.Equal(t, "running", new_status)
	})

	t.Run("test FirewalldSourceAdd", func(t *testing.T) {
		osobj.FirewalldSourceAdd(dropzone, "192.168.75.200")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_sources := new_zoneconfig.Sources
		assert.Equal(t, []string{"192.168.75.200"}, new_sources)
	})

	t.Run("test FirewalldSourceRemove", func(t *testing.T) {
		osobj.FirewalldSourceRemove(dropzone, "192.168.75.200")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_sources := new_zoneconfig.Sources
		assert.Equal(t, []string{""}, new_sources)
	})

	t.Run("test FirewalldServiceAdd", func(t *testing.T) {
		osobj.FirewalldServiceAdd(dropzone, "dns")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_service := new_zoneconfig.Service
		assert.Equal(t, []string{"dns"}, new_service)
	})

	t.Run("test FirewalldServiceRemove", func(t *testing.T) {
		osobj.FirewalldServiceRemove(dropzone, "dns")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_service := new_zoneconfig.Service
		assert.Equal(t, []string{""}, new_service)
	})

	t.Run("test AddZonePort", func(t *testing.T) {
		osobj.AddZonePort(dropzone, "53", "tcp")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_ports := new_zoneconfig.Ports.([]map[string]string)
		assert.Equal(t, []map[string]string{{"port": "53", "protocol": "tcp"}}, new_ports)
	})

	t.Run("test DelZonePort", func(t *testing.T) {
		osobj.DelZonePort(dropzone, "53", "tcp")
		new_zoneconfig, err := osobj.FirewalldZoneConfig(dropzone)
		assert.Nil(t, err)
		new_ports := new_zoneconfig.Ports.([]map[string]string)
		assert.Equal(t, []map[string]string{}, new_ports)
	})
}
