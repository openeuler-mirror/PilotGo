package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTime(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetTime()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
