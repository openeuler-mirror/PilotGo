package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo/app/server/cmd/options"
	"gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/network"
	"gitee.com/openeuler/PilotGo/app/server/network/websocket"
	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"gitee.com/openeuler/PilotGo/app/server/service/eventbus"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/dbmanager"
	"gitee.com/openeuler/PilotGo/dbmanager/redismanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

const flagconfig = "conf"

func NewServerCommand() *cobra.Command {
	s := options.NewServerOptions()
	cmd := &cobra.Command{
		Use:  "pilotgo",
		Long: `Run the pilotgo API server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(s, signals.SetupSignalHandler(), cmd)
		},
		SilenceUsage: true,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}
	cmd.ResetFlags()
	cmd.Flags().StringP(flagconfig, "c", "./config_server.yaml", "yaml")
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of pilotgo server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("2.1.1")
		},
	}
	cmd.AddCommand(versionCmd)
	return cmd
}
func run(_ *options.ServerOptions, _ context.Context, cmd *cobra.Command) error {
	config_file, err := cmd.Flags().GetString(flagconfig)
	if err != nil {
		return errors.Wrapf(err, "error accessing flag %s for command %s", flagconfig, cmd.Name())
	}
	err = config.Init(config_file)
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		return err
	}
	if config.Config().Storage.Path == "" {
		fmt.Println("Please set the path for file service storage in yaml")
		return errors.New("storage path is nil")
	}
	if err := logger.Init(&config.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		return err
	}
	logger.Info("Thanks to choose PilotGo!")

	// redis db初始化
	if err := dbmanager.RedisdbInit(&config.Config().RedisDBinfo); err != nil {
		logger.Error("redis db init failed, please check again: %s", err)
		return err
	}

	// mysql db初始化
	if err := dbmanager.MysqldbInit(&config.Config().MysqlDBinfo); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		return err
	}

	// 启动agent socket server
	if err := network.SocketServerInit(&config.Config().SocketServer); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		return err
	}

	//此处启动前端及REST http server
	err = network.HttpServerInit(&config.Config().HttpServer)
	if err != nil {
		logger.Error("socket server init failed, error:%v", err)
		return err
	}

	// 启动所有功能模块服务
	if err := startServices(); err != nil {
		logger.Error("start services error: %s", err)
		return err
	}

	// 前端推送告警
	go websocket.SendWarnMsgToWeb()

	logger.Info("start to serve.")

	// 信号监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			// TODO: DO EXIT

			redismanager.Redis().Close()

			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
	return nil

}
func startServices() error {
	// 鉴权模块初始化
	auth.Casbin(&config.Config().MysqlDBinfo)

	// 初始化plugin服务
	plugin.Init()

	//初始化eventbus
	eventbus.Init()

	return nil
}
