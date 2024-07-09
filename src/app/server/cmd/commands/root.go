package commands

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"gitee.com/openeuler/PilotGo/app/server/cmd/options"
	"gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/network"
	"gitee.com/openeuler/PilotGo/app/server/network/websocket"
	"gitee.com/openeuler/PilotGo/app/server/service/auth"
	"gitee.com/openeuler/PilotGo/app/server/service/eventbus"
	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/dbmanager"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/spf13/cobra"
)

const flagconfig = "conf"

var conut int64

func NewServerCommand() *cobra.Command {

	s := options.NewServerOptions()
	conf, err := options.TryLoadFromDisk()
	if err == nil {
		s.ServerConfig = conf
		klog.InfoS("TryLoadFromDisk pilotgo config !", "HttpServer", *s.ServerConfig.HttpServer)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "SocketServer", *s.ServerConfig.SocketServer)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "JWT", *s.ServerConfig.JWT)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "Logopts", *s.ServerConfig.Logopts)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "RedisDBinfo", *s.ServerConfig.RedisDBinfo)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "MysqlDBinfo", *s.ServerConfig.MysqlDBinfo)
		klog.InfoS("TryLoadFromDisk pilotgo config !", "Storage", *s.ServerConfig.Storage)
	} else {
		klog.Fatal("Failed to load configuration from disk", err)
	}
	cmd := &cobra.Command{
		Use:  "pilotgo",
		Long: `Run the pilotgo API server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(s, signals.SetupSignalHandler(), cmd, options.WatchConfigChange())
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

func run(_ *options.ServerOptions, ctx context.Context, cmd *cobra.Command) error {
	if atomic.LoadInt64(&conut) > 0 {
		return nil
	}
	atomic.AddInt64(&conut, 1)

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
	if err := dbmanager.RedisdbInit(&config.Config().RedisDBinfo, ctx.Done()); err != nil {
		logger.Error("redis db init failed, please check again: %s", err)
		return err
	}

	// mysql db初始化
	if err := dbmanager.MysqldbInit(&config.Config().MysqlDBinfo); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		return err
	}

	// 启动agent socket server
	if err := network.SocketServerInit(&config.Config().SocketServer, ctx.Done()); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		return err
	}

	//此处启动前端及REST http server
	err = network.HttpServerInit(&config.Config().HttpServer, ctx.Done())
	if err != nil {
		logger.Error("HttpServerInit socket server init failed, error:%v", err)
		return err
	}

	// 启动所有功能模块服务
	if err := startServices(ctx.Done()); err != nil {
		logger.Error("start services error: %s", err)
		return err
	}

	// 前端推送告警
	go websocket.SendWarnMsgToWeb(ctx.Done())

	logger.Info("start to serve")
	atomic.AddInt64(&conut, -1)
	// 信号监听 redis
	return nil

}
func startServices(stopCh <-chan struct{}) error {
	// 鉴权模块初始化
	auth.Casbin(&config.Config().MysqlDBinfo)

	// 初始化plugin服务
	plugin.Init(stopCh)

	//初始化eventbus
	eventbus.Init(stopCh)

	return nil
}

func Run(s *options.ServerOptions, ctx context.Context, cmd *cobra.Command, configCh <-chan options.ServerConfig) error {

	cctx, cancelFunc := context.WithCancel(context.TODO())
	errCh := make(chan error)
	defer close(errCh)
	go func() {
		if err := runer(s, cctx, cmd); err != nil {
			klog.Errorf("runner start error:%v", err)
			errCh <- err
		}
	}()
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			klog.Warningln("pilotgo exit bye")
			return nil
		case cfg := <-configCh:
			klog.Warningln("config is change")
			cancelFunc()
			s.ServerConfig = &cfg
			cctx, cancelFunc = context.WithCancel(context.TODO())
			go func() {
				if err := runer(s, cctx, cmd); err != nil {
					klog.Errorf("config is change , reboot server error:%v", err)
					errCh <- err
				}
			}()
		case err := <-errCh:
			cancelFunc()
			return err
		}
	}

}

func runer(s *options.ServerOptions, ctx context.Context, cmd *cobra.Command) error {
	return run(s, ctx, cmd)
}
