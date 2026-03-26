package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Load 加载配置文件并返回强类型 AppConfig
// 流程: Default() -> ReadInConfig() -> Unmarshal()
func Load(configFile string) (*AppConfig, error) {
	v := viper.New()

	// 1. 先填充默认值
	defaults := GetDefault()
	v.SetConfigType("yaml")

	// 将 AppConfig 默认值写入 viper
	if err := setDefaultsFromAppConfig(v, defaults); err != nil {
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
		// 配置文件不存在，使用纯默认值
	}

	// 3. Unmarshal 到强类型
	var cfg AppConfig
	// 使用自定义 DecoderConfig，正确处理 mapstructure tag
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
	)
	if err := v.Unmarshal(&cfg, viper.DecodeHook(decodeHook)); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaultsFromAppConfig 将 AppConfig 的默认值设置到 viper 中
// 实现智能合并：配置文件中未指定的字段自动使用默认值
func setDefaultsFromAppConfig(v *viper.Viper, cfg AppConfig) error {
	// Log
	setNestedDefault(v, "log.level", cfg.Log.Level)
	setNestedDefault(v, "log.format", cfg.Log.Format)
	setNestedDefault(v, "log.output", cfg.Log.Output)
	setNestedDefault(v, "log.max_size", cfg.Log.MaxSize)
	setNestedDefault(v, "log.max_backups", cfg.Log.MaxBackups)
	setNestedDefault(v, "log.max_age", cfg.Log.MaxAge)
	setNestedDefault(v, "log.compress", cfg.Log.Compress)

	// HTTP
	setNestedDefault(v, "http.host", cfg.HTTP.Host)
	setNestedDefault(v, "http.port", cfg.HTTP.Port)
	setNestedDefault(v, "http.read_timeout", cfg.HTTP.ReadTimeout.String())
	setNestedDefault(v, "http.write_timeout", cfg.HTTP.WriteTimeout.String())
	setNestedDefault(v, "http.idle_timeout", cfg.HTTP.IdleTimeout.String())
	setNestedDefault(v, "http.max_header_bytes", cfg.HTTP.MaxHeaderBytes)
	setNestedDefault(v, "http.max_body_size", cfg.HTTP.MaxBodySize)

	// Database (多实例)
	for name, dbCfg := range cfg.Database {
		prefix := "database." + name
		v.SetDefault(prefix+".driver", dbCfg.Driver)
		v.SetDefault(prefix+".dsn", dbCfg.DSN)
		v.SetDefault(prefix+".max_open_conns", dbCfg.MaxOpenConns)
		v.SetDefault(prefix+".max_idle_conns", dbCfg.MaxIdleConns)
		v.SetDefault(prefix+".conn_max_lifetime", dbCfg.ConnMaxLifetime)
		v.SetDefault(prefix+".conn_max_idle_time", dbCfg.ConnMaxIdleTime)
		v.SetDefault(prefix+".log_level", dbCfg.LogLevel)
	}

	// Redis (多实例)
	for name, rCfg := range cfg.Redis {
		prefix := "redis." + name
		v.SetDefault(prefix+".addr", rCfg.Addr)
		v.SetDefault(prefix+".password", rCfg.Password)
		v.SetDefault(prefix+".db", rCfg.DB)
		v.SetDefault(prefix+".pool_size", rCfg.PoolSize)
		v.SetDefault(prefix+".min_idle_conns", rCfg.MinIdleConns)
		v.SetDefault(prefix+".dial_timeout", rCfg.DialTimeout)
		v.SetDefault(prefix+".read_timeout", rCfg.ReadTimeout)
		v.SetDefault(prefix+".write_timeout", rCfg.WriteTimeout)
	}

	// JWT
	setNestedDefault(v, "jwt.secret", cfg.JWT.Secret)
	setNestedDefault(v, "jwt.expires", cfg.JWT.Expires.String())
	setNestedDefault(v, "jwt.issuer", cfg.JWT.Issuer)

	return nil
}

// setNestedDefault 辅助函数：设置嵌套的默认值
func setNestedDefault(v *viper.Viper, key string, value interface{}) {
	v.SetDefault(key, value)
}
