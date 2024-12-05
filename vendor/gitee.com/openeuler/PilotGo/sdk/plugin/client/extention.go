/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Nov 7 15:01:17 2023 +0800
 */
package client

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
)

func (c *Client) RegisterExtention(exts []common.Extention) {
	c.extentions = append(c.extentions, exts...)
}
