package http

// Config HTTP 服务配置
//
// YAML 配置示例:
//
//	http:
//	  enable_debug: false                              # 是否启用 Debug 模式
//	  host: "0.0.0.0"                                  # 监听地址
//	  port: 8080                                       # 监听端口
//	  read_timeout: "10s"                              # 读取请求超时时间
//	  write_timeout: "10s"                             # 写入响应超时时间
//	  idle_timeout: "60s"                              # 空闲连接超时时间
//	  max_header_bytes: "1MB"                          # 最大允许的请求头大小 (支持 KB, MB, GB)
//	  max_body_size: "10MB"                            # 最大允许的请求体大小 (支持 KB, MB, GB)
//	  cors:                                           # CORS 跨域配置（可选）
//	    allow_origins: ["https://example.com"]         # 允许的源
//	    allow_methods: ["GET", "POST", "PUT", "DELETE"] # 允许的 HTTP 方法
//	    allow_headers: ["Content-Type", "Authorization"] # 允许的请求头
//	    expose_headers: ["X-Request-ID"]               # 暴露给客户端的响应头
//	    max_age: 86400                                 # 预检请求缓存时间(秒)
//	  rate_limit:                                      # 限流配置（可选）
//	    requests_per_second: 100                       # 每秒允许的请求数
//	    burst: 20                                      # 突发请求允许的最大数量
type Config struct {
	EnableDebug   bool            `mapstructure:"enable_debug" yaml:"enable_debug" json:"enable_debug"`     // 是否启用 Debug
	Host          string          `mapstructure:"host" yaml:"host" json:"host"`                                // 监听地址
	Port          int             `mapstructure:"port" yaml:"port" json:"port"`                                // 监听端口
	ReadTimeout   string          `mapstructure:"read_timeout" yaml:"read_timeout" json:"read_timeout"`       // 读取超时
	WriteTimeout  string          `mapstructure:"write_timeout" yaml:"write_timeout" json:"write_timeout"`    // 写入超时
	IdleTimeout   string          `mapstructure:"idle_timeout" yaml:"idle_timeout" json:"idle_timeout"`       // 空闲超时
	MaxHeaderBytes string         `mapstructure:"max_header_bytes" yaml:"max_header_bytes" json:"max_header_bytes"` // 最大请求头
	MaxBodySize   string          `mapstructure:"max_body_size" yaml:"max_body_size" json:"max_body_size"`    // 最大请求体
	CORS          CORSConfig      `mapstructure:"cors" yaml:"cors" json:"cors"`                                // CORS 跨域配置
	RateLimit     RateLimitConfig `mapstructure:"rate_limit" yaml:"rate_limit" json:"rate_limit"`             // 限流配置
}

// CORSConfig CORS 跨域配置
type CORSConfig struct {
	AllowOrigins  []string `mapstructure:"allow_origins" yaml:"allow_origins" json:"allow_origins"`       // 允许的源
	AllowMethods  []string `mapstructure:"allow_methods" yaml:"allow_methods" json:"allow_methods"`       // 允许的方法
	AllowHeaders  []string `mapstructure:"allow_headers" yaml:"allow_headers" json:"allow_headers"`       // 允许的请求头
	ExposeHeaders []string `mapstructure:"expose_headers" yaml:"expose_headers" json:"expose_headers"`   // 暴露的响应头
	MaxAge        int      `mapstructure:"max_age" yaml:"max_age" json:"max_age"`                         // 预检缓存时间(秒)
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerSecond float64 `mapstructure:"requests_per_second" yaml:"requests_per_second" json:"requests_per_second"` // 每秒请求数
	Burst             int     `mapstructure:"burst" yaml:"burst" json:"burst"`                                                 // 突发最大数量
}

// DefaultConfig 返回 HTTP 默认配置
func DefaultConfig() Config {
	return Config{
		Host:           "0.0.0.0",
		Port:           8080,
		ReadTimeout:    "10s",
		WriteTimeout:   "10s",
		IdleTimeout:    "60s",
		MaxHeaderBytes: "1MB",
		MaxBodySize:    "10MB",
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			MaxAge:       86400,
		},
	}
}
