/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 13:56:05 2024 +0800
 */
package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/etcd/client"
)

// Gateway represents the API gateway
type Gateway struct {
	etcdClient  *client.Client
	services    map[string][]*client.ServiceInfo
	serviceLock sync.RWMutex
	watcher     *client.Watcher
}

// NewGateway creates a new API gateway
func NewGateway(etcdClient *client.Client) *Gateway {
	g := &Gateway{
		etcdClient: etcdClient,
		services:   make(map[string][]*client.ServiceInfo),
	}

	// Start watching for service changes
	g.watchServices()
	return g
}

// watchServices watches for service changes in etcd
func (g *Gateway) watchServices() {
	callback := func(eventType client.EventType, key, value string) {
		switch eventType {
		case client.EventTypePut:
			var service client.ServiceInfo
			if err := json.Unmarshal([]byte(value), &service); err != nil {
				fmt.Printf("Failed to unmarshal service info: %v\n", err)
				return
			}
			g.addService(&service)

		case client.EventTypeDelete:
			g.removeService(key)
		}
	}

	g.watcher = client.NewWatcher(g.etcdClient, "/services/", callback)
	g.watcher.Start()
}

// addService adds a service to the gateway
func (g *Gateway) addService(service *client.ServiceInfo) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	services := g.services[service.ServiveName]
	// Check if service already exists
	for i, s := range services {
		if s.Address == service.Address {
			services[i] = service
			return
		}
	}
	// Add new service
	g.services[service.ServiveName] = append(services, service)
}

// removeService removes a service from the gateway
func (g *Gateway) removeService(key string) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	for name, services := range g.services {
		for i, service := range services {
			if fmt.Sprintf("/services/%s", service.ServiveName) == key {
				g.services[name] = append(services[:i], services[i+1:]...)
				return
			}
		}
	}
}

// getService returns a service instance using round-robin load balancing
func (g *Gateway) getService(name string) (*client.ServiceInfo, error) {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()

	services := g.services[name]
	if len(services) == 0 {
		return nil, fmt.Errorf("no available services for %s", name)
	}

	// Simple round-robin load balancing
	service := services[time.Now().UnixNano()%int64(len(services))]
	return service, nil
}

// ProxyHandler handles the proxying of requests to services
func (g *Gateway) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract service name from path
	serviceName := r.Header.Get("X-Service-Name")
	if serviceName == "" {
		http.Error(w, "Service name not specified", http.StatusBadRequest)
		return
	}

	service, err := g.getService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Create the target URL
	target := fmt.Sprintf("http://%s:%v", service.Address, service.Port)
	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Update the headers to allow for SSL redirection
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	proxy.ServeHTTP(w, r)
}

// Start starts the gateway server
func (g *Gateway) Start(addr string) error {
	http.HandleFunc("/", g.ProxyHandler)
	return http.ListenAndServe(addr, nil)
}

// Stop stops the gateway
func (g *Gateway) Stop() {
	if g.watcher != nil {
		g.watcher.Stop()
	}
}
