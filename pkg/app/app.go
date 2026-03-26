package app

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// App 基础设施结构体，包含与业务逻辑相关的资源
type App struct {
	DB    *gorm.DB
	Redis *redis.Client
	Conf  *viper.Viper
}
