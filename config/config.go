package config

import (
	"aicode/pkg/database"
	pkgdiscovery "aicode/pkg/discovery"
	pkghttp "aicode/pkg/http"
	pkgjwt "aicode/pkg/jwt"
	pkgkafka "aicode/pkg/kafka"
	pkglogger "aicode/pkg/logger"
	pkgrabbitmq "aicode/pkg/rabbitmq"
	pkgredis "aicode/pkg/redis"
	pkgtracing "aicode/pkg/tracing"
)

// AppConfig 应用全局配置（强类型）
// 所有子配置均引用 pkg 层各自的 Config 结构体，统一由各包的 DefaultConfig() 提供默认值
type AppConfig struct {
	Log             pkglogger.Config         `mapstructure:"log" yaml:"log"`
	HTTP            pkghttp.Config            `mapstructure:"http" yaml:"http"`
	Database        database.Config           `mapstructure:"database" yaml:"database"`
	LogDatabase     *database.Config          `mapstructure:"log_database" yaml:"log_database"`
	AnalyticsDB     *database.Config          `mapstructure:"analytics_database" yaml:"analytics_database"`
	Redis           pkgredis.Config           `mapstructure:"redis" yaml:"redis"`
	CacheRedis      *pkgredis.Config          `mapstructure:"cache_redis" yaml:"cache_redis"`
	SessionRedis    *pkgredis.Config          `mapstructure:"session_redis" yaml:"session_redis"`
	JWT             pkgjwt.Config             `mapstructure:"jwt" yaml:"jwt"`
	Kafka           *pkgkafka.Config          `mapstructure:"kafka" yaml:"kafka"`
	RabbitMQ        *pkgrabbitmq.Config       `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	ServiceDiscovery pkgdiscovery.Config       `mapstructure:"service_discovery" yaml:"service_discovery"`
	Tracing         pkgtracing.Config         `mapstructure:"tracing" yaml:"tracing"`
}

// DefaultConfig 返回全局默认配置
// 各子配置的默认值由 pkg 层各自的 DefaultConfig() 提供
func DefaultConfig() AppConfig {
	return AppConfig{
		Log:             pkglogger.DefaultConfig(),
		HTTP:            pkghttp.DefaultConfig(),
		Database:        database.DefaultConfig(),
		Redis:           pkgredis.DefaultConfig(),
		JWT:             pkgjwt.DefaultConfig(),
		ServiceDiscovery: pkgdiscovery.DefaultConfig(),
		Tracing:         pkgtracing.DefaultConfig(),
	}
}
