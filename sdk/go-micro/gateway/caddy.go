/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 09 10:36:05 2025 +0800
 */
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/headers"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

type CaddyGateway struct {
	registry      registry.Registry
	services      map[string][]*registry.ServiceInfo
	serviceStatus map[string]bool
	serviceLock   *sync.RWMutex
	cancel        context.CancelFunc
	httpAddr      string
}

func NewCaddyGateway(reg registry.Registry, httpAddr string) *CaddyGateway {
	return &CaddyGateway{
		registry:      reg,
		services:      make(map[string][]*registry.ServiceInfo),
		serviceStatus: make(map[string]bool),
		serviceLock:   &sync.RWMutex{},
		httpAddr:      httpAddr,
	}
}
func (g *CaddyGateway) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 1)

	if err := caddy.Run(&caddy.Config{}); err != nil {
		return fmt.Errorf("failed to start caddy: %v", err)
	}
	logger.Info("start gateway service on: https://%s", g.httpAddr)

	// 初始加载已有服务
	if err := g.loadExistingServices(); err != nil {
		return err
	}

	// 启动服务监听
	go func() {
		if err := g.watchServices(ctx); err != nil {
			errChan <- err
		}
	}()
	select {
	case err := <-errChan:
		return fmt.Errorf("caddy gateway error: %v", err)
	case sig := <-sigChan:
		logger.Info("Received signal: %v, Shutting down caddy gateway...", sig)
		if err := g.Stop(); err != nil {
			return fmt.Errorf("error stopping caddy gateway: %v", err)
		}
	}
	return nil
}

func (g *CaddyGateway) Stop() error {
	if g.cancel != nil {
		g.cancel()
	}

	if err := g.registry.Close(); err != nil {
		logger.Error("failed to close registry: %v", err.Error())
	}

	if err := caddy.Stop(); err != nil {
		logger.Error("failed to stop caddy: %v", err.Error())
		return err
	}

	return nil
}

func (g *CaddyGateway) updateCaddyConfig() error {
	config, err := g.generateCaddyConfig()
	if err != nil {
		return fmt.Errorf("failed to generate caddy config: %v", err)
	}
	configBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal caddy config: %v", err)
	}

	if err := caddy.Load(configBytes, true); err != nil {
		return fmt.Errorf("failed to load caddy config: %v", err)
	}
	logger.Info("Caddy config reloaded successfully")
	return nil
}

