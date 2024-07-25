/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-11-18 10:25:52
 * LastEditTime: 2023-06-12 15:18:57
 * Description: server main
 ******************************************************************************/
package main

import (
	"gitee.com/openeuler/PilotGo/app/server/cmd/commands"
	_ "gitee.com/openeuler/PilotGo/app/server/docs"
)

// @title          PilotGo Swagger  API
// @version         1.0
// @description     This is a pilotgo server API docs.
// @license.name  MulanPSL2
// @license.url   http://license.coscl.org.cn/MulanPSL2
// @host      localhost:8888
// @BasePath  /api/v1
// SwaggerUI: http://localhost:8888/swagger/index.html
func main() {
	commands.Execute()
}
