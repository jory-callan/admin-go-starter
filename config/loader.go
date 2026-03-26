package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Load 加载配置文件并返回强类型 AppConfig
// 流程: DefaultConfig() → ReadInConfig() → Unmarshal()
// 配置文件中未指定的字段自动使用 DefaultConfig() 中的默认值
func Load(configFile string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 1. 获取默认值并设置到 viper（智能合并的基础）
	defaults := DefaultConfig()
	if err := v.Unmarshal(&defaults); err != nil {
		return nil, fmt.Errorf("prepare defaults: %w", err)
	}

	// 将默认值逐 key 写入 viper.SetDefault，确保未配置的字段也能被 Unmarshal 填充
	if err := setDefaultsFromStruct(v, defaults); err != nil {
		return nil, fmt.Errorf("set defaults: %w", err)
	}

	// 2. 查找并读取配置文件
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// 默认查找顺序
		for _, p := range []string{".", "./config", "./conf"} {
			v.AddConfigPath(p)
		}
		v.SetConfigName("config")
	}

	// ReadInConfig 可选：配置文件不存在不报错（使用默认值）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	// 3. Unmarshal 到强类型
	var cfg AppConfig
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
	)
	if err := v.Unmarshal(&cfg, viper.DecodeHook(decodeHook)); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaultsFromStruct 将 AppConfig 的默认值设置到 viper 中
// 确保 Unmarshal 时未在配置文件中指定的字段也能被填充
func setDefaultsFromStruct(v *viper.Viper, cfg AppConfig) error {
	// Log
	v.SetDefault("log.level", cfg.Log.Level)
	v.SetDefault("log.format", cfg.Log.Format)
	v.SetDefault("log.output", cfg.Log.Output)
	v.SetDefault("log.file_path", cfg.Log.FilePath)
	v.SetDefault("log.max_size", cfg.Log.MaxSize)
	v.SetDefault("log.max_backups", cfg.Log.MaxBackups)
	v.SetDefault("log.max_age", cfg.Log.MaxAge)
	v.SetDefault("log.compress", cfg.Log.Compress)

	// HTTP
	v.SetDefault("http.enable_debug", cfg.HTTP.EnableDebug)
	v.SetDefault("http.host", cfg.HTTP.Host)
	v.SetDefault("http.port", cfg.HTTP.Port)
	v.SetDefault("http.read_timeout", cfg.HTTP.ReadTimeout)
	v.SetDefault("http.write_timeout", cfg.HTTP.WriteTimeout)
	v.SetDefault("http.idle_timeout", cfg.HTTP.IdleTimeout)
	v.SetDefault("http.max_header_bytes", cfg.HTTP.MaxHeaderBytes)
	v.SetDefault("http.max_body_size", cfg.HTTP.MaxBodySize)

	// Database (默认主库)
	v.SetDefault("database.driver", cfg.Database.Driver)
	v.SetDefault("database.dsn", cfg.Database.DSN)
	v.SetDefault("database.max_open_conns", cfg.Database.MaxOpenConns)
	v.SetDefault("database.max_idle_conns", cfg.Database.MaxIdleConns)
	v.SetDefault("database.conn_max_lifetime", cfg.Database.ConnMaxLifetime)
	v.SetDefault("database.conn_max_idle_time", cfg.Database.ConnMaxIdleTime)
	v.SetDefault("database.log_level", cfg.Database.LogLevel)

	// Redis (默认)
	v.SetDefault("redis.addr", cfg.Redis.Addr)
	v.SetDefault("redis.password", cfg.Redis.Password)
	v.SetDefault("redis.db", cfg.Redis.DB)
	v.SetDefault("redis.pool_size", cfg.Redis.PoolSize)
	v.SetDefault("redis.min_idle_conns", cfg.Redis.MinIdleConns)
	v.SetDefault("redis.dial_timeout", cfg.Redis.DialTimeout)
	v.SetDefault("redis.read_timeout", cfg.Redis.ReadTimeout)
	v.SetDefault("redis.write_timeout", cfg.Redis.WriteTimeout)

	// JWT
	v.SetDefault("jwt.secret", cfg.JWT.Secret)
	v.SetDefault("jwt.expires", cfg.JWT.Expires)
	v.SetDefault("jwt.issuer", cfg.JWT.Issuer)

	// ServiceDiscovery
	v.SetDefault("service_discovery.enabled", cfg.ServiceDiscovery.Enabled)
	v.SetDefault("service_discovery.driver", cfg.ServiceDiscovery.Driver)
	v.SetDefault("service_discovery.address", cfg.ServiceDiscovery.Address)
	v.SetDefault("service_discovery.service_name", cfg.ServiceDiscovery.ServiceName)
	v.SetDefault("service_discovery.service_port", cfg.ServiceDiscovery.ServicePort)
	v.SetDefault("service_discovery.health_check_path", cfg.ServiceDiscovery.HealthCheckPath)

	// Tracing
	v.SetDefault("tracing.enabled", cfg.Tracing.Enabled)
	v.SetDefault("tracing.driver", cfg.Tracing.Driver)
	v.SetDefault("tracing.endpoint", cfg.Tracing.Endpoint)
	v.SetDefault("tracing.service_name", cfg.Tracing.ServiceName)
	v.SetDefault("tracing.sample_rate", cfg.Tracing.SampleRate)

	return nil
}
