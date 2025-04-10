package gateway

import (
	"errors"
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/go-micro/registry"
)

// 设置服务状态
func (g *CaddyGateway) SetServiceStatus(serviceName string, status bool) error {
	if _, ok := g.serviceStatus[serviceName]; !ok {
		return errors.New("service not found")
	}
	g.serviceLock.Lock()
	g.serviceStatus[serviceName] = status
	g.serviceLock.Unlock()

	if err := g.updateCaddyConfig(); err != nil {
		return fmt.Errorf("failed to reload caddy config: %w", err)
	}
	return nil
}

// 获取服务状态
func (g *CaddyGateway) GetServiceStatus(serviceName string) bool {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()
	return g.serviceStatus[serviceName]
}

// 获取某个服务
func (g *CaddyGateway) GetService(key string) *registry.ServiceInfo {
	g.serviceLock.Lock()
	defer g.serviceLock.Unlock()

	var service *registry.ServiceInfo
	for _, services := range g.services {
		for _, s := range services {
			if fmt.Sprintf("/services/%s", s.ServiceName) == key {
				service = s
				break
			}
		}
	}
	return service
}

// 获取所有服务
func (g *CaddyGateway) GetAllServices() []map[string]interface{} {
	g.serviceLock.RLock()
	defer g.serviceLock.RUnlock()

	result := make([]map[string]interface{}, 0)

	for serviceName, services := range g.services {
		if len(services) == 0 {
			continue
		}

		svc := services[0] // 只取第一个实例
		result = append(result, map[string]interface{}{
			"serviceName": svc.ServiceName,
			"address":     svc.Address,
			"port":        svc.Port,
			"version":     svc.Metadata["version"],
			"status":      g.serviceStatus[serviceName],
			"extentions":  svc.Metadata["extentions"],
			"permissions": svc.Metadata["permissions"],
		})
	}

	return result
}
