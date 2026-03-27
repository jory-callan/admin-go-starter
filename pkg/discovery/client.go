package discovery

import (
	"fmt"
	"sync"
)

// Client 服务发现客户端（统一入口）
type Client struct {
	cfg    Config
	driver Driver
	mu     sync.RWMutex
}

// NewClient 创建服务发现客户端
func NewClient(cfg Config) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil // 未启用服务发现
	}

	client := &Client{
		cfg: cfg,
	}

	// 根据 driver 类型创建对应的驱动
	var driver Driver
	var err error

	switch cfg.Driver {
	case "consul":
		driver, err = NewConsulDriver(cfg)
	default:
		return nil, fmt.Errorf("unsupported discovery driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, err
	}

	client.driver = driver
	return client, nil
}

// Register 注册服务实例
func (c *Client) Register(instance ServiceInstance) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.driver == nil {
		return nil // 未启用服务发现
	}

	return c.driver.Register(instance)
}

// Deregister 注销服务实例
func (c *Client) Deregister(instanceID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.driver == nil {
		return nil
	}

	return c.driver.Deregister(instanceID)
}

// Discover 发现服务实例
func (c *Client) Discover(serviceName string) ([]ServiceInstance, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.driver == nil {
		return nil, nil
	}

	return c.driver.Discover(serviceName)
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(serviceID string) (*HealthStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.driver == nil {
		return nil, nil
	}

	return c.driver.HealthCheck(serviceID)
}

// Close 关闭客户端
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.driver == nil {
		return nil
	}

	return c.driver.Close()
}

// GetDriver 获取底层驱动（高级用法）
func (c *Client) GetDriver() Driver {
	return c.driver
}
