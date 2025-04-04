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
	"errors"
	"time"
)

// ServiceInfo represents the basic information of a service
type ServiceInfo struct {
	ServiceName string            `json:"serviceName"`
	Address     string            `json:"address"`
	Port        string            `json:"port"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Options represents the configuration options for service registry
type Options struct {
	Endpoints   []string      // etcd address
	ServiceName string        // 服务地址
	ServiceAddr string        // 服务名称
	Version     string        // 服务版本
	DialTimeout time.Duration // 超时时间
}

// EventType represents the type of service registry events
type EventType int32

const (
	EventTypePut    EventType = 0
	EventTypeDelete EventType = 1
)

// Event represents a service registry event
type Event struct {
	Type  EventType
	Key   string
	Value string
}

// WatchCallback is the callback function for watch events
type WatchCallback func(event Event)

// Registry defines the interface for service registry operations
type Registry interface {
	// Register registers a service
	Register(info *ServiceInfo) error
	// Deregister removes a service registration
	Deregister() error
	// Get retrieves service information
	Get(key string) (string, error)
	// Put stores service information
	Put(key string, value string) error
	// Delete removes service information
	Delete(key string) error
	// List retrieves all service information
	List() ([]*ServiceInfo, error)
	// Watch watches for service changes
	Watch(ctx context.Context, key string, callback WatchCallback) error
	// Close closes the registry client
	Close() error
}

// RegistryType represents the type of registry
type RegistryType string

const (
	RegistryTypeEtcd RegistryType = "etcd"
)

// NewRegistry creates a new registry client based on the registry type
func NewRegistry(registryType RegistryType, opts *Options) (Registry, error) {
	switch registryType {
	case RegistryTypeEtcd:
		return newEtcdRegistry(opts)
	default:
		return nil, errors.New("unsupported registry type")
	}
}
