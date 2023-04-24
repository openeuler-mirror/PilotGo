package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYumsource(t *testing.T) {
	tmp, err := GetRepoSource()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
