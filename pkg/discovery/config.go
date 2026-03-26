package discovery

// Config 服务发现与注册配置
//
// YAML 配置示例:
//
//	service_discovery:
//	  enabled: true                              # 是否启用服务发现
//	  driver: "consul"                           # 驱动: consul, etcd, nacos
//	  address: "127.0.0.1:8500"                  # 服务发现服务器地址
//	  service_name: "my-go-app"                  # 当前服务名称
//	  service_port: 8080                         # 当前服务端口
//	  health_check_path: "/health"               # 健康检查接口路径
type Config struct {
	Enabled         bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`                   // 是否启用
	Driver          string `mapstructure:"driver" yaml:"driver" json:"driver"`                       // 驱动: consul, etcd, nacos
	Address         string `mapstructure:"address" yaml:"address" json:"address"`                   // 服务器地址
	ServiceName     string `mapstructure:"service_name" yaml:"service_name" json:"service_name"`     // 当前服务名称
	ServicePort     int    `mapstructure:"service_port" yaml:"service_port" json:"service_port"`     // 当前服务端口
	HealthCheckPath string `mapstructure:"health_check_path" yaml:"health_check_path" json:"health_check_path"` // 健康检查路径
}

// DefaultConfig 返回服务发现默认配置
func DefaultConfig() Config {
	return Config{
		Enabled:         false,
		Driver:          "consul",
		Address:         "127.0.0.1:8500",
		ServiceName:     "aicode",
		ServicePort:     8080,
		HealthCheckPath: "/health",
	}
}
