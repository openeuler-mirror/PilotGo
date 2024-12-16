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
	"sync"
	"syscall"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/proxy"
	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

// Gateway represents the API gateway
type Gateway struct {
	registry    registry.Registry
	services    map[string][]*registry.ServiceInfo
	serviceLock sync.RWMutex
	balancer    proxy.LoadBalancer
	proxy       *proxy.Proxy
	server      *http.Server
	cancel      context.CancelFunc
}

// NewGateway creates a new API gateway
func NewGateway(reg registry.Registry) *Gateway {
	return &Gateway{
		registry: reg,
		services: make(map[string][]*registry.ServiceInfo),
		balancer: proxy.NewRoundRobinBalancer(),
		proxy:    proxy.NewProxy(reg),
	}
}

// Run starts the gateway and handles graceful shutdown
func (g *Gateway) Run(addr string, router *gin.Engine) error {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// 启动gateway
	go func() {
		if err := g.Start(addr, router); err != nil {
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

// Start starts the gateway server
func (g *Gateway) Start(addr string, router *gin.Engine) error {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel

	// Start watching for service changes
	if err := g.watchServices(router); err != nil {
		return fmt.Errorf("failed to start service watcher: %v", err)
	}

	g.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Info("Starting gateway server on %s\n", addr)
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	return nil
}

// Stop stops the gateway
func (g *Gateway) Stop() error {
	if g.cancel != nil {
		g.cancel()
	}

	// 创建一个带超时的上下文用于关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务器
	if g.server != nil {
		if err := g.server.Shutdown(ctx); err != nil {
			logger.Error("HTTP server Shutdown error: %v", err)
		}
	}

	// 关闭注册中心连接
	if err := g.registry.Close(); err != nil {
		logger.Error("Registry close error: %v", err)
		return err
	}

	return nil
}

// watchServices watches for service changes in the registry
func (g *Gateway) watchServices(router *gin.Engine) error {
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
			g.removeService(event.Key)
			logger.Info("Service removed: %s\n", event.Key)
		}
	}

	return g.registry.Watch("/services/", callback)
}

// updateRouter updates the router with the given service's routes
func (g *Gateway) updateRouter(router *gin.Engine, serviceName string) {
	// 动态绑定服务名到代理
	router.Any(fmt.Sprintf("/%s/*path", serviceName), func(c *gin.Context) {
		targetService, err := g.getService(serviceName)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("Service %s unavailable", serviceName)})
			return
		}

		// 构造目标URL
		targetURL := fmt.Sprintf("http://%s:%s%s", targetService.Address, targetService.Port, c.Param("path"))
		logger.Info("Proxying request to: %s", targetURL)

		// 使用反向代理转发请求
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = fmt.Sprintf("%s:%s", targetService.Address, targetService.Port)
				req.URL.Path = c.Param("path")
				req.Header = c.Request.Header
			},
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})
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
			if fmt.Sprintf("/services/%s/%s:%s", service.ServiceName, service.Address, service.Port) == key {
				g.services[name] = append(services[:i], services[i+1:]...)
				if len(g.services[name]) == 0 {
					delete(g.services, name)
				}
				return
			}
		}
	}
}

// getService returns a service instance using the load balancer
func (g *Gateway) getService(name string) (*registry.ServiceInfo, error) {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()

	services := g.services[name]
	if len(services) == 0 {
		return nil, fmt.Errorf("no available services for %s", name)
	}

	return g.balancer.Next(services), nil
}
