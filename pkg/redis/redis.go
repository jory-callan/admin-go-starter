package redis

import (
	"aicode/pkg/logger"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// New 根据 Config 创建 *redis.Client 实例
func New(cfg Config) *redis.Client {
	rdLog := logger.C("redis")

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
		panic(fmt.Errorf("failed to connect redis (addr=%s): %w", cfg.Addr, err))
	}

	rdLog.Info("redis connection established", "addr", cfg.Addr, "db", cfg.DB, "pool_size", cfg.PoolSize)
	return client
}

func Shutdown(client *redis.Client) error {
	return client.Close()
}
