package app

import (
	"aicode/pkg/database"

	"gorm.io/gorm"
)

// initDatabases 初始化所有数据库实例
func (a *App) initDatabases() {
	cfg := a.Config

	// 默认主数据库（必填）
	db := database.New(cfg.Database)
	a.DB = db
	a.registerDBCloser("database", db)
	a.Log.Info("database initialized", "name", "database", "driver", cfg.Database.Driver)
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
