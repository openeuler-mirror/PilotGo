package global

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

type ConfigMessage struct {
	ConfigName string
	ConfigType string
	ConfigPath string
}

func FileGetViper(ConMess ConfigMessage) error {
	viper.SetConfigName(ConMess.ConfigName)
	viper.SetConfigType(ConMess.ConfigType)
	viper.AddConfigPath(ConMess.ConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Info("read config failed: %v", err)
		return err
	}
	/*
		//输出配置文件的一些参数
		fmt.Println("mysql ip: ", viper.Get("mysql.host_name"))
		fmt.Println("mysql port: ", viper.Get("mysql.port"))
		fmt.Println("mysql user: ", viper.Get("mysql.user_name"))
		fmt.Println("redis ip: ", viper.Get("redis.redis_conn"))
	*/
	err = Configfsnotify(ConMess)
	if err != nil {
		logger.Info("listening config failed: %v", err)
		return err
	}
	return nil
}

//配置文件的监听器
func Configfsnotify(ConMess ConfigMessage) error {
	//创建一个监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("NewWatcher failed: ", err)
		return err
	}
	defer watcher.Close()
	done := make(chan bool)
	var event fsnotify.Event
	var ok bool
	go func() {
		defer close(done)
		for {
			select {
			case event, ok = <-watcher.Events:
				logger.Fatal("%s %s\n", event.Name, event.Op)
				if !ok {
					goto here
				}
			case err, ok = <-watcher.Errors:
				logger.Fatal("error:", err)
				if !ok {
					goto here
				}
			}
		}
	here:
		err = watcher.Add(ConMess.ConfigPath + ConMess.ConfigName + "." + ConMess.ConfigType)
		if err != nil {
			log.Fatal("Add failed:", err)
		}
	}()
	<-done
	return nil
}
