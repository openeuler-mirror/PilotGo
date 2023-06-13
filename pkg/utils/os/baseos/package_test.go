package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func TestYumsource(t *testing.T) {
	tmp, err := common.GetRepoSource()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
