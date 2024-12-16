/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Dec 16 10:36:05 2024 +0800
 */
package proxy

import (
	"sync"
	"sync/atomic"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
)

// LoadBalancer interface defines methods for load balancing
type LoadBalancer interface {
	Next(services []*registry.ServiceInfo) *registry.ServiceInfo
	UpdateWeight(address string, weight int)
	GetWeight(address string) int
}

// RoundRobinBalancer implements round-robin load balancing
type RoundRobinBalancer struct {
	counter uint64
	weights map[string]int
	mu      sync.RWMutex
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		weights: make(map[string]int),
	}
}

func (b *RoundRobinBalancer) UpdateWeight(address string, weight int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.weights[address] = weight
}

func (b *RoundRobinBalancer) GetWeight(address string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.weights[address]
}

func (b *RoundRobinBalancer) Next(services []*registry.ServiceInfo) *registry.ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// Get total weight
	totalWeight := 0
	for _, service := range services {
		weight := b.weights[service.Address]
		if weight == 0 {
			weight = 1
		}
		totalWeight += weight
	}

	// Use atomic counter for thread safety
	count := atomic.AddUint64(&b.counter, 1)
	targetWeight := count % uint64(totalWeight)

	// Find service based on weight
	currentWeight := 0
	for _, service := range services {
		weight := b.weights[service.Address]
		if weight == 0 {
			weight = 1
		}
		currentWeight += weight
		if uint64(currentWeight) > targetWeight {
			return service
		}
	}

	return services[0]
}

// WeightedRoundRobinBalancer implements weighted round-robin load balancing
type WeightedRoundRobinBalancer struct {
	mu      sync.RWMutex
	current int
	weights map[string]int
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
	return &WeightedRoundRobinBalancer{
		weights: make(map[string]int),
	}
}

func (b *WeightedRoundRobinBalancer) UpdateWeight(address string, weight int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.weights[address] = weight
}

func (b *WeightedRoundRobinBalancer) GetWeight(address string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.weights[address]
}

func (b *WeightedRoundRobinBalancer) Next(services []*registry.ServiceInfo) *registry.ServiceInfo {
	if len(services) == 0 {
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Simple weighted round-robin implementation
	totalWeight := 0
	maxWeight := 0
	for _, service := range services {
		weight := b.weights[service.Address]
		if weight == 0 {
			weight = 1
		}
		if weight > maxWeight {
			maxWeight = weight
		}
		totalWeight += weight
	}

	for {
		b.current = (b.current + 1) % maxWeight
		for _, service := range services {
			weight := b.weights[service.Address]
			if weight == 0 {
				weight = 1
			}
			if weight >= b.current {
				return service
			}
		}
		if b.current >= maxWeight {
			b.current = 0
		}
	}
}
