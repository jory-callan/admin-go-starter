package internal

import (
	"aicode/internal/model"
	"aicode/pkg/logger"
	"log/slog"

	"gorm.io/gorm"
)

var log *slog.Logger

// Migrate 自动迁移表结构
func Migrate(db *gorm.DB) {
	log = logger.C("db")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		panic(err)
	}

	log.Info("database migrated successfully")
}
