/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Dec 10 13:56:05 2024 +0800
 */
package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/etcd/client"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Proxy represents the API proxy
type Proxy struct {
	client      *clientv3.Client
	services    map[string][]*client.ServiceInfo
	serviceLock sync.RWMutex
}

// NewProxy creates a new API proxy
func NewProxy(cli *clientv3.Client) *Proxy {
	return &Proxy{
		client:   cli,
		services: make(map[string][]*client.ServiceInfo),
	}
}

// Start starts the proxy server
func (p *Proxy) Start(addr string) error {
	http.HandleFunc("/", p.ProxyHandler)
	return http.ListenAndServe(addr, nil)
}

// ProxyHandler handles the proxying of requests to services
func (p *Proxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.Header.Get("X-Service-Name")
	if serviceName == "" {
		http.Error(w, "Service name not specified", http.StatusBadRequest)
		return
	}

	service, err := p.getService(serviceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting service: %v", err), http.StatusInternalServerError)
		return
	}

	targetURL, err := url.Parse(fmt.Sprintf("http://%s:%s", service.Address, service.Port))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing service URL: %v", err), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

// getService returns a service instance using round-robin load balancing
func (p *Proxy) getService(name string) (*client.ServiceInfo, error) {
	p.serviceLock.RLock()
	defer p.serviceLock.RUnlock()

	services, ok := p.services[name]
	if !ok || len(services) == 0 {
		return nil, fmt.Errorf("no available services for %s", name)
	}

	// Simple round-robin for now
	service := services[0]
	// Move the first service to the end
	p.services[name] = append(services[1:], service)

	return service, nil
}
