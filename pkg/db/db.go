package db

import (
	"aicode/pkg/logger"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log *slog.Logger

// New 根据 driver 类型初始化数据库连接
// 支持: mysql / postgres / sqlite，默认 sqlite
func New(conf *viper.Viper) *gorm.DB {
	log = logger.C("db")

	driver := conf.GetString("database.driver")
	if driver == "" {
		driver = "sqlite"
	}

	dsn := conf.GetString("database.dsn")

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		log.Error("unsupported database driver", "driver", driver)
		panic(fmt.Sprintf("unsupported database driver: %s", driver))
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		panic(err)
	}

	log.Info("database connected", "driver", driver)
	return db
}
