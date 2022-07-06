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
 * Date: 2022-07-05 13:03:16
 * LastEditTime: 2022-07-05 14:10:23
 * Description: http server init
 ******************************************************************************/
package initialization

import (
	"net/http"

	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func HttpServerInit(conf *sconfig.HttpServer) error {
	if err := SessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := SetupRouter()
		r.Run(conf.Addr)

		err := http.ListenAndServe(conf.Addr, nil) // listen and serve
		if err != nil {
			logger.Error("failed to start http server, error:%v", err)
		}
	}()

	return nil
}
func SessionManagerInit(conf *sconfig.HttpServer) error {
	var sessionManage service.SessionManage
	sessionManage.Init(conf.SessionMaxAge, conf.SessionCount)
	return nil
}
