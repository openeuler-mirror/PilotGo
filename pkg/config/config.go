package config

/**
 * @Author: zhang han
 * @Date: 2021/11/1 16:08
 * @Description:
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	configType = "yaml"
)

var pilogo_config_file_name = "pkg/config/config.json"

type LogOpts struct {
	LogLevel  string `json:"log_level"`
	LogDriver string `json:"log_driver"`
	LogPath   string `json:"log_path"`
	MaxFile   int    `json:"max_file"`
	MaxSize   int    `json:"max_size"`
}
type Server struct {
	ServerPort int `json:"server_port"`
}
type DbInfo struct {
	HostName string `json:"host_name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	DataBase string `json:"data_base"`
	Port     int    `json:"port"`
}

type Configure struct {
	Logopts      LogOpts `json:"log_opts"`
	S            Server  `json:"server"`
	MaxAge       int     `json:"max_age"`
	SessionCount int     `json:"session_count"`
	Dbinfo       DbInfo  `json:"db_info"`
}

func Load() (*Configure, error) {
	var config Configure
	bytes, err := ioutil.ReadFile(pilogo_config_file_name)
	if err != nil {
		fmt.Printf("open %s failed! err = %s\n", pilogo_config_file_name, err.Error())
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("json.Unmarshal %s failed!\n", string(bytes))
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
		_, _ = fmt.Fprintf(output, "Config file changed %s \n", e.Name)
	})
	return nil
}

func MustInit(output io.Writer, conf string) { // MustInit if fail panic
	if err := Init(output, conf); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
