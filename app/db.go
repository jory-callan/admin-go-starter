package app

import (
	"aicode/pkg/database"
	"fmt"
)

// initDatabases 初始化所有数据库实例
// 从 App.Config.Database 读取配置，创建 *gorm.DB 并存入 App.Databases 和 App.DB
func (a *App) initDatabases() error {
	if len(a.Config.Database) == 0 {
		a.Log.Warn("no database configured")
		return nil
	}

	for name, cfg := range a.Config.Database {
		db, err := database.Open(cfg, a.Log)
		if err != nil {
			return fmt.Errorf("init database[%s]: %w", name, err)
		}

		a.Databases[name] = db

		// 默认实例同时挂载到 App.DB
		if name == "default" {
			a.DB = db
		}

		// 注册关闭钩子
		a.registerCloser(func() error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			a.Log.Info("closing database", "name", name)
			return sqlDB.Close()
		})
	}

	if a.DB == nil {
		return fmt.Errorf("database 'default' instance is required but not configured")
	}

	a.Log.Info("all databases initialized", "count", len(a.Databases))
	return nil
}
