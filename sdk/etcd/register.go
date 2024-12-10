/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 10:17:05 2024 +0800
 */
package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type RegisterOptions struct {
	Endpoints   []string
	ServiceName string
	ServiceAddr string
	Version     string
	DialTimeout time.Duration
}

type ServiceInfo struct {
	ServiveName string            `json:"serviceName"`
	Address     string            `json:"address"`
	Port        string            `json:"port"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type ServiceRegister struct {
	client        *Client
	leaseID       clientv3.LeaseID
	servicePath   string
	serviceInfo   *ServiceInfo
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	cancel        context.CancelFunc // To cancel the keep-alive goroutine
}

// SetupEtcdRegistration initializes etcd registration with graceful shutdown
func Register(ctx context.Context, opts *RegisterOptions) error {
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
			serviceRegister.deregister()
		}
	}()

	return nil
}

// Deregister removes the service from etcd

// GetAllServices returns all registered services
func (c *Client) GetAllServices() (map[string][]*ServiceInfo, error) {
	services := make(map[string][]*ServiceInfo)

	// Get all services with prefix "/services/"
	resp, err := c.client.Get(context.Background(), "/services/", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	// Process each service
	for _, kv := range resp.Kvs {
		var serviceInfo ServiceInfo
		err := json.Unmarshal(kv.Value, &serviceInfo)
		if err != nil {
			fmt.Printf("Failed to unmarshal service info: %v\n", err)
			continue
		}

		services[serviceInfo.ServiveName] = append(services[serviceInfo.ServiveName], &serviceInfo)
	}

	return services, nil
}

// RegisterToEtcd registers the server to etcd
func registerService(opts *RegisterOptions) (*ServiceRegister, error) {
	// Create etcd client
	etcdClient, err := NewClient(opts.Endpoints, opts.DialTimeout)
	if err != nil {
		return nil, errors.New("failed to create etcd client: " + err.Error())
	}

	// Create service info
	serviceInfo := &ServiceInfo{
		ServiveName: opts.ServiceName,
		Address:     strings.Split(opts.ServiceAddr, ":")[0],
		Port:        strings.Split(opts.ServiceAddr, ":")[1],
		Metadata: map[string]string{
			"version": opts.Version,
		},
	}

	// Register service with TTL
	serviceRegister, err := newServiceRegister(etcdClient, serviceInfo, 10)
	if err != nil {
		etcdClient.Close()
		return nil, errors.New("failed to register service: " + err.Error())
	}
	logger.Info("service registered to etcd successfully")

	return serviceRegister, nil
}

// NewServiceRegister creates a new service register
func newServiceRegister(client *Client, info *ServiceInfo, ttl int64) (*ServiceRegister, error) {
	// 首先检查客户端是否可用
	if client == nil || client.client == nil {
		return nil, errors.New("etcd client is not initialized")
	}

	ctx, cancel := context.WithCancel(client.ctx)
	sr := &ServiceRegister{
		client:      client,
		serviceInfo: info,
		servicePath: fmt.Sprintf("/services/%s", info.ServiveName),
		cancel:      cancel,
	}

	// 检查etcd连接状态
	ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 尝试简单的etcd操作来验证连接
	_, err := client.client.Get(ctx, "/health", clientv3.WithCountOnly())
	if err != nil {
		return nil, err
	}

	var opts []clientv3.OpOption
	if ttl > 0 {
		// Create lease
		resp, err := client.client.Grant(ctx, ttl)
		if err != nil {
			return nil, err
		}
		sr.leaseID = resp.ID

		// Keep lease alive
		keepAliveChan, err := client.client.KeepAlive(context.Background(), resp.ID)
		if err != nil {
			return nil, err
		}
		sr.keepAliveChan = keepAliveChan

		opts = append(opts, clientv3.WithLease(sr.leaseID))
	}

	// Register service
	if err := sr.register(opts...); err != nil {
		return nil, errors.New("failed to register service:" + err.Error())
	}

	// Start keepalive goroutine if using TTL
	if ttl > 0 {
		go sr.keepAlive()
	}

	return sr, nil
}

// register puts service info into etcd
func (sr *ServiceRegister) register(opts ...clientv3.OpOption) error {
	value, err := json.Marshal(sr.serviceInfo)
	if err != nil {
		return err
	}

	_, err = sr.client.client.Put(
		context.Background(),
		sr.servicePath,
		string(value),
		opts...,
	)
	return err
}
func (sr *ServiceRegister) deregister() error {
	// Cancel the context to stop keep-alive goroutine
	if sr.cancel != nil {
		sr.cancel()
	}

	// Delete the service key from etcd
	_, err := sr.client.client.Delete(
		context.Background(),
		sr.servicePath,
	)
	if err != nil {
		return fmt.Errorf("failed to deregister service %s: %v", sr.serviceInfo.ServiveName, err)
	}

	// Revoke the lease if it exists
	if sr.leaseID != 0 {
		_, err = sr.client.client.Revoke(context.Background(), sr.leaseID)
		if err != nil {
			return fmt.Errorf("failed to revoke lease for service %s: %v", sr.serviceInfo.ServiveName, err)
		}
	}

	return nil
}

// keepAlive keeps the lease alive
func (sr *ServiceRegister) keepAlive() {
	for {
		select {
		case resp, ok := <-sr.keepAliveChan:
			if !ok {
				fmt.Printf("Keep-alive channel closed for service %s/%s\n", sr.serviceInfo.ServiveName, sr.serviceInfo.ServiveName)
				return
			}
			if resp == nil {
				fmt.Printf("Keep-alive response is nil for service %s/%s\n", sr.serviceInfo.ServiveName, sr.serviceInfo.ServiveName)
				return
			}
		case <-sr.client.ctx.Done():
			fmt.Printf("Context done for service %s/%s\n", sr.serviceInfo.ServiveName, sr.serviceInfo.ServiveName)
			return
		}
	}
}
