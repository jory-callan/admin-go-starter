package redis

import (
	"context"
	"aicode/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log/slog"
)

var log *slog.Logger

// New 初始化 Redis
func New(conf *viper.Viper) *redis.Client {
	log = logger.C("redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("redis.addr"),
		Password: conf.GetString("redis.password"),
		DB:       conf.GetInt("redis.db"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Error("failed to connect to redis", "error", err)
		panic(err)
	}

	log.Info("redis connected successfully")
	return rdb
}
