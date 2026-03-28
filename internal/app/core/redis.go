package core

import (
	"aicode/pkg/redis"
)

// initRedis 初始化所有 Redis 实例
func (a *App) initRedis() {
	cfg := a.Config

	// 默认 Redis（必填）
	rdb := redis.New(cfg.Redis)
	a.Redis = rdb
	a.registerRedisCloser("redis", rdb)
	a.Log.Info("redis initialized", "name", "redis", "addr", cfg.Redis.Addr)
}

// registerRedisCloser 注册 Redis 关闭钩子
func (a *App) registerRedisCloser(name string, client interface{ Close() error }) {
	a.registerCloser(func() error {
		a.Log.Info("closing redis", "name", name)
		return client.Close()
	})
}
