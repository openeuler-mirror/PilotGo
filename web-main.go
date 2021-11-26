package main

import (
	"fmt"
	"openeluer.org/PilotGo/PilotGo/pkg/cmd"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"os"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}

	logger.Init(conf)
	logger.Info("Thanks to choose PilotGo!")

	err = cmd.Start(conf)
	if err != nil {
		logger.Info("server start failed:%s", err.Error())
		os.Exit(-1)
	}
}
