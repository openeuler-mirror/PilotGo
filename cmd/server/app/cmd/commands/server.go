/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package commands

import (
	"context"
	"fmt"
	"sync/atomic"

	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/cmd/server/app/config"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/websocket"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/pkg/dbmanager"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

const flagconfig = "conf"
const cliName = "start"

var conut int64

func NewServerCommand() *cobra.Command {
	s := options.NewServerOptions()
	cmd := &cobra.Command{
		Example: `
		# run  the pilotgo-server
		pilotgo-server 
		or
		pilotgo-server start
		`,
		Use:   cliName,
		Short: "Start the api-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config_file, err := cmd.Flags().GetString(flagconfig)
			if err != nil {
				return errors.Wrapf(err, "error accessing flag %s for command %s", flagconfig, cmd.Name())
			}
			klog.Infof("load configuration config_file is:%s", config_file)
			isExist, _ := utils.IsFileExist(config_file)
			if !isExist {
				return errors.Errorf("%s not found,please check the file. Flags: -h, --help   help for pilotgo", config_file)
			}
			conf, err := options.TryLoadFromDisk(config_file)
			if err == nil {
				s.ServerConfig = conf
				config.OptionsConfig = conf
			} else {
				klog.Fatal("Failed to load configuration from disk", err)
			}
			if errs := s.ServerConfig.Validate(); len(errs) != 0 {
				klog.Errorf("please check current config, errors is:%v", errs)
				return errors.New("please check config_server.yaml")
			}
			return Run(s, signals.SetupSignalHandler(), cmd, options.WatchConfigChange())
		},
	}
	cmd.ResetFlags()
	cmd.Flags().StringP(flagconfig, "c", options.DefaultConfigFilePath, "choose config yaml file")
	return cmd
}
func run(opts *options.ServerOptions, ctx context.Context, _ *cobra.Command) error {
	if atomic.LoadInt64(&conut) > 0 {
		return nil
	}
	atomic.AddInt64(&conut, 1)
	config := opts.ServerConfig
	if config.Storage.Path == "" {
		fmt.Println("Please set the path for file service storage in yaml")
		return errors.New("storage path is nil")
	}
	if err := logger.Init(config.Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		return err
	}
	logger.Info("Thanks to choose PilotGo!")

	// redis db初始化
	if err := dbmanager.RedisdbInit(config.RedisDBinfo, ctx.Done()); err != nil {
		if err == context.Canceled {
			return nil
		}
		logger.Error("redis db init failed, please check again: %s", err)
		return err
	}

	// mysql db初始化
	if err := dbmanager.MysqldbInit(config.MysqlDBinfo); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		return err
	}

	// 启动agent socket server
	if err := network.SocketServerInit(config.SocketServer, ctx.Done()); err != nil {
		logger.Error("socket server init failed, error:%v", err)
		return err
	}

	//此处启动前端及REST http server
	err := network.HttpServerInit(config.HttpServer, ctx.Done())
	if err != nil {
		logger.Error("HttpServerInit socket server init failed, error:%v", err)
		return err
	}

	// 启动所有功能模块服务
	if err := startServices(config.MysqlDBinfo, ctx.Done()); err != nil {
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
func startServices(mysqlInfo *options.MysqlDBInfo, stopCh <-chan struct{}) error {
	// 鉴权模块初始化
	auth.Casbin(mysqlInfo)

	// 初始化plugin服务
	plugin.Init(stopCh)

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
			config.OptionsConfig = &cfg
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
