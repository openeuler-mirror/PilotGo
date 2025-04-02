/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Dec 16 10:36:05 2024 +0800
 */
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// Gateway represents the API gateway
type Gateway struct {
	registry        registry.Registry
	services        map[string][]*registry.ServiceInfo
	serviceLock     *sync.RWMutex
	cancel          context.CancelFunc
	registeredPaths map[string]bool
}

// NewGateway creates a new API gateway
func NewGateway(reg registry.Registry) *Gateway {
	return &Gateway{
		registry:        reg,
		services:        make(map[string][]*registry.ServiceInfo),
		serviceLock:     &sync.RWMutex{},
		cancel:          nil,
		registeredPaths: make(map[string]bool),
	}
}

// Run starts the gateway and handles graceful shutdown
func (g *Gateway) Run(router *gin.Engine) error {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel // 保存cancel函数
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// 启动gateway
	go func() {
		if err := g.watchServices(ctx, router); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("gateway error: %v", err)
	case sig := <-sigChan:
		logger.Info("Received signal: %v, Shutting down gateway...", sig)
		if err := g.Stop(); err != nil {
			return fmt.Errorf("error stopping gateway: %v", err)
		}
	}

	return nil
}

// Stop stops the gateway
func (g *Gateway) Stop() error {
	if g.cancel != nil {
		g.cancel()
	}

	// 关闭注册中心连接
	if err := g.registry.Close(); err != nil {
		logger.Error("Registry close error: %v", err)
		return err
	}

	return nil
}

// watchServices watches for service changes in the registry
func (g *Gateway) watchServices(ctx context.Context, router *gin.Engine) error {
	// 新增：启动时加载已有子服务
	Services, err := g.registry.List()
	if err == nil {
		for _, service := range Services {
			if service.ServiceName != "pilotgo-server" {
				g.addService(service)
				g.updateRouter(router, service.ServiceName)

				logger.Info("Reload existing service: %s at %s:%s", service.ServiceName, service.Address, service.Port)
			}
		}
	} else {
		logger.Error("Failed to list existing services: %v", err)
	}

	callback := func(event registry.Event) {
		switch event.Type {
		case registry.EventTypePut:
			var service registry.ServiceInfo
			if err := json.Unmarshal([]byte(event.Value), &service); err != nil {
				logger.Error("Failed to unmarshal service info: %v\n", err)
				return
			}
			g.addService(&service)
			g.updateRouter(router, service.ServiceName)
			logger.Info("Service added/updated: %s at %s:%s\n", service.ServiceName, service.Address, service.Port)

		case registry.EventTypeDelete:
			// key格式: /services/{serviceName}
			parts := strings.Split(event.Key, "/")
			if len(parts) >= 3 {
				serviceName := parts[2]
				g.removeService(event.Key)
				g.updateRouter(router, serviceName)
				logger.Info("Service removed: %s", event.Key)
			} else {
				logger.Error("Invalid service key format: %s", event.Key)
			}
		}
	}

	return g.registry.Watch(ctx, "/services/", callback)
}

// addService adds a service to the gateway
func (g *Gateway) addService(service *registry.ServiceInfo) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	services := g.services[service.ServiceName]
	// Check if service already exists
	for i, s := range services {
		if s.Address == service.Address && s.Port == service.Port {
			services[i] = service
			return
		}
	}
	// Add new service
	g.services[service.ServiceName] = append(services, service)
}

// removeService removes a service from the gateway
func (g *Gateway) removeService(key string) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	for name, services := range g.services {
		for i, service := range services {
			if fmt.Sprintf("/services/%s", service.ServiceName) == key {
				g.services[name] = append(services[:i], services[i+1:]...)
				if len(g.services[name]) == 0 {
					delete(g.services, name)
				}
				return
			}
		}
	}
}
func (g *Gateway) getService(serviceName string) (*registry.ServiceInfo, error) {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()

	// 获取服务列表
	services, exists := g.services[serviceName]
	if !exists || len(services) == 0 {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	// 简单的负载均衡逻辑：轮询选择服务实例
	selectedService := services[0] // 默认选择第一个实例
	return selectedService, nil
}

func (g *Gateway) updateRouter(router *gin.Engine, serviceName string) {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	// 动态生成基础路径
	var basePath string
	switch {
	case serviceName == "pilotgo-server":
		basePath = "/api/v1"
	default:
		serviceType := strings.TrimSuffix(serviceName, "-service")
		basePath = fmt.Sprintf("/plugin/%s", serviceType)
	}
	path := fmt.Sprintf("%s/*path", basePath)

	// 动态绑定服务名到代理
	services, exists := g.services[serviceName]
	if !exists && len(services) == 0 { // 如果服务不存在且路径已注册，则标记为false
		g.registeredPaths[path] = false
		return
	}
	if exists && len(services) > 0 { // 子服务重启后，路由不需重新注册，可通过路由的状态来判断子服务是第一次注册还是重启注册，
		if _, exist := g.registeredPaths[path]; !exist { // 只有路径未注册时才实际注册路由处理函数
			router.Any(path, func(c *gin.Context) {
				targetService, err := g.getService(serviceName)
				if err != nil {
					logger.Error(fmt.Sprintf("Service %s unavailable", serviceName))
					response.Unavailable(c, nil, fmt.Sprintf("Service %s unavailable", serviceName))
					return
				}
				// 使用反向代理转发请求
				proxy := &httputil.ReverseProxy{
					Director: func(req *http.Request) {
						req.URL.Scheme = "http"
						req.URL.Host = fmt.Sprintf("%s:%s", targetService.Address, targetService.Port)
						req.URL.Path = basePath + c.Param("path")
						req.Header = c.Request.Header
					},
				}
				proxy.ServeHTTP(c.Writer, c.Request)
			})
			logger.Info("First time route registration for service: %s", serviceName)
		}
		g.registeredPaths[path] = true
	}
}
