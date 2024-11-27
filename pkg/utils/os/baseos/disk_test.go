/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Tue Apr 4 11:36:42 2023 +0800
 */
package baseos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test disk
func TestGetDiskInfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetDiskInfo()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}

func TestGetDiskUsageInfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetDiskUsageInfo()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
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