func (g *CaddyGateway) generateCaddyConfig() (*caddy.Config, error) {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()

	var routes []caddyhttp.Route
	for serviceName, services := range g.services {
		if len(services) == 0 {
			continue
		}
		if !g.serviceStatus[serviceName] {
			continue
		}
		basePath := g.getBasePath(serviceName)
		matcherBytes, err := json.Marshal([]string{
			basePath,
			basePath + "/*",
			basePath + "/*/*",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal matcher config: %v", err)
		}
		matcherModuleMap := caddy.ModuleMap{
			"path": matcherBytes,
		}
		// 构造反向代理处理器配置
		upstreams := make([]map[string]interface{}, len(services))
		for i, s := range services {
			upstreams[i] = map[string]interface{}{
				"dial": fmt.Sprintf("%s:%s", s.Address, s.Port),
			}
			logger.Info("update upstream for %s: %s:%s", serviceName, s.Address, s.Port)
		}
		proxyHandlerConfig := map[string]interface{}{
			"handler":   "reverse_proxy",
			"upstreams": upstreams,
			"transport": map[string]interface{}{
				"protocol": "http",
				"tls":      nil,
			},
		}

		// 添加独立的 headers 处理器
		headersConfig := map[string]interface{}{
			"handler": "headers",
			"request": map[string]interface{}{
				"set": map[string][]string{
					"Host":               {"{http.reverse_proxy.upstream.host}"},
					"X-Forwarded-Prefix": {basePath},
				},
			},
			"response": map[string]interface{}{
				"set":      map[string][]string{},
				"deferred": true,
			},
		}

		// 创建路由
		route := caddyhttp.Route{
			MatcherSetsRaw: caddyhttp.RawMatcherSets{matcherModuleMap},
			HandlersRaw: []json.RawMessage{
				toRawMessage(headersConfig),
				toRawMessage(proxyHandlerConfig)},
		}
		routes = append(routes, route)
	}

	routes = append(routes, g.createMainServiceRoute())
	// 构造HTTP应用配置
	httpApp := map[string]interface{}{
		"servers": map[string]interface{}{
			"srv0": map[string]interface{}{
				"listen": []string{g.httpAddr},
				"routes": routes,
			},
		},
	}
	httpConfig, err := json.Marshal(httpApp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal http app config: %v", err)
	}

	return &caddy.Config{
		AppsRaw: map[string]json.RawMessage{
			"http": httpConfig,
		},
	}, nil
}

func toRawMessage(config map[string]interface{}) json.RawMessage {
	configBytes, _ := json.Marshal(config)
	return json.RawMessage(configBytes)
}

func (g *CaddyGateway) loadExistingServices() error {
	services, err := g.registry.List()
	if err != nil {
		return err
	}

	for _, service := range services {
		if service.ServiceName != "pilotgo-server" {
			g.addService(service)
		}
	}
	return g.updateCaddyConfig()
}

func (g *CaddyGateway) getBasePath(serviceName string) string {
	switch {
	case serviceName == "pilotgo-server":
		return "/"
	default:
		serviceType := strings.TrimSuffix(serviceName, "-service")
		return fmt.Sprintf("/plugin/%s", serviceType)
	}
}

func (g *CaddyGateway) createMainServiceRoute() caddyhttp.Route {
	matcherBytes, _ := json.Marshal([]string{
		"/*",
	})

	upstreams := make([]map[string]interface{}, 0)
	services, _ := g.registry.List()
	for _, service := range services {
		if service.ServiceName == "pilotgo-server" {
			upstreams = append(upstreams, map[string]interface{}{
				"dial": fmt.Sprintf("%s:%s", service.Address, service.Port),
			})
		}
	}

	return caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{
			caddy.ModuleMap{
				"path": matcherBytes,
			},
		},
		HandlersRaw: []json.RawMessage{
			toRawMessage(map[string]interface{}{
				"handler":   "reverse_proxy",
				"upstreams": upstreams,
				"headers": map[string]interface{}{
					"request": map[string]interface{}{
						"set": map[string][]string{
							"Host":               {"{http.reverse_proxy.upstream.host}"},
							"X-Forwarded-Prefix": {"/"},
						},
					},
				},
			}),
		},
	}
}
func (g *CaddyGateway) watchServices(ctx context.Context) error {
	callback := func(event registry.Event) {
		switch event.Type {
		case registry.EventTypePut:
			var service registry.ServiceInfo
			if err := json.Unmarshal([]byte(event.Value), &service); err != nil {
				logger.Error("Failed to unmarshal service info: %v", err.Error())
				return
			}
			g.addService(&service)
			logger.Info("Service added/updated: %s at %s:%s", service.ServiceName, service.Address, service.Port)
		case registry.EventTypeDelete:
			g.removeService(event.Key)
			logger.Info("Service removed: %s", event.Key)
		}

		if err := g.updateCaddyConfig(); err != nil {
			logger.Error("Failed to reload caddy config: %v", err.Error())
		}
	}

	return g.registry.Watch(ctx, "/services/", callback)
}
func (g *CaddyGateway) addService(service *registry.ServiceInfo) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	services := g.services[service.ServiceName]
	for i, s := range services {
		if s.Address == service.Address && s.Port == service.Port { // 检查是否已存在相同的服务实例, 如果存在则更新
			services[i] = service
			g.services[service.ServiceName] = services
			return
		}
	}
	g.services[service.ServiceName] = append(services, service)
	g.serviceStatus[service.ServiceName] = false
}

func (g *CaddyGateway) removeService(key string) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	for name, services := range g.services {
		for i, service := range services {
			if fmt.Sprintf("/services/%s", service.ServiceName) == key {
				g.services[name] = append(services[:i], services[i+1:]...)
				if len(g.services[name]) == 0 {
					delete(g.services, name)
					delete(g.serviceStatus, name)
				}
				return
			}
		}
	}
}
