/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package os

import (
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils/os/centos"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/kylin"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/nestos"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/openeuler"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

const (
	OpenEuler = "openeuler"
	Kylin     = "kylin"
	NestOS    = "nestos"
	CentOS    = "centos"
)

func OS() common.OSOperator {
	osinfo, err := common.InitOSName()
	if err != nil {
		logger.Error("osname init failed: %s", err)
		return nil
	}
	switch strings.ToLower(osinfo.OSName) {
	case OpenEuler:
		return &openeuler.OpenEuler{}
	case Kylin:
		return &kylin.Kylin{}
	case NestOS:
		switch strings.ToLower(osinfo.ID) {
		case "nestos for container":
			return &nestos.NestOS4Container{}
		case "nestos for virt":
			return &nestos.NestOS4Virt{}
		}
	case CentOS:
		return &centos.CentOS{}
	}
	return nil
}
