package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	Driver            string `mapstructure:"driver" yaml:"driver"`                       // mysql, postgres, sqlite
	DSN               string `mapstructure:"dsn" yaml:"dsn"`                             // 连接串
	MaxOpenConns      int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`       // 最大打开连接数
	MaxIdleConns      int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`       // 最大空闲连接数
	ConnMaxLifetime   int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 连接最大存活时间(秒)
	ConnMaxIdleTime   int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 连接最大空闲时间(秒)
	LogLevel          string `mapstructure:"log_level" yaml:"log_level"`                 // gorm 日志级别
}

// GetDefault 返回数据库默认配置
func GetDefault() Config {
	return Config{
		Driver:          "sqlite",
		DSN:             "demo.sqlite.db",
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: 1800, // 30 分钟
		ConnMaxIdleTime: 300,  // 5 分钟
		LogLevel:        "warn",
	}
}

// parseGormLogLevel 将字符串转为 gorm logger level
func parseGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}

// Open 根据 Config 创建 *gorm.DB 实例
func Open(cfg Config, log *slog.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(parseGormLogLevel(cfg.LogLevel)),
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database (driver=%s): %w", cfg.Driver, err)
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	log.Info("database connected", "driver", cfg.Driver)

	return db, nil
}
