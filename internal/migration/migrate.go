package migration

import (
	"log/slog"

	"aicode/internal/model"

	"gorm.io/gorm"
)

var log *slog.Logger

// Migrate 自动迁移表结构
func Migrate(db *gorm.DB) {

	if err := db.AutoMigrate(
		// user
		&model.User{},
		&model.Role{},
		&model.UserRole{},
	); err != nil {
		slog.Error("failed to migrate database", "error", err)
	}

	slog.Info("database migrated successfully")
	SeedDefaultData(db)
}

// 通过 db.Create 创建初始必备数据，忽略报错，说明这个主键的数据已经存在了
func SeedDefaultData(db *gorm.DB) {
	// 初始化默认管理员角色
	adminRole := model.Role{
		ID:   "019d394a9ba47243a8bd0d587028deaf",
		Code: "admin",
		Name: "管理员",
	}
	if err := db.Create(&adminRole).Error; err != nil {
		slog.Error("failed to seed admin role", "error", err)
	}

	// 初始化默认用户 (密码: admin123)
	// 密码哈希: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
	// 密码 123456
	// 哈希：$2a$10$gH6bm5s9tRG72FVMW/nVYeliwqTChISyggta7A4D/5JsRbx6b6iie
	adminUser := model.User{
		ID:       "019d3949a8ed75aea9dcdba4a3b8a665",
		Username: "admin",
		Password: "$2a$10$gH6bm5s9tRG72FVMW/nVYeliwqTChISyggta7A4D/5JsRbx6b6iie",
		Email:    "admin@example.com",
		Status:   1,
	}
	if err := db.Create(&adminUser).Error; err != nil {
		slog.Error("failed to seed admin user", "error", err)
	}

	// 绑定默认管理员用户和角色关系
	userRole := model.UserRole{
		UserID: adminUser.ID,
		RoleID: adminRole.ID,
	}
	if err := db.Create(&userRole).Error; err != nil {
		slog.Error("failed to seed user_role relation", "error", err)
	}

}
