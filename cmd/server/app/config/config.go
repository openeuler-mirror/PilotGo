/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package config

import "gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"

// type HttpServer struct {
// 	Addr          string `yaml:"addr"`
// 	SessionCount  int    `yaml:"session_count"`
// 	SessionMaxAge int    `yaml:"session_max_age"`
// 	Debug         bool   `yaml:"debug"`
// 	UseHttps      bool   `yaml:"use_https"`
// 	CertFile      string `yaml:"cert_file"`
// 	KeyFile       string `yaml:"key_file"`
// }

// type SocketServer struct {
// 	Addr string `yaml:"addr"`
// }

// type JWTConfig struct {
// 	SecretKey string `yaml:"secret_key"`
// }

// type MysqlDBInfo struct {
// 	HostName string `yaml:"host_name"`
// 	UserName string `yaml:"user_name"`
// 	Password string `yaml:"password"`
// 	DataBase string `yaml:"data_base"`
// 	Port     int    `yaml:"port"`
// }

// type RedisDBInfo struct {
// 	RedisConn   string        `yaml:"redis_conn"`
// 	UseTLS      bool          `yaml:"use_tls"`
// 	RedisPwd    string        `yaml:"redis_pwd"`
// 	DefaultDB   int           `yaml:"defaultDB"`
// 	DialTimeout time.Duration `yaml:"dialTimeout"`
// 	EnableRedis bool          `yaml:"enableRedis"`
// }

// type Storage struct {
// 	Path string `yaml:"path"`
// }

// type ServerConfig struct {
// 	HttpServer   HttpServer     `yaml:"http_server"`
// 	SocketServer SocketServer   `yaml:"socket_server"`
// 	JWT          JWTConfig      `api:"jwt"`
// 	Logopts      logger.LogOpts `yaml:"log"`
// 	MysqlDBinfo  MysqlDBInfo    `yaml:"mysql"`
// 	RedisDBinfo  RedisDBInfo    `yaml:"redis"`
// 	Storage      Storage        `yaml:"storage"`
// }

// var global_config ServerConfig

// func Init(path string) error {
// 	return config.Load(path, &global_config)
// }

//	func Config() *ServerConfig {
//		return &global_config
//	}
var OptionsConfig *options.ServerConfig
