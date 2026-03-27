package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// ConsulDriver Consul 服务发现驱动
type ConsulDriver struct {
	cfg    Config
	client *consul.Client
	mu     sync.RWMutex
}

// NewConsulDriver 创建 Consul 驱动实例
func NewConsulDriver(cfg Config) (*ConsulDriver, error) {
	client, err := consul.NewClient(&consul.Config{
		Address: cfg.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("create consul client: %w", err)
	}

	driver := &ConsulDriver{
		cfg:    cfg,
		client: client,
	}

	return driver, nil
}

// Register 注册服务实例
func (d *ConsulDriver) Register(instance ServiceInstance) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	registration := new(consul.AgentServiceRegistration)
	registration.ID = instance.ID
	registration.Name = instance.Name
	registration.Address = instance.Address
	registration.Port = instance.Port
	registration.Tags = instance.Tags
	registration.Meta = instance.Metadata

	// 添加健康检查
	registration.Check = &consul.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", instance.Address, instance.Port, "/health"),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}

	if err := d.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("register service: %w", err)
	}

	return nil
}

// Deregister 注销服务实例
func (d *ConsulDriver) Deregister(instanceID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.client.Agent().ServiceDeregister(instanceID); err != nil {
		return fmt.Errorf("deregister service: %w", err)
	}

	return nil
}

// Discover 发现服务实例
func (d *ConsulDriver) Discover(serviceName string) ([]ServiceInstance, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	services, _, err := d.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("discover service: %w", err)
	}

	instances := make([]ServiceInstance, 0, len(services))
	for _, s := range services {
		instance := ServiceInstance{
			ID:       s.Service.ID,
			Name:     s.Service.Service,
			Address:  s.Service.Address,
			Port:     s.Service.Port,
			Tags:     s.Service.Tags,
			Metadata: s.Service.Meta,
			Healthy:  true,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// HealthCheck 健康检查
func (d *ConsulDriver) HealthCheck(serviceID string) (*HealthStatus, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	checks, _, err := d.client.Health().Checks(serviceID, nil)
	if err != nil {
		return nil, fmt.Errorf("health check: %w", err)
	}

	healthy := true
	message := "healthy"
	for _, check := range checks {
		if check.Status != consul.HealthPassing {
			healthy = false
			message = check.Notes
			break
		}
	}

	return &HealthStatus{
		ServiceID: serviceID,
		Healthy:   healthy,
		Message:   message,
	}, nil
}

// Close 关闭连接
func (d *ConsulDriver) Close() error {
	return nil
}

// WatchService 监听服务变化（返回一个 channel）
func (d *ConsulDriver) WatchService(ctx context.Context, serviceName string) (<-chan []ServiceInstance, error) {
	stream := make(chan []ServiceInstance)

	go func() {
		defer close(stream)
		lastIndex := uint64(0)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				services, meta, err := d.client.Health().Service(serviceName, "", true, &consul.QueryOptions{
					WaitIndex: lastIndex,
				})
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				if lastIndex == meta.LastIndex {
					continue
				}

				lastIndex = meta.LastIndex

				instances := make([]ServiceInstance, 0, len(services))
				for _, s := range services {
					instance := ServiceInstance{
						ID:       s.Service.ID,
						Name:     s.Service.Service,
						Address:  s.Service.Address,
						Port:     s.Service.Port,
						Tags:     s.Service.Tags,
						Metadata: s.Service.Meta,
						Healthy:  true,
					}
					instances = append(instances, instance)
				}

				select {
				case stream <- instances:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return stream, nil
}
