package migration

import (
	"aicode/internal/model"
	"log/slog"

	"gorm.io/gorm"
)

var log *slog.Logger

// Migrate 自动迁移表结构
func Migrate(db *gorm.DB) {

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		panic(err)
	}

	slog.Info("database migrated successfully")
}
