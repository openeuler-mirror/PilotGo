package options

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

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
var (
	defaultConfigurationName = "config_server"
	DefaultConfigFilePath    = fmt.Sprintf("./%s.yaml", defaultConfigurationName)
)

type HttpServer struct {
	Addr          string `yaml:"addr" mapstructure:"addr"`
	SessionCount  int    `yaml:"session_count" mapstructure:"session_count"`
	SessionMaxAge int    `yaml:"session_max_age" mapstructure:"session_max_age"`
	Debug         bool   `yaml:"debug" mapstructure:"debug"`
	UseHttps      bool   `yaml:"use_https" mapstructure:"use_https"`
	CertFile      string `yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile       string `yaml:"key_file" mapstructure:"key_file"`
}

func (hs *HttpServer) Validate() []error {
	var errors []error
	if len(hs.Addr) == 0 {
		errors = append(errors, fmt.Errorf("HttpServer config addr is nil"))
	}
	return errors
}

type SocketServer struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
}

func (hs *SocketServer) Validate() []error {
	var errors []error
	if len(hs.Addr) == 0 {
		errors = append(errors, fmt.Errorf("SocketServer config addr is nil"))
	}
	return errors
}

type JWTConfig struct {
	SecretKey string `yaml:"secret_key" mapstructure:"secret_key"`
}

type MysqlDBInfo struct {
	HostName string `yaml:"host_name" mapstructure:"host_name"`
	UserName string `yaml:"user_name" mapstructure:"user_name"`
	Password string `yaml:"password" mapstructure:"password"`
	DataBase string `yaml:"data_base" mapstructure:"data_base"`
	Port     int    `yaml:"port" mapstructure:"port"`
}

type RedisDBInfo struct {
	RedisConn   string        `yaml:"redis_conn" mapstructure:"redis_conn"`
	UseTLS      bool          `yaml:"use_tls" mapstructure:"use_tls"`
	RedisPwd    string        `yaml:"redis_pwd" mapstructure:"redis_pwd"`
	DefaultDB   int           `yaml:"defaultDB" mapstructure:"defaultDB"`
	DialTimeout time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
	EnableRedis bool          `yaml:"enableRedis" mapstructure:"enableRedis"`
}

type Storage struct {
	Path string `yaml:"path" mapstructure:"path"`
}

type ServerConfig struct {
	HttpServer   *HttpServer     `yaml:"http_server" mapstructure:"http_server"`
	SocketServer *SocketServer   `yaml:"socket_server" mapstructure:"socket_server"`
	JWT          *JWTConfig      `api:"jwt" yaml:"jwt" mapstructure:"jwt"`
	Logopts      *logger.LogOpts `yaml:"log" mapstructure:"log"`
	MysqlDBinfo  *MysqlDBInfo    `yaml:"mysql" mapstructure:"mysql"`
	RedisDBinfo  *RedisDBInfo    `yaml:"redis" mapstructure:"redis"`
	Storage      *Storage        `yaml:"storage" mapstructure:"storage"`
}

func (s *ServerConfig) Validate() []error {
	var errors []error
	errors = append(errors, s.HttpServer.Validate()...)
	errors = append(errors, s.SocketServer.Validate()...)
	return errors
}

type ServerOptions struct {
	Config       string
	ServerConfig *ServerConfig
}

func NewServerOptions() *ServerOptions {
	s := &ServerOptions{
		Config:       DefaultConfigFilePath,
		ServerConfig: New(),
	}
	return s
}
func (s *ServerOptions) NewAPIServer(stopCh <-chan struct{}) {
}

var (
	_config = defaultConfig()
)

type pilotgoConfig struct {
	cfg         *ServerConfig
	cfgChangeCh chan ServerConfig
	watchOnce   sync.Once
	loadOnce    sync.Once
}

func WatchConfigChange() <-chan ServerConfig {
	return _config.watchConfig()
}
func (c *pilotgoConfig) watchConfig() <-chan ServerConfig {
	c.watchOnce.Do(func() {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			cfg := New()
			if err := viper.Unmarshal(cfg); err != nil {
				klog.Errorf("config reload error", err)
			} else {
				if in.Op&fsnotify.Write != 0 && len(viper.AllKeys()) > 0 {
					c.cfgChangeCh <- *cfg
				}
			}
		})
	})
	return c.cfgChangeCh
}

func New() *ServerConfig {
	return &ServerConfig{}
}
func TryLoadFromDisk(configfile string) (*ServerConfig, error) {
	viper.SetConfigFile(configfile)
	viper.SetEnvPrefix("pilotgo")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return _config.loadFromDisk()
}

func (c *pilotgoConfig) loadFromDisk() (*ServerConfig, error) {
	var err error
	c.loadOnce.Do(func() {
		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				err = fmt.Errorf("error parsing configuration file %s", err)
			}
		}
		err = viper.Unmarshal(c.cfg)
		if err != nil {
			klog.Errorf("viper error: %v /n", err)
		}
	})
	return c.cfg, err
}
func defaultConfig() *pilotgoConfig {

	return &pilotgoConfig{
		cfg:         New(),
		cfgChangeCh: make(chan ServerConfig),
		watchOnce:   sync.Once{},
		loadOnce:    sync.Once{},
	}
}
