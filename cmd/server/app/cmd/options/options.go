/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
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

var (
	defaultConfigurationName = "config_server"
	DefaultConfigFilePath    = fmt.Sprintf("./%s.yaml", defaultConfigurationName)
)

func NewHttpServerOptions() *HttpServer {
	return &HttpServer{
		Addr:          "0.0.0.0:8888",
		SessionCount:  100,
		SessionMaxAge: 1800,
		Debug:         false,
		UseHttps:      false,
		CertFile:      "",
		KeyFile:       "",
	}
}

type HttpServer struct {
	Addr          string `yaml:"addr" mapstructure:"addr"`
	SessionCount  int    `yaml:"session_count" mapstructure:"session_count"`
	SessionMaxAge int    `yaml:"session_max_age" mapstructure:"session_max_age"`
	Debug         bool   `yaml:"debug" mapstructure:"debug" comment:"if true, will start pprof on :6060"`
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

func NewSocketServerOptions() *SocketServer {
	return &SocketServer{
		Addr: "0.0.0.0:8879",
	}
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

func NewJWTConfigOptions() *JWTConfig {
	return &JWTConfig{
		SecretKey: "",
	}
}
func NewLogOptsOptions() *logger.LogOpts {
	return &logger.LogOpts{
		Level:   "debug",
		Driver:  "file",
		Path:    "./log/pilotgo_server.log",
		MaxFile: 1,
		MaxSize: 10485760,
	}
}

type MysqlDBInfo struct {
	HostName string `yaml:"host_name" mapstructure:"host_name"`
	UserName string `yaml:"user_name" mapstructure:"user_name" comment:"this is the username of database"`
	Password string `yaml:"password" mapstructure:"password"`
	DataBase string `yaml:"data_base" mapstructure:"data_base"`
	Port     int    `yaml:"port" mapstructure:"port"`
}

func NewMysqlDBInfoOpts() *MysqlDBInfo {
	return &MysqlDBInfo{
		HostName: "localhost",
		UserName: "root",
		Password: "",
		DataBase: "PilotGo",
		Port:     3306,
	}
}

type RedisDBInfo struct {
	RedisConn   string        `yaml:"redis_conn" mapstructure:"redis_conn"`
	UseTLS      bool          `yaml:"use_tls" mapstructure:"use_tls"`
	RedisPwd    string        `yaml:"redis_pwd" mapstructure:"redis_pwd"`
	DefaultDB   int           `yaml:"defaultDB" mapstructure:"defaultDB"`
	DialTimeout time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout" comment:"redis连接超时时间.默认5s"`
	EnableRedis bool          `yaml:"enableRedis" mapstructure:"enableRedis" comment:"是否启用redis"`
}

func NewRedisDBInfoOpts() *RedisDBInfo {
	return &RedisDBInfo{
		RedisConn:   "localhost:6379",
		UseTLS:      false,
		RedisPwd:    "",
		DefaultDB:   0,
		DialTimeout: time.Duration(5 * time.Second),
		EnableRedis: true,
	}
}

type Storage struct {
	Path string `yaml:"path" mapstructure:"path" comment:"文件服务存储路径"`
}

func NewStorageOptions() *Storage {
	return &Storage{
		Path: "",
	}
}

type Etcd struct {
	Endpoints   []string      `yaml:"endpoints" mapstructure:"endpoints"`
	ServiveName string        `yaml:"service_name" mapstructure:"service_name"`
	Version     string        `yaml:"version" mapstructure:"version"`
	DialTimeout time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
}

func NewEtcdOptions() *Etcd {
	return &Etcd{
		Endpoints:   []string{"localhost:2379"},
		ServiveName: "pilotgo-server",
		Version:     "3.0",
		DialTimeout: 5 * time.Second,
	}
}

type ServerConfig struct {
	HttpServer   *HttpServer     `yaml:"http_server" mapstructure:"http_server"`
	SocketServer *SocketServer   `yaml:"socket_server" mapstructure:"socket_server"`
	JWT          *JWTConfig      `api:"jwt" yaml:"jwt" mapstructure:"jwt"`
	Logopts      *logger.LogOpts `yaml:"log" mapstructure:"log"`
	MysqlDBinfo  *MysqlDBInfo    `yaml:"mysql" mapstructure:"mysql"`
	RedisDBinfo  *RedisDBInfo    `yaml:"redis" mapstructure:"redis"`
	Storage      *Storage        `yaml:"storage" mapstructure:"storage"`
	Etcd         *Etcd           `yaml:"etcd" mapstructure:"etcd"`
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
