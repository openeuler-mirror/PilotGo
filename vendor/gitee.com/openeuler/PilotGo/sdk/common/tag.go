/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Thu Nov 2 16:38:53 2023 +0800
 */
package common

const (
	TypeOk    = "ok"
	TypeWarn  = "warn"
	TypeError = "error"
)

type Tag struct {
	UUID       string `json:"machineuuid"`
	PluginName string `json:"plugin_name"`
	Type       string `json:"type"`
	Data       string `json:"data"`
}
