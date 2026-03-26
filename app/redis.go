package app

import (
	"aicode/pkg/redis"
	"fmt"
)

// initRedis 初始化所有 Redis 实例
// 从 App.Config.Redis 读取配置，创建 *redis.Client 并存入 App.RedisMap 和 App.Redis
func (a *App) initRedis() error {
	if len(a.Config.Redis) == 0 {
		a.Log.Warn("no redis configured")
		return nil
	}

	for name, cfg := range a.Config.Redis {
		client, err := redis.Open(cfg, a.Log)
		if err != nil {
			return fmt.Errorf("init redis[%s]: %w", name, err)
		}

		a.RedisMap[name] = client

		// 默认实例同时挂载到 App.Redis
		if name == "default" {
			a.Redis = client
		}

		// 注册关闭钩子
		a.registerCloser(func() error {
			a.Log.Info("closing redis", "name", name)
			return client.Close()
		})
	}

	if a.Redis == nil {
		return fmt.Errorf("redis 'default' instance is required but not configured")
	}

	a.Log.Info("all redis instances initialized", "count", len(a.RedisMap))
	return nil
}
