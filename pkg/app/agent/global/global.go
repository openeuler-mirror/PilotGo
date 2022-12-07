package global

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"openeuler.org/PilotGo/PilotGo/pkg/app/agent/network"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/message/protocol"
)

type ConfigMessage struct {
	ConfigName    string
	ConfigContent string
	ConfigChange  string
	Machine_uuid  string
}

//配置文件的监听器
func Configfsnotify(ConMess ConfigMessage, client *network.SocketClient) error {
	//创建一个监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("NewWatcher failed: ", err)
		return err
	}
	defer watcher.Close()
	done := make(chan bool)
	err = watcher.Add(ConMess.ConfigName)
	if err != nil {
		logger.Error("Add failed:", err)
	}
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Debug("在文件 %s 上进行 : %s", event.Name, event.Op)
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Write == fsnotify.Write {
					ConMess.ConfigChange = event.Op.String()
					ConMess.ConfigContent, err = utils.FileReadString(ConMess.ConfigName)
					if err != nil {
						logger.Debug("error:", err)
					}
					msg := &protocol.Message{
						UUID:   uuid.New().String(),
						Type:   protocol.ConfigFileMonitor,
						Status: 0,
						Data:   ConMess,
					}
					if err := client.Send(msg); err != nil {
						logger.Debug("send message failed, error:", err)
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					_, err := os.Stat(ConMess.ConfigName)
					if err == nil {
						err = watcher.Add(ConMess.ConfigName)
					}
				}
			case err, ok := <-watcher.Errors:
				logger.Debug("error:", err)
				if !ok {
					return
				}
			}
		}
	}()
	<-done
	return nil
}
