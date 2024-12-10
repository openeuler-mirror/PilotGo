/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 13:56:05 2024 +0800
 */
package client

import (
	"context"
	"encoding/json"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EventType represents the type of event
type EventType int32

const (
	EventTypePut    EventType = 0
	EventTypeDelete EventType = 1
)

type WatchCallback func(eventType EventType, key, value string)

type Watcher struct {
	client   *Client
	prefix   string
	callback WatchCallback
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewWatcher creates a new watcher
func NewWatcher(client *Client, prefix string, callback WatchCallback) *Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Watcher{
		client:   client,
		prefix:   prefix,
		callback: callback,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start starts watching for changes
func (w *Watcher) Start() {
	watchChan := w.client.client.Watch(w.ctx, w.prefix, clientv3.WithPrefix())
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				var eventType EventType
				switch event.Type {
				case clientv3.EventTypePut:
					eventType = EventTypePut
				case clientv3.EventTypeDelete:
					eventType = EventTypeDelete
				}
				w.callback(eventType, string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}()
}

// Stop stops watching
func (w *Watcher) Stop() {
	w.cancel()
}

// WatchService watches for service changes
func WatchService(client *Client, serviceName string) (*Watcher, error) {
	servicePath := fmt.Sprintf("/services/%s/", serviceName)
	services := make(map[string]*ServiceInfo)

	callback := func(eventType EventType, key, value string) {
		switch eventType {
		case EventTypePut:
			var service ServiceInfo
			if err := json.Unmarshal([]byte(value), &service); err != nil {
				fmt.Printf("Failed to unmarshal service info: %v\n", err)
				return
			}
			services[key] = &service
			fmt.Printf("Service added/updated: %s at %s:%v\n", service.ServiveName, service.Address, service.Port)

		case EventTypeDelete:
			if service, ok := services[key]; ok {
				delete(services, key)
				fmt.Printf("Service removed: %s\n", service.ServiveName)
			}
		}
	}

	watcher := NewWatcher(client, servicePath, callback)
	watcher.Start()
	return watcher, nil
}

// WatchConfig watches for configuration changes
func WatchConfig(client *Client, prefix string) (*Watcher, error) {
	callback := func(eventType EventType, key, value string) {
		switch eventType {
		case EventTypePut:
			fmt.Printf("Configuration updated - Key: %s, Value: %s\n", key, value)
		case EventTypeDelete:
			fmt.Printf("Configuration deleted - Key: %s\n", key)
		}
	}

	watcher := NewWatcher(client, prefix, callback)
	watcher.Start()
	return watcher, nil
}
