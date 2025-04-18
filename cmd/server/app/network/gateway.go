/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Fri Mar 21 16:18:53 2025 +0800
 */
package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/cmd/server/app/network/websocket"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auth"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/gateway"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"k8s.io/klog/v2"
)

func HttpGatewayServerInit(conf *options.HttpServer, stopCh <-chan struct{}) error {
	if err := SessionManagerInit(conf); err != nil {
		return err
	}

	go func() {
		r := SetupRouter()
		// start websocket server
		go websocket.CliManager.Start(stopCh)

		shutdownCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		srv := &http.Server{
			Handler: r,
		}
		ln, err := net.Listen("tcp", ":0") // 随机端口
		if err != nil {
			logger.Error("Failed to create listener: %v", err)
			return
		}
		addr := ln.Addr().(*net.TCPAddr)
		addr.IP = net.ParseIP(strings.Split(conf.Addr, ":")[0])

		if err := startGateway(shutdownCtx, conf, addr); err != nil {
			logger.Error("failed to start gateway, error:%v", err)
		}

		go func() {
			<-stopCh
			klog.Warningln("httpserver prepare stop")
			_ = srv.Shutdown(shutdownCtx)
		}()

		if conf.UseHttps {
			if conf.CertFile == "" || conf.KeyFile == "" {
				logger.Error("https cert or key not configd")
				return
			}

			if err := srv.ServeTLS(ln, conf.CertFile, conf.KeyFile); err != nil {
				if err != http.ErrServerClosed {
					logger.Error("ListenAndServeTLS start http server failed:%v", err)
					return
				}
			}
		} else {
			if err := srv.Serve(ln); err != nil {
				if err != http.ErrServerClosed {
					logger.Error("ListenAndServe start http server failed:%v", err)

				}

			}
		}

	}()
	if conf.Debug {
		go func() {
			// pprof
			portIndex := strings.Index(conf.Addr, ":")
			addr := conf.Addr[:portIndex] + ":6060"
			logger.Debug("start pprof service on: %s", addr)
			if conf.UseHttps {
				if conf.CertFile == "" || conf.KeyFile == "" {
					logger.Error("https cert or key not configd")
					return
				}

				err := http.ListenAndServeTLS(addr, conf.CertFile, conf.KeyFile, nil)
				if err != nil {
					logger.Error("failed to start pprof, error:%v", err)
				}
			} else {
				err := http.ListenAndServe(addr, nil)
				if err != nil {
					logger.Error("failed to start pprof, error:%v", err)
				}
			}
		}()
	}

	return nil
}

func startGateway(ctx context.Context, conf *options.HttpServer, addr *net.TCPAddr) error {
	reg, err := registry.NewServiceRegistrar(&registry.Options{
		Endpoints:   []string{"localhost:2379"},
		ServiceAddr: addr.String(),
		ServiceName: "pilotgo-server",
		Version:     "v3.0",
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %v", err)
	}

	watchCallback := func(eventType registry.EventType, service *registry.ServiceInfo) {
		if perms, ok := service.Metadata["permissions"]; ok {

			jsonData, err := json.Marshal(perms)
			if err != nil {
				return
			}
			var permissions []common.Permission
			if err := json.Unmarshal(jsonData, &permissions); err != nil {
				return
			}

			switch eventType {
			case registry.EventTypePut:
				if err := auth.AddPluginServicePermission("admin", permissions, service.ServiceName); err != nil {
					logger.Error("Failed to add permissions for service %s: %v", service.ServiceName, err)
				}
			case registry.EventTypeDelete:
				if err := auth.DeletePluginServicePermission(permissions, service.ServiceName); err != nil {
					logger.Error("Failed to remove permissions for service %s: %v", service.ServiceName, err)
				}
			}
		}
	}

	global.GW = gateway.NewCaddyGateway(reg, conf.Addr, watchCallback)

	go func() {
		if err := global.GW.Run(); err != nil {
			logger.Error("Gateway encountered an error: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Info("Gateway stopped successfully")

		if err := global.GW.Stop(); err != nil {
			logger.Error("Failed to stop Gateway: %v", err)
		} else {
			logger.Info("Gateway stopped successfully")
		}
	}()
	return nil
}
