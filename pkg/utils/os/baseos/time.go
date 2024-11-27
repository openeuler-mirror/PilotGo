/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package baseos

import (
	"fmt"

	"gitee.com/openeuler/PilotGo/pkg/utils"
)

// 获取机器时间
func (b *BaseOS) GetTime() (string, error) {
	exitc, nowtime, stde, err := utils.RunCommand("date +%s")
	if exitc == 0 && nowtime != "" && stde == "" && err == nil {
		return nowtime, nil
	}
	return "", fmt.Errorf("failed to get unix time: %d, %s, %s, %v", exitc, nowtime, stde, err)
}
