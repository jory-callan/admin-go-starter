package core

import (
	"aicode/pkg/redis"
	"fmt"
)

// initRedis 初始化所有 Redis 实例
func (a *App) initRedis() {
	cfg := a.Config
	var msg string

	// 默认 Redis（必填）
	rdb := redis.New(cfg.Redis)
	a.Redis = rdb
	msg = fmt.Sprintf("redis initialized, name: %s, addr: %s", "redis", cfg.Redis.Addr)
	a.Log.Info(msg)
	a.registerRedisCloser("redis", rdb)
}

// registerRedisCloser 注册 Redis 关闭钩子
func (a *App) registerRedisCloser(name string, client interface{ Close() error }) {
	var msg string
	a.registerCloser(func() error {
		msg = fmt.Sprintf("closing redis, name: %s", name)
		a.Log.Info(msg)
		client.Close()
		msg = fmt.Sprintf("redis closed, name: %s", name)
		a.Log.Info(msg)
		return nil
	})
}
