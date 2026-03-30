package migration

import (
	"log/slog"

	"gorm.io/gorm"
)

var log *slog.Logger

// Migrate 自动迁移表结构
func Migrate(db *gorm.DB) {

	if err := db.AutoMigrate(
	// user

	); err != nil {
		log.Error("failed to migrate database", "error", err)
		panic(err)
	}

	slog.Info("database migrated successfully")
	SeedDefaultData(db)
}

// 通过 db.Create 创建初始必备数据，忽略报错，说明这个主键的数据已经存在了
func SeedDefaultData(db *gorm.DB) {
	// // 初始化
	// if err := db.Create(&model.UserRole{
	// 	UserID: "019d3949a8ed75aea9dcdba4a3b8a665",
	// 	RoleID: "019d394a9ba47243a8bd0d587028deaf",
	// }).Error; err != nil {
	// 	slog.Error("failed to seed user_role data", "error", err)
	// }

}
