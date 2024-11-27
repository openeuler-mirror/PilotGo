/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Wed Dec 27 19:38:17 2023 +0800
 */
package client

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func (c *Client) Wait4Bind() {
	c.cond.L.Lock()
	logger.Debug("等待bind中...")
	for c.server == "" {
		c.cond.Wait() // 等待条件变量被通知
	}
	c.cond.L.Unlock()
	logger.Debug("bind 成功！")
}
