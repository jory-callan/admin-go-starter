package app

import (
	"aicode/pkg/redis"
	"fmt"
)

// initRedis 初始化所有 Redis 实例
func (a *App) initRedis() error {
	cfg := a.Config

	// 默认 Redis（必填）
	rdb, err := redis.Open(cfg.Redis)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	a.Redis = rdb
	a.registerRedisCloser("redis", rdb)
	a.Log.Info("redis initialized", "name", "redis", "addr", cfg.Redis.Addr)

	// 缓存专用 Redis（可选）
	if cfg.CacheRedis != nil {
		cacheRDB, err := redis.Open(*cfg.CacheRedis)
		if err != nil {
			return fmt.Errorf("init cache_redis: %w", err)
		}
		a.CacheRedis = cacheRDB
		a.registerRedisCloser("cache_redis", cacheRDB)
		a.Log.Info("redis initialized", "name", "cache_redis", "addr", cfg.CacheRedis.Addr)
	}

	// Session 专用 Redis（可选）
	if cfg.SessionRedis != nil {
		sessionRDB, err := redis.Open(*cfg.SessionRedis)
		if err != nil {
			return fmt.Errorf("init session_redis: %w", err)
		}
		a.SessionRedis = sessionRDB
		a.registerRedisCloser("session_redis", sessionRDB)
		a.Log.Info("redis initialized", "name", "session_redis", "addr", cfg.SessionRedis.Addr)
	}

	return nil
}

// registerRedisCloser 注册 Redis 关闭钩子
func (a *App) registerRedisCloser(name string, client interface{ Close() error }) {
	a.registerCloser(func() error {
		a.Log.Info("closing redis", "name", name)
		return client.Close()
	})
}
