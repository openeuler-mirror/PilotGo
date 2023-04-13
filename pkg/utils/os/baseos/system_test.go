package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostInfo(t *testing.T) {
	var osobj BaseOS
	tmp := osobj.GetHostInfo()
	assert.NotNil(t, tmp.IP)
	assert.NotNil(t, tmp.HostId)
	assert.NotNil(t, tmp.KernelArch)
	assert.NotNil(t, tmp.KernelVersion)
	assert.NotNil(t, tmp.Platform)
	assert.NotNil(t, tmp.PlatformVersion)
	assert.NotNil(t, tmp.Uptime)
}
