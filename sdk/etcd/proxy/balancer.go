/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 14:36:05 2024 +0800
 */
package proxy

import (
	"sync"
	"sync/atomic"

	"gitee.com/openeuler/PilotGo/sdk/etcd/client"
)

// LoadBalancer interface defines methods for load balancing
type LoadBalancer interface {
	Next(services []*client.ServiceInfo) *client.ServiceInfo
}

// RoundRobinBalancer implements round-robin load balancing
type RoundRobinBalancer struct {
	counter uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (b *RoundRobinBalancer) Next(services []*client.ServiceInfo) *client.ServiceInfo {
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

func (b *WeightedRoundRobinBalancer) Next(services []*client.ServiceInfo) *client.ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Simple weighted round-robin implementation
	totalWeight := 0
	for _, service := range services {
		weight := b.weights[service.Address]
		if weight == 0 {
			weight = 1
		}
		totalWeight += weight
	}

	b.current = (b.current + 1) % totalWeight
	for _, service := range services {
		weight := b.weights[service.Address]
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
