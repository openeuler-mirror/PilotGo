/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 14:36:05 2024 +0800
 */
package etcd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitee.com/openeuler/PilotGo/sdk/etcd/client"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// SetupEtcdRegistration initializes etcd registration with graceful shutdown
func Register(ctx context.Context, opts *client.RegisterOptions) error {
	serviceRegister, err := registerService(opts)
	if err != nil {
		return err
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
		case <-ctx.Done():
		}
		if serviceRegister != nil {
			serviceRegister.Deregister()
		}
	}()

	return nil
}

// RegisterToEtcd registers the server to etcd
func registerService(opts *client.RegisterOptions) (*client.ServiceRegister, error) {
	// 1. Create etcd client
	etcdClient, err := client.NewClient(opts.Endpoints, opts.DialTimeout)
	if err != nil {
		return nil, errors.New("failed to create etcd client: " + err.Error())
	}

	serviceInfo := &client.ServiceInfo{
		ServiveName: opts.ServiceName,
		Address:     strings.Split(opts.ServiceAddr, ":")[0],
		Port:        strings.Split(opts.ServiceAddr, ":")[1],
		Metadata: map[string]string{
			"version": opts.Version,
		},
	}

	// 2. Register service with TTL
	serviceRegister, err := client.NewServiceRegister(etcdClient, serviceInfo, 10)
	if err != nil {
		etcdClient.Close()
		return nil, errors.New("failed to register service: " + err.Error())
	}
	logger.Info("service registered to etcd successfully")

	return serviceRegister, nil
}
