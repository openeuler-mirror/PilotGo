/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Dec 16 10:36:05 2024 +0800
 */
package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
)

// Proxy represents the API proxy
type Proxy struct {
	registry    registry.Registry
	services    map[string][]*registry.ServiceInfo
	serviceLock sync.RWMutex
	balancer    LoadBalancer
	retries     int
}

// NewProxy creates a new API proxy
func NewProxy(reg registry.Registry) *Proxy {
	return &Proxy{
		registry: reg,
		services: make(map[string][]*registry.ServiceInfo),
		balancer: NewRoundRobinBalancer(),
		retries:  3, // Default number of retries
	}
}

// Start starts the proxy server
func (p *Proxy) Start(addr string) error {
	// Watch for service changes
	if err := p.watchServices(); err != nil {
		return fmt.Errorf("failed to start service watcher: %v", err)
	}

	log.Printf("Starting proxy server on %s\n", addr)
	http.HandleFunc("/", p.ProxyHandler)
	return http.ListenAndServe(addr, nil)
}

// watchServices watches for service changes in the registry
func (p *Proxy) watchServices() error {
	callback := func(event registry.Event) {
		switch event.Type {
		case registry.EventTypePut:
			var service registry.ServiceInfo
			if err := json.Unmarshal([]byte(event.Value), &service); err != nil {
				log.Printf("Failed to unmarshal service info: %v\n", err)
				return
			}
			p.addService(&service)
			log.Printf("Service added/updated: %s at %s:%s\n", service.ServiceName, service.Address, service.Port)

		case registry.EventTypeDelete:
			p.removeService(event.Key)
			log.Printf("Service removed: %s\n", event.Key)
		}
	}

	return p.registry.Watch(context.Background(), "/services/", callback)
}

// addService adds a service to the proxy
func (p *Proxy) addService(service *registry.ServiceInfo) {
	p.serviceLock.Lock()
	defer p.serviceLock.Unlock()

	services := p.services[service.ServiceName]
	// Check if service already exists
	for i, s := range services {
		if s.Address == service.Address && s.Port == service.Port {
			services[i] = service
			return
		}
	}
	// Add new service
	p.services[service.ServiceName] = append(services, service)
}

// removeService removes a service from the proxy
func (p *Proxy) removeService(key string) {
	p.serviceLock.Lock()
	defer p.serviceLock.Unlock()

	for name, services := range p.services {
		for i, service := range services {
			if fmt.Sprintf("/services/%s", service.ServiceName) == key {
				p.services[name] = append(services[:i], services[i+1:]...)
				if len(p.services[name]) == 0 {
					delete(p.services, name)
				}
				return
			}
		}
	}
}

// ProxyHandler handles the proxying of requests to services
func (p *Proxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.Header.Get("X-Service-Name")
	if serviceName == "" {
		http.Error(w, "Service name not specified", http.StatusBadRequest)
		return
	}

	var lastErr error
	for i := 0; i <= p.retries; i++ {
		service, err := p.getService(serviceName)
		if err != nil {
			lastErr = err
			continue
		}

		target := fmt.Sprintf("http://%s:%s", service.Address, service.Port)
		targetURL, err := url.Parse(target)
		if err != nil {
			lastErr = err
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Add custom error handling
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			lastErr = err
			// Don't write error here, let retry logic handle it
		}

		// Add custom request modification
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Header.Set("X-Real-IP", req.RemoteAddr)
			req.Header.Set("X-Retry-Count", fmt.Sprintf("%d", i))
		}

		// Attempt to proxy the request
		var proxyErr error
		done := make(chan bool)
		go func() {
			proxy.ServeHTTP(w, r)
			done <- true
		}()

		select {
		case <-done:
			return // Request completed successfully
		case <-time.After(10 * time.Second):
			proxyErr = fmt.Errorf("request timed out")
		}

		if proxyErr != nil {
			lastErr = proxyErr
			continue
		}

		return // Request completed successfully
	}

	// All retries failed
	http.Error(w, fmt.Sprintf("Service unavailable after %d retries: %v", p.retries, lastErr), http.StatusServiceUnavailable)
}

// getService returns a service instance using the configured load balancer
func (p *Proxy) getService(name string) (*registry.ServiceInfo, error) {
	p.serviceLock.RLock()
	defer p.serviceLock.RUnlock()

	services, ok := p.services[name]
	if !ok || len(services) == 0 {
		return nil, fmt.Errorf("no available services for %s", name)
	}

	return p.balancer.Next(services), nil
}
