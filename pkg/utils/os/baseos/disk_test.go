package baseos

import (
	"os"
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
	mountpath := "/root/mountdir"

	t.Run("test DiskMount", func(t *testing.T) {
		tmp, err := osobj.DiskMount(diskpath, mountpath)
		assert.Nil(t, err)
		assert.Equal(t, "", tmp)
	})

	t.Run("test DiskUMount", func(t *testing.T) {
		tmp, err := osobj.DiskUMount(diskpath)
		assert.Nil(t, err)
		assert.Equal(t, "", tmp)
		Err := os.RemoveAll(mountpath)
		assert.Nil(t, Err)
	})

	t.Run("test DiskFormat", func(t *testing.T) {
		tmp, err := osobj.DiskFormat("ext4", diskpath)
		assert.Nil(t, err)
		assert.NotNil(t, tmp)
	})
}
