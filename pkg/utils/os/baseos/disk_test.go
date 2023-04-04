package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test disk
func TestGetDiskInfo(t *testing.T) {
	var osobj BaseOS
	assert.NotNil(t, osobj.GetDiskInfo())
}

func TestGetDiskUsageInfo(t *testing.T) {
	var osobj BaseOS
	assert.NotNil(t, osobj.GetDiskUsageInfo())
}

func TestDiskConfig(t *testing.T) {
	var osobj BaseOS
	diskpath := "/dev/nvme0n1"
	mountpath := "/root/PilotGo-plugins/PilotGo/pkg/utils/os/testapi/mountdir"

	t.Run("test DiskMount", func(t *testing.T) {
		assert.NotNil(t, osobj.DiskMount(diskpath, mountpath))
	})

	t.Run("test DiskUMount", func(t *testing.T) {
		assert.NotNil(t, osobj.DiskUMount(diskpath))
	})

	t.Run("test DiskFormat", func(t *testing.T) {
		assert.NotNil(t, osobj.DiskFormat("ext4", diskpath))
	})
}
