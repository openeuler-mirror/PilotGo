package config

/**
 * @Author: zhang han
 * @Date: 2021/11/1 16:08
 * @Description:
 */

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configType    = "yaml"
	Log_FILE_PATH = "/var/log/pilotgo"
	LOG_FILE_NAME = "pilotgo.log"
)

var pilogo_config_file_name = "./config.yaml"

type LogOpts struct {
	LogLevel  string `yaml:"log_level"`
	LogDriver string `yaml:"log_driver"`
	LogPath   string `yaml:"log_path"`
	MaxFile   int    `yaml:"max_file"`
	MaxSize   int    `yaml:"max_size"`
}
type Server struct {
	ServerPort int `yaml:"server_port"`
}
type DbInfo struct {
	HostName string `yaml:"host_name"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	DataBase string `yaml:"data_base"`
	Port     int    `yaml:"port"`
}

type Configure struct {
	Logopts      LogOpts `yaml:"log_opts"`
	S            Server  `yaml:"server"`
	MaxAge       int     `yaml:"max_age"`
	SessionCount int     `yaml:"session_count"`
	Dbinfo       DbInfo  `yaml:"db_info"`
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
