/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: yangzhao1
 * Date: 2022-04-05 10:28:34
 * LastEditTime: 2022-04-06 15:08:42
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/

package config

import (
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type Server struct {
	Addr string `yaml:"addr"`
}

type AgentConfig struct {
	Server  Server         `yaml:"server"`
	Logopts logger.LogOpts `yaml:"log"`
}

const config_file = "./config_agent.yaml"

var global_config AgentConfig

func Init() error {
	return utils.Load(config_file, &global_config)
}

func Config() *AgentConfig {
	return &global_config
}
