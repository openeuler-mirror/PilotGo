/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 14:36:05 2024 +0800
 */
package gateway

import (
	"sync"
	"sync/atomic"

	"gitee.com/openeuler/PilotGo/sdk/etcd"
)

// LoadBalancer interface defines methods for load balancing
type LoadBalancer interface {
	Next(services []*etcd.ServiceInfo) *etcd.ServiceInfo
}

// RoundRobinBalancer implements round-robin load balancing
type RoundRobinBalancer struct {
	counter uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (b *RoundRobinBalancer) Next(services []*etcd.ServiceInfo) *etcd.ServiceInfo {
	if len(services) == 0 {
		return nil
	}
	count := atomic.AddUint64(&b.counter, 1)
	return services[count%uint64(len(services))]
}

// WeightedRoundRobinBalancer implements weighted round-robin load balancing
type WeightedRoundRobinBalancer struct {
	mu      sync.Mutex
	current int
	weights map[string]int
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
	return &WeightedRoundRobinBalancer{
		weights: make(map[string]int),
	}
}

func (b *WeightedRoundRobinBalancer) SetWeight(serviceID string, weight int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.weights[serviceID] = weight
}

func (b *WeightedRoundRobinBalancer) Next(services []*etcd.ServiceInfo) *etcd.ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Default weight is 1 if not specified
	totalWeight := 0
	for _, service := range services {
		weight := b.weights[service.ServiveName]
		if weight == 0 {
			weight = 1
		}
		totalWeight += weight
	}

	b.current = (b.current + 1) % totalWeight
	for _, service := range services {
		weight := b.weights[service.ServiveName]
		if weight == 0 {
			weight = 1
		}
		if b.current < weight {
			return service
		}
		b.current -= weight
	}

	return services[0]
}

// LeastConnectionBalancer implements least connection load balancing
type LeastConnectionBalancer struct {
	mu          sync.RWMutex
	connections map[string]int
}

func NewLeastConnectionBalancer() *LeastConnectionBalancer {
	return &LeastConnectionBalancer{
		connections: make(map[string]int),
	}
}

func (b *LeastConnectionBalancer) IncrementConnections(serviceID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.connections[serviceID]++
}

func (b *LeastConnectionBalancer) DecrementConnections(serviceID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.connections[serviceID] > 0 {
		b.connections[serviceID]--
	}
}

func (b *LeastConnectionBalancer) Next(services []*etcd.ServiceInfo) *etcd.ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	var minConnections int = -1
	var selectedService *etcd.ServiceInfo

	for _, service := range services {
		connections := b.connections[service.ServiveName]
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selectedService = service
		}
	}

	return selectedService
}
