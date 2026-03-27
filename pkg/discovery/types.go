package discovery

// ServiceInstance 服务实例信息
type ServiceInstance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	Healthy   bool              `json:"healthy"`
}

// HealthStatus 健康状态
type HealthStatus struct {
	ServiceID string
	Healthy   bool
	Message   string
}

// Driver 服务发现驱动接口
type Driver interface {
	// Register 注册服务实例
	Register(instance ServiceInstance) error
	
	// Deregister 注销服务实例
	Deregister(instanceID string) error
	
	// Discover 发现服务实例
	Discover(serviceName string) ([]ServiceInstance, error)
	
	// HealthCheck 健康检查
	HealthCheck(serviceID string) (*HealthStatus, error)
	
	// Close 关闭连接
	Close() error
}
