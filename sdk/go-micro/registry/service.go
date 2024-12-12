/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Thu Dec 12 17:36:05 2024 +0800
 */

package registry

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// ServiceRegistrar handles service registration and lifecycle management
type ServiceRegistrar struct {
	registry Registry
	info     *ServiceInfo
}

// NewServiceRegistrar creates a new service registrar
func NewServiceRegistrar(opts *Options) error {
	// Create registry client
	reg, err := NewRegistry(RegistryTypeEtcd, opts)
	if err != nil {
		return fmt.Errorf("failed to create registry: %v", err)
	}

	// Create service info
	info := &ServiceInfo{
		ServiceName: opts.ServiceName,
		Address:     strings.Split(opts.ServiceAddr, ":")[0],
		Port:        strings.Split(opts.ServiceAddr, ":")[1],
		Metadata:    map[string]string{"version": opts.Version},
	}

	sr := &ServiceRegistrar{
		registry: reg,
		info:     info,
	}

	if err := sr.Start(); err != nil {
		return err
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan // 等待信号

		if err := sr.Stop(); err != nil {
			logger.Error("failed to stop service registrar: %v", err)
		} else {
			logger.Info("service deregistered successfully")
		}
		os.Exit(0)
	}()

	logger.Info("service registered to etcd successfully")
	return nil
}

// Start registers the service and starts keeping it alive
func (s *ServiceRegistrar) Start() error {
	if err := s.registry.Register(s.info); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	return nil
}

// Stop deregisters the service and cleans up resources
func (s *ServiceRegistrar) Stop() error {
	if err := s.registry.Deregister(); err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}
	return s.registry.Close()
}
