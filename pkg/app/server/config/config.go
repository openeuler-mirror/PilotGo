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
 * LastEditTime: 2023-09-04 16:16:36
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/
package config

import (
	"time"

	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/sdk/logger"
)

type HttpServer struct {
	Addr          string `yaml:"addr"`
	SessionCount  int    `yaml:"session_count"`
	SessionMaxAge int    `yaml:"session_max_age"`
	Debug         bool   `yaml:"debug"`
	UseHttps      bool   `yaml:"use_https"`
	CertFile      string `yaml:"cert_file"`
	KeyFile       string `yaml:"key_file"`
}

type SocketServer struct {
	Addr string `yaml:"addr"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
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
	JWT          JWTConfig      `api:"jwt"`
	Logopts      logger.LogOpts `yaml:"log"`
	MysqlDBinfo  MysqlDBInfo    `yaml:"mysql"`
	RedisDBinfo  RedisDBInfo    `yaml:"redis"`
}

var global_config ServerConfig

func Init(path string) error {
	return utils.Load(path, &global_config)
}

func Config() *ServerConfig {
	return &global_config
}
