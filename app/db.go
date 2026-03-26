package app

import (
	"aicode/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

// initDatabases 初始化所有数据库实例
func (a *App) initDatabases() error {
	cfg := a.Config

	// 默认主数据库（必填）
	db, err := database.Open(cfg.Database)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	a.DB = db
	a.registerDBCloser("database", db)
	a.Log.Info("database initialized", "name", "database", "driver", cfg.Database.Driver)

	// 日志专用数据库（可选）
	if cfg.LogDatabase != nil {
		logDB, err := database.Open(*cfg.LogDatabase)
		if err != nil {
			return fmt.Errorf("init log_database: %w", err)
		}
		a.LogDB = logDB
		a.registerDBCloser("log_database", logDB)
		a.Log.Info("database initialized", "name", "log_database", "driver", cfg.LogDatabase.Driver)
	}

	// 分析专用数据库（可选）
	if cfg.AnalyticsDB != nil {
		analyticsDB, err := database.Open(*cfg.AnalyticsDB)
		if err != nil {
			return fmt.Errorf("init analytics_database: %w", err)
		}
		a.AnalyticsDB = analyticsDB
		a.registerDBCloser("analytics_database", analyticsDB)
		a.Log.Info("database initialized", "name", "analytics_database", "driver", cfg.AnalyticsDB.Driver)
	}

	return nil
}

// registerDBCloser 注册数据库关闭钩子
func (a *App) registerDBCloser(name string, db *gorm.DB) {
	a.registerCloser(func() error {
		a.Log.Info("closing database", "name", name)
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	})
}
