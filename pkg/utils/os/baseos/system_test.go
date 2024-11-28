/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Thu Apr 13 10:35:55 2023 +0800
 */
package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostInfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetHostInfo()
	assert.Nil(t, err)
	assert.NotNil(t, tmp.IP)
	assert.NotNil(t, tmp.HostId)
	assert.NotNil(t, tmp.KernelArch)
	assert.NotNil(t, tmp.KernelVersion)
	assert.NotNil(t, tmp.Platform)
	assert.NotNil(t, tmp.PlatformVersion)
	assert.NotNil(t, tmp.Uptime)
}
