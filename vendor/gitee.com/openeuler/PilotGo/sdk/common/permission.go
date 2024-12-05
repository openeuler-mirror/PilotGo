/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Thu Jan 11 20:24:26 2024 +0800
 */
package common

type Permission struct {
	Resource string `json:"resource"` //字符串中不允许包含"/"
	Operate  string `json:"operate"`
}
