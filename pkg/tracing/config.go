package tracing

// Config 链路追踪配置
//
// YAML 配置示例:
//
//	tracing:
//	  enabled: true                                         # 是否启用链路追踪
//	  driver: "jaeger"                                      # 驱动: jaeger, zipkin, otel
//	  endpoint: "http://jaeger:14268/api/traces"            # 收集器地址
//	  service_name: "my-go-app"                             # 当前服务名称
//	  sample_rate: 0.1                                      # 采样率 (0.0 ~ 1.0)
type Config struct {
	Enabled     bool    `mapstructure:"enabled" yaml:"enabled" json:"enabled"`         // 是否启用
	Driver      string  `mapstructure:"driver" yaml:"driver" json:"driver"`           // 驱动: jaeger, zipkin, otel
	Endpoint    string  `mapstructure:"endpoint" yaml:"endpoint" json:"endpoint"`     // 收集器地址
	ServiceName string  `mapstructure:"service_name" yaml:"service_name" json:"service_name"` // 服务名称
	SampleRate  float64 `mapstructure:"sample_rate" yaml:"sample_rate" json:"sample_rate"` // 采样率
}

// DefaultConfig 返回链路追踪默认配置
func DefaultConfig() Config {
	return Config{
		Enabled:     false,
		Driver:      "jaeger",
		ServiceName: "aicode",
		SampleRate:  0.1,
	}
}
