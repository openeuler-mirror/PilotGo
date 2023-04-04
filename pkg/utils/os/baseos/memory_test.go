package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemoryConfig(t *testing.T) {
	var osobj BaseOS
	assert.NotNil(t, osobj.GetMemoryConfig())
}
