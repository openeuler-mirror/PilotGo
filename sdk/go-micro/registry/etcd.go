/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Thu Dec 12 17:36:05 2024 +0800
 */
package registry

import (
	"context"
	"encoding/json"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdRegistry struct {
	client      *clientv3.Client
	leaseID     clientv3.LeaseID
	ctx         context.Context
	servicePath string
	serviceInfo *ServiceInfo
	keepAlive   <-chan *clientv3.LeaseKeepAliveResponse
	options     *Options
	cancel      context.CancelFunc
}

func newEtcdRegistry(opts *Options) (Registry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opts.Endpoints,
		DialTimeout: opts.DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &etcdRegistry{
		client:  cli,
		ctx:     ctx,
		cancel:  cancel,
		options: opts,
	}, nil
}

func (e *etcdRegistry) Register(info *ServiceInfo) error {
	e.serviceInfo = info

	// Create lease
	lease, err := e.client.Grant(e.ctx, 10)
	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}
	e.leaseID = lease.ID

	key := fmt.Sprintf("/services/%s", info.ServiceName)
	value, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %v", err)
	}
	e.servicePath = key

	// Register service with lease
	_, err = e.client.Put(e.ctx, key, string(value), clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to put service info: %v", err)
	}

	// Keep lease alive
	keepAliveChan, err := e.client.KeepAlive(e.ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("failed to keep lease alive: %v", err)
	}
	e.keepAlive = keepAliveChan

	go e.keepAliveLoop()

	return nil
}

func (e *etcdRegistry) Deregister() error {
	// Delete the service key from etcd
	_, err := e.client.Delete(e.ctx, e.servicePath)
	if err != nil {
		return fmt.Errorf("failed to deregister service %s: %v", e.serviceInfo.ServiceName, err)
	}

	// Revoke the lease if it exists
	if e.leaseID != 0 {
		_, err := e.client.Revoke(e.ctx, e.leaseID)
		if err != nil {
			return fmt.Errorf("failed to revoke lease: %v", err)
		}
	}

	e.cancel()
	return nil
}

func (e *etcdRegistry) Get(key string) (*ServiceInfo, error) {
	resp, err := e.client.Get(e.ctx, "/services/"+key)
	if err != nil {
		return &ServiceInfo{}, err
	}
	if len(resp.Kvs) == 0 {
		return &ServiceInfo{}, fmt.Errorf("未找到服务实例")
	}
	var service ServiceInfo
	if err := json.Unmarshal(resp.Kvs[0].Value, &service); err != nil {
		return &ServiceInfo{}, err
	}
	return &service, nil
}

func (e *etcdRegistry) Put(key string, value string) error {
	_, err := e.client.Put(e.ctx, key, value)
	return err
}

func (r *etcdRegistry) List() ([]*ServiceInfo, error) {
	resp, err := r.client.Get(context.Background(), "/services/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	services := make([]*ServiceInfo, 0)
	for _, kv := range resp.Kvs {
		var service ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err == nil {
			services = append(services, &service)
		}
	}
	return services, nil
}
func (e *etcdRegistry) Delete(key string) error {
	_, err := e.client.Delete(e.ctx, key)
	return err
}

func (e *etcdRegistry) Close() error {
	e.cancel()
	return e.client.Close()
}

func (e *etcdRegistry) Watch(ctx context.Context, key string, callback WatchCallback) error {
	watchChan := e.client.Watch(e.ctx, key, clientv3.WithPrefix())
	go func() {
		for resp := range watchChan {
			for _, ev := range resp.Events {
				var eventType EventType
				switch ev.Type {
				case clientv3.EventTypePut:
					eventType = EventTypePut
				case clientv3.EventTypeDelete:
					eventType = EventTypeDelete
				}
				callback(Event{
					Type:  eventType,
					Key:   string(ev.Kv.Key),
					Value: string(ev.Kv.Value),
				})
			}
		}
	}()
	return nil
}

func (e *etcdRegistry) keepAliveLoop() {
	for {
		select {
		case resp, ok := <-e.keepAlive:
			if !ok {
				return
			}
			if resp == nil {
				return
			}
		case <-e.ctx.Done():
			return
		}
	}
}
