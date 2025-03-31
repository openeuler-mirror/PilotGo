/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Fei Mar 21 16:18:53 2025 +0800
 */
package network

import (
	"context"
	"fmt"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/cmd/options"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/gateway"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

// StartGateway initializes and starts the API Gateway
func StartGateway(ctx context.Context, conf *options.HttpServer, router *gin.Engine) error {
	// 初始化注册中心
	reg, err := registry.NewServiceRegistrar(&registry.Options{
		Endpoints:   []string{"localhost:2379"},
		ServiceAddr: conf.Addr,
		ServiceName: "pilotgo-server",
		Version:     "v3.0",
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %v", err)
	}
	// 创建 Gateway 实例
	gw := gateway.NewGateway(reg)

	// 启动 Gateway
	go func() {
		if err := gw.Run(router); err != nil {
			logger.Error("Gateway encountered an error: %v", err)
		}
	}()

	// 监听上下文取消信号
	go func() {
		<-ctx.Done()
		logger.Info("Gateway stopped successfully")

		// 停止 Gateway
		if err := gw.Stop(); err != nil {
			logger.Error("Failed to stop Gateway: %v", err)
		} else {
			logger.Info("Gateway stopped successfully")
		}
	}()
	return nil
}
