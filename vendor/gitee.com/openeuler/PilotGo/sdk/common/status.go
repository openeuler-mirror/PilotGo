/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
// 提供公共数据结构定义
package common

const (
	// 插件正在运行
	StatusRunning = "running"
	// 插件已加载，但未运行
	StatusLoaded = "loaded"
	// 插件离线，无法访问
	StatusOffline = "offline"
)
