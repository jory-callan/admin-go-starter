package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config Redis 配置
type Config struct {
	Addr         string `mapstructure:"addr" yaml:"addr"`                   // Redis 地址 host:port
	Password     string `mapstructure:"password" yaml:"password"`           // 密码
	DB           int    `mapstructure:"db" yaml:"db"`                       // 数据库索引
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size"`         // 连接池大小
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"` // 最小空闲连接数
	DialTimeout  int    `mapstructure:"dial_timeout" yaml:"dial_timeout"`   // 连接超时(秒)
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`   // 读超时(秒)
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"` // 写超时(秒)
}

// GetDefault 返回 Redis 默认配置
func GetDefault() Config {
	return Config{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
	}
}

// Open 根据 Config 创建 *redis.Client 实例
func Open(cfg Config, log *slog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DialTimeout)*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis (addr=%s): %w", cfg.Addr, err)
	}

	log.Info("redis connected", "addr", cfg.Addr, "db", cfg.DB)
	return client, nil
}
