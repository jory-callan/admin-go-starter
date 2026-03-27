package database

import (
	"aicode/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// parseGormLogLevel 将字符串转为 gorm logger level
func parseGormLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn", "warning":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}

// New 根据 Config 创建 *gorm.DB 实例
func New(cfg Config) *gorm.DB {
	dbLog := logger.C("database")

	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgres", "postgresql", "pg", "pgsql":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(cfg.DSN)
	default:
		panic(fmt.Errorf("unsupported database driver: %s", cfg.Driver))
	}

	gormCfg := &gorm.Config{
		Logger: gormlogger.Default.LogMode(parseGormLogLevel(cfg.LogLevel)),
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect database (driver=%s): %w", cfg.Driver, err))
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get underlying sql.DB: %w", err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	dbLog.Info("database connection established", "driver", cfg.Driver, "max_open", cfg.MaxOpenConns, "max_idle", cfg.MaxIdleConns)

	return db
}
