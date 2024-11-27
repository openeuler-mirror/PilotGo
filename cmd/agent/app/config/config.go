/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package config

import (
	"flag"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/config"
)

type Server struct {
	Addr string `yaml:"addr"`
}

type AgentConfig struct {
	Server  Server         `yaml:"server"`
	Logopts logger.LogOpts `yaml:"log"`
}

var Config_file string

var global_config AgentConfig

func Init() error {
	flag.StringVar(&Config_file, "conf", "./config_agent.yaml", "pilotgo-agent configuration file")
	flag.Parse()
	return config.Load(Config_file, &global_config)
}

func Config() *AgentConfig {
	return &global_config
}
