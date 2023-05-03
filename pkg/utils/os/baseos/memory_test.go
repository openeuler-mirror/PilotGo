package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemoryConfig(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetMemoryConfig()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
