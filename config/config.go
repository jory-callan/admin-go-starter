package config

import (
	"time"

	"aicode/pkg/database"
	pkglogger "aicode/pkg/logger"
	pkgredis "aicode/pkg/redis"
)

// AppConfig 应用全局配置（强类型）
type AppConfig struct {
	Log      pkglogger.Config               `mapstructure:"log" yaml:"log"`
	HTTP     HTTPConfig                      `mapstructure:"http" yaml:"http"`
	Database map[string]database.Config      `mapstructure:"database" yaml:"database"` // 多实例数据库
	Redis    map[string]pkgredis.Config      `mapstructure:"redis" yaml:"redis"`       // 多实例 Redis
	JWT      JWTConfig                       `mapstructure:"jwt" yaml:"jwt"`
}

// HTTPConfig HTTP 服务配置
type HTTPConfig struct {
	Host           string        `mapstructure:"host" yaml:"host"`
	Port           int           `mapstructure:"port" yaml:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes" yaml:"max_header_bytes"` // 字节
	MaxBodySize    int           `mapstructure:"max_body_size" yaml:"max_body_size"`       // 字节
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret  string        `mapstructure:"secret" yaml:"secret"`
	Expires time.Duration `mapstructure:"expires" yaml:"expires"`
	Issuer  string        `mapstructure:"issuer" yaml:"issuer"`
}

// GetDefault 返回全局默认配置
// 各子配置的默认值由 pkg 层各自的 GetDefault() 提供
func GetDefault() AppConfig {
	return AppConfig{
		Log: pkglogger.GetDefault(),
		HTTP: HTTPConfig{
			Host:           "0.0.0.0",
			Port:           8080,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,  // 1MB
			MaxBodySize:    10 << 20, // 10MB
		},
		Database: map[string]database.Config{
			"default": database.GetDefault(),
		},
		Redis: map[string]pkgredis.Config{
			"default": pkgredis.GetDefault(),
		},
		JWT: JWTConfig{
			Secret:  "change-me-in-production-32chars!",
			Expires: 24 * time.Hour,
			Issuer:  "aicode",
		},
	}
}
