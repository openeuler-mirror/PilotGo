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
 * Date: 2021-11-18 13:03:16
 * LastEditTime: 2022-04-02 18:03:06
 * Description: provide configure yaml functions.
 ******************************************************************************/
package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configType = "yaml"
)

var pilogo_config_file_name = "./config_server.yaml"

type LogOpts struct {
	Level   string `yaml:"level"`
	Driver  string `yaml:"driver"`
	Path    string `yaml:"path"`
	MaxFile int    `yaml:"max_file"`
	MaxSize int    `yaml:"max_size"`
}
type HttpServer struct {
	Addr          string `yaml:"addr"`
	SessionCount  int    `yaml:"session_count"`
	SessionMaxAge int    `yaml:"session_max_age"`
}
type SocketServer struct {
	Addr string `yaml:"addr"`
}
type Monitor struct {
	PrometheusAddr   string `yaml:"prometheus_addr"`
	AlertManagerAddr string `yaml:"alertmanager_addr"`
}
type DbInfo struct {
	HostName string `yaml:"host_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	DataBase string `yaml:"data_base"`
	Port     int    `yaml:"port"`
}

type Configure struct {
	HttpServer   HttpServer   `yaml:"http_server"`
	SocketServer SocketServer `yaml:"socket_server"`
	Monitor      Monitor      `yaml:"monitor"`
	Logopts      LogOpts      `yaml:"log"`
	Dbinfo       DbInfo       `yaml:"database"`
}

func Load() (*Configure, error) {
	config := Configure{}
	bytes, err := ioutil.ReadFile(pilogo_config_file_name)
	if err != nil {
		fmt.Printf("open %s failed! err = %s\n", pilogo_config_file_name, err.Error())
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("yaml Unmarshal %s failed!\n", string(bytes))
		return nil, err
	}
	return &config, nil
}

func Init(output io.Writer, configFile string) error {
	if output == nil {
		output = ioutil.Discard
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType(configType) // or viper.SetConfigType("YAML")
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Fprintf(output, "Config file changed %s \n", e.Name)
	})
	return nil
}

func MustInit(output io.Writer, conf string) { // MustInit if fail panic
	if err := Init(output, conf); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
