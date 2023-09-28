package baseos

import (
	"testing"

	"gitee.com/PilotGo/PilotGo/utils/os/common"
	"github.com/stretchr/testify/assert"
)

func TestYumsource(t *testing.T) {
	tmp, err := common.GetRepoSource()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
