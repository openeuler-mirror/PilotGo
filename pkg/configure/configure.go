package configure

//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//)
//
//var pilogo_config_file_name = "./config.json"
//
//type LogOpts struct {
//	LogLevel  string `json:"log_level"`
//	LogDriver string `json:"log_driver"`
//	LogPath   string `json:"log_path"`
//	MaxFile   int    `json:"max_file"`
//	MaxSize   int    `json:"max_size"`
//}
//type Server struct {
//	ServerPort int `json:"server_port"`
//}
//type DbInfo struct {
//	HostName string `json:"host_name"`
//	UserName string `json:"user_name"`
//	Password string `json:"password"`
//	DataBase string `json:"data_base"`
//	Port     int    `json:"port"`
//}
//
//type Configure struct {
//	Logopts      LogOpts `json:"log_opts"`
//	S            Server  `json:"server"`
//	MaxAge       int     `json:"max_age"`
//	SessionCount int     `json:"session_count"`
//	Dbinfo       DbInfo  `json:"db_info"`
//}
//
//func Load() (*Configure, error) {
//	var config Configure
//	bytes, err := ioutil.ReadFile(pilogo_config_file_name)
//	if err != nil {
//		fmt.Printf("open %s failed! err = %s\n", pilogo_config_file_name, err.Error())
//		return nil, err
//	}
//
//	err = json.Unmarshal(bytes, &config)
//	if err != nil {
//		fmt.Printf("json.Unmarshal %s failed!\n", string(bytes))
//		return nil, err
//	}
//	return &config, nil
//}
