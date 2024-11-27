/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Tue Apr 4 10:12:19 2023 +0800
 */
package baseos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetcpuinfo(t *testing.T) {
	var osobj BaseOS
	tmp, err := osobj.GetCPUInfo()
	assert.Nil(t, err)
	assert.NotNil(t, tmp)
}
