package model

import "time"

// UserRole 用户角色关联表
type UserRole struct {
	UserID    string    `gorm:"type:varchar(36);primaryKey;comment:用户ID" json:"user_id"`
	RoleID    string    `gorm:"type:varchar(36);primaryKey;comment:角色ID" json:"role_id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       string    `gorm:"type:varchar(36);primaryKey;comment:角色ID" json:"role_id"`
	PermissionID string    `gorm:"type:varchar(36);primaryKey;comment:权限ID" json:"permission_id"`
	CreatedAt    time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
