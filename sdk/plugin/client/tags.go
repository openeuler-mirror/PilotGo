/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Aug 06 17:35:12 2025 +0800
 */
package client

func (c *Client) OnGetTags(callback GetTagsCallback) {
	c.getTagsCallback = callback
}
