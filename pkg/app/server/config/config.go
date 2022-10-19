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
 * Date: 2022-04-06 13:27:45
 * LastEditTime: 2022-04-20 14:32:28
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/
package config

import (
	"time"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type HttpServer struct {
	Addr          string `yaml:"addr"`
	SessionCount  int    `yaml:"session_count"`
	SessionMaxAge int    `yaml:"session_max_age"`
}
type SocketServer struct {
	Addr string `yaml:"addr"`
}
type Monitor struct {
	PrometheusAddr    string `yaml:"prometheus_addr"`
	AlertManagerAddr  string `yaml:"alertmanager_addr"`
	AlertRulesPath    string `yaml:"alert_rules_path"`
	PrometheusYmlPath string `yaml:"prometheus_yml_path"`
}
type MysqlDBInfo struct {
	HostName string `yaml:"host_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	DataBase string `yaml:"data_base"`
	Port     int    `yaml:"port"`
}

type RedisDBInfo struct {
	RedisConn   string        `yaml:"redis_conn"`
	RedisPwd    string        `yaml:"redis_pwd"`
	DefaultDB   int           `yaml:"defaultDB"`
	DialTimeout time.Duration `yaml:"dialTimeout"`
	EnableRedis bool          `yaml:"enableRedis"`
}

type ServerConfig struct {
	HttpServer   HttpServer     `yaml:"http_server"`
	SocketServer SocketServer   `yaml:"socket_server"`
	Monitor      Monitor        `yaml:"monitor"`
	Logopts      logger.LogOpts `yaml:"log"`
	MysqlDBinfo  MysqlDBInfo    `yaml:"mysql"`
	RedisDBinfo  RedisDBInfo    `yaml:"redis"`
}

const config_file = "./config_server.yaml"

var global_config ServerConfig

func Init() error {
	return utils.Load(config_file, &global_config)
}

func Config() *ServerConfig {
	return &global_config
}
