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
 * Description: monitor init
 ******************************************************************************/
package initialization

import (
	"time"

	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/controller"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/global"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func MonitorInit(conf *sconfig.Monitor) error {
	go func() {
		logger.Info("start monitor")
		err := controller.InitPromeYml()
		if err != nil {
			logger.Error("初始化promethues配置文件失败")
		}
		for {
			// TODO: 重构为事件触发机制
			a := make([]map[string]string, 0)
			var m []model.MachineNode
			global.PILOTGO_DB.Find(&m)
			for _, value := range m {
				r := map[string]string{}
				r[value.MachineUUID] = value.IP
				a = append(a, r)
			}
			err := controller.WriteYml(a)
			if err != nil {
				logger.Error("%s", err.Error())
			}
			time.Sleep(100 * time.Second)
		}

	}()

	return nil
}
