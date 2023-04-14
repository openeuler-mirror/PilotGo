package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetcpuinfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetCPUInfo()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
