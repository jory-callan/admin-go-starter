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
		&model.UserRole{},
		&model.RolePermission{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		panic(err)
	}

	slog.Info("database migrated successfully")
	SeedDefaultData(db)

}

// 通过 uuidv7 id 主键，判断，如果数据存在，则跳过，不存在则插入初始化数据
func SeedDefaultData(db *gorm.DB) {
	// 初始化 user
	if err := db.Create(&model.User{
		ID:       "019d3949a8ed75aea9dcdba4a3b8a665",
		Username: "admin",
		Password: "123456",
		Nickname: "管理员",
		Email:    "1219946450@qq.com",
		Phone:    "13012345678",
		Avatar:   "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
		Status:   1,
	}).Error; err != nil {
		slog.Error("failed to seed user data", "error", err)
	}

	// 初始化 role
	if err := db.Create(&model.Role{
		ID:          "019d394a9ba47243a8bd0d587028deaf",
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员",
		Sort:        1,
		Status:      1,
	}).Error; err != nil {
		slog.Error("failed to seed role data", "error", err)
	}

	// 初始化 permission
	if err := db.Create(&model.Permission{
		ID:          "019d394b3d4b7bf2abef7375de4672b9",
		Name:        "管理权限",
		Code:        "*",
		Description: "拥有所有权限",
		Sort:        1,
		Status:      1,
	}).Error; err != nil {
		slog.Error("failed to seed permission data", "error", err)
	}

	// 初始化 user_role 关联
	if err := db.Create(&model.UserRole{
		UserID: "019d3949a8ed75aea9dcdba4a3b8a665",
		RoleID: "019d394a9ba47243a8bd0d587028deaf",
	}).Error; err != nil {
		slog.Error("failed to seed user_role data", "error", err)
	}

	// 初始化 role_permission 关联
	if err := db.Create(&model.RolePermission{
		RoleID:       "019d394a9ba47243a8bd0d587028deaf",
		PermissionID: "019d394b3d4b7bf2abef7375de4672b9",
	}).Error; err != nil {
		slog.Error("failed to seed role_permission data", "error", err)
	}

}
