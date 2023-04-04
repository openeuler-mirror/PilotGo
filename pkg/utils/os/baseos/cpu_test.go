package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var osobj BaseOS

func TestGetcpuinfo(t *testing.T) {
	assert.NotNil(t, osobj.GetCPUInfo())
}
