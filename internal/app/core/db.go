package core

import (
	"aicode/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

// initDatabases 初始化所有数据库实例
func (a *App) initDatabases() {
	cfg := a.Config
	var msg string
	// 默认主数据库（必填）
	db := database.New(cfg.Database)
	a.DB = db
	msg = fmt.Sprintf("database initialized, name: %s, driver: %s", "database", cfg.Database.Driver)
	a.Log.Info(msg)
	a.registerDBCloser("database", db)
}

// registerDBCloser 注册数据库关闭钩子
func (a *App) registerDBCloser(name string, db *gorm.DB) {
	var msg string
	a.registerCloser(func() error {
		msg = fmt.Sprintf("closing database, name: %s", name)
		a.Log.Info(msg)
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		msg = fmt.Sprintf("database closed, name: %s", name)
		a.Log.Info(msg)
		return nil
	})
}
