/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package main

import (
	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/commands"
	_ "gitee.com/openeuler/PilotGo/cmd/server/app/docs"
)

// @title          PilotGo Swagger  API
// @version         2.0
// @description     This is a pilotgo server API docs.
// @license.name  MulanPSL2
// @license.url   http://license.coscl.org.cn/MulanPSL2
// @host      localhost:8888
// @BasePath  /api/v1
// SwaggerUI: http://localhost:8888/swagger/index.html
func main() {
	commands.Execute()
}
